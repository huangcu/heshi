package main

import (
	"context"
	"database/sql"
	"db/mysql"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var (
	db                  *sql.DB
	ctx                 context.Context
	cancelFn            context.CancelFunc
	serverIsInterrupted bool
	store               sessions.CookieStore
	activeCurrencyRate  *currency
)
var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var env string

func main() {
	flag.StringVar(&env, "env", "dev", "specifiy env dev or pro, default env - dev.")
	flag.Parse()
	os.Setenv("stage", env)

	ticker := time.NewTicker(time.Hour * 8)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := getLatestRates(); err != nil {
					util.FailToGetCurrencyExchangeAlert()
				}
				var err error
				activeCurrencyRate, err = getAcitveCurrencyRate()
				if err != nil {
					util.Println("fail to get latest active currency rate")
				}
			case <-stop:
				return
			}
		}
	}()
	defer func() {
		ticker.Stop()
		stop <- true
	}()

	// log.Fatal(startWebServer(":443"))
	log.Fatal(startWebServer(":8080"))
}

func startWebServer(port string) error {
	// cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return err
	// }
	// config := &tls.Config{Certificates: []tls.Certificate{cer}}

	r := gin.New()
	if os.Getenv("stage") != "pro" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		r.Use(gin.Logger())
	}

	r.Use(gin.Recovery())
	configRoute(r)
	// webServer := &http.Server{Addr: port, Handler: r, TLSConfig: config}
	webServer := &http.Server{Addr: port, Handler: r}
	return webServer.ListenAndServe()
	// return webServer.ListenAndServeTLS("server.crt", "server.key")
}

func configRoute(r *gin.Engine) {
	api := r.Group("/api")
	if os.Getenv("stage") == "pro" {
		api.Use(AuthMiddleWare())
		api.Use(sessions.Sessions("sessionid", store))
	}

	{
		apiAdmin := api.Group("admin")
		apiAdmin.Use(AdminSessionMiddleWare())
		{
			//admin user -
			apiAdmin.POST("/users", newAdminAgentUser)
			apiAdmin.GET("/users/:id", getUser)
			apiAdmin.GET("/users", getAllUsers)
			apiAdmin.PATCH("/users/:id", configAgent)
			apiAdmin.DELETE("/users/:id", removeUser)

			//currency rate
			apiAdmin.GET("/exchangerate", getCurrencyRate)
			apiAdmin.POST("/exchangerate", currencyRateReqValidator(newCurrencyRate))

			//discount
			apiAdmin.GET("/discounts/:id", getDiscount)
			apiAdmin.GET("/discounts", getAllDiscounts)
			apiAdmin.POST("/discount", newDiscount)

			//config
			apiAdmin.GET("/config", getConfig)
			apiAdmin.GET("/configs", getAllConfigs)
			apiAdmin.POST("/config", newConfig)

			//products with customize header
			apiAdmin.POST("/upload", uploadAndGetFileHeaders)
			apiAdmin.POST("/process/diamond", processDiamonds)

			//upload products by csv file
			apiAdmin.POST("/products/upload", uploadAndProcessProducts)

			//supplier
			apiAdmin.POST("/supplier", newSupplier)
			apiAdmin.GET("/suppliers", getAllSuppliers)
			apiAdmin.PUT("/suppliers/:id", updateSupplier)
			apiAdmin.DELETE("/suppliers/:id", removeSupplier)
		}
		//agent, customer
		api.POST("/users", newUser)
		api.PATCH("/users", UserSessionMiddleWare(), updateUser)
		api.GET("/users", UserSessionMiddleWare(), getUser)
		api.POST("/login", userLogin)
		api.POST("/logout", userLogout)
		api.GET("/users/:id/contactinfo", UserSessionMiddleWare(), agentContactInfo)

		//action- > add, delete
		api.POST("/shoppingList/:action", UserSessionMiddleWare(), toShoppingList)

		//products
		api.GET("/products", getAllProducts)
		api.GET("/products/diamonds", getAllDiamonds)
		api.GET("/products/small_diamonds", getAllSmallDiamonds)
		api.GET("/products/jewelrys", getAllJewelrys)
		api.GET("/products/diamonds/:id", getDiamond)
		api.GET("/products/small_diamonds/:id", getSmallDiamond)
		api.GET("/products/jewelrys/:id", getJewelry)
		api.POST("/products/diamonds", newDiamond)
		api.POST("/products/small_diamonds", newSmallDiamond)
		api.POST("/products/jewelrys", newJewelry)
		api.POST("/products/search", searchProducts)
		api.POST("/products/diamonds/search", searchProducts)
		api.POST("/products/jewelrys/search", searchProducts)
		//wechat
		api.GET("/wechat/auth", wechatAuth)
		api.GET("/wechat/token", wechatToken)
		api.GET("/wechat/qrcode", wechatQrCode)
		api.GET("/wechat/temp_qrcode", wechatTempQrCode)
		api.POST("/wechat/status", wechatQrCodeStatus)
		api.GET("/wechat/callback", wechatCallback)
		api.POST("/wechat/callback", wechatCallback)

		//token
		api.GET("/token", GetToken)
		api.POST("/token", VerifyToken)
	}
	r.Static("../webpage", "webpage")
}

func init() {
	// u, err := qrCodePic()
	// if err != nil {
	// 	log.Fatalf("qr code pic error %s", err.Error())
	// }
	// log.Println(u)
	os.Setenv("TRACE", "true")
	//running dir
	// if err := chdir(); err != nil {
	// 	log.Fatal(err)
	// }
	lf, err := os.OpenFile("heshi.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer lf.Close()

	if util.ShouldTrace() {
		log.SetOutput(io.MultiWriter(os.Stdout, lf))
		util.Logger = log.New(io.MultiWriter(os.Stdout, lf), "", log.LstdFlags)
	}
	log.SetFlags(log.LstdFlags)

	db, err = mysql.OpenDB()
	if db == nil && err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}

	activeConfig = config{Rate: 0.01, CreatedBy: "system", CreatedAt: time.Now().Local()}
	val, err := redisClient.FlushAll().Result()
	if err != nil {
		util.Printf("fail to flush redis db. err: %s", err.Error())
	}
	util.Printf("flushed redis db. %s", val)

	store = sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: int(30 * time.Minute), //30min
		Path:   "/",
	})
	// if err := getLatestRates(); err != nil {
	// 	log.Fatalf("init fail. err: %s;", err.Error())
	// }
}

func chdir() error {
	pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	return os.Chdir(pwd)
}
