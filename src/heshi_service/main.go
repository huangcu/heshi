package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"db/mysql"
	"flag"
	"fmt"
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
	os.Setenv("STAGE", env)
	os.Setenv("TRACE", "true")

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
	r := gin.New()
	//if PRO
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	if os.Getenv("STAGE") != "dev" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r.Use(gin.Recovery())
	configRoute(r)
	if os.Getenv("STAGE") != "dev" {
		webServer := &http.Server{Addr: port, Handler: r, TLSConfig: config}
		return webServer.ListenAndServeTLS("server.crt", "server.key")
	}
	webServer := &http.Server{Addr: port, Handler: r}
	return webServer.ListenAndServe()
}

func configRoute(r *gin.Engine) {
	api := r.Group("/api")
	if os.Getenv("STAGE") != "dev" {
		api.Use(AuthMiddleWare())
	}

	store = sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: int(30 * time.Minute), //30min
		Path:   "/",
	})
	api.Use(sessions.Sessions("SESSIONID", store))
	//Cross-Site Request Forgery (CSRF)
	// api.Use(csrf.Middleware(csrf.Options{
	// 	Secret: "secret",
	// 	ErrorFunc: func(c *gin.Context) {
	// 		c.String(400, "CSRF token mismatch")
	// 		c.Abort()
	// 	},
	// }))
	api.Use(CORSMiddleware())
	api.Use(RequestLogger())

	jwtMiddleware := AuthenticateMiddleWare()

	{
		apiAdmin := api.Group("admin")
		apiAdmin.Use(AdminSessionMiddleWare())
		apiAdmin.Use(jwtMiddleware.MiddlewareFunc())
		{
			//admin agent user
			apiAdmin.POST("/users", newAdminAgentUser)
			apiAdmin.PATCH("/users/:id", updateAdminAgent)

			//user - by id, can be admin. agent.customer
			apiAdmin.GET("/users/:id", getUser)

			//pass ?type="user_type" for admin/agent/customer
			apiAdmin.GET("/users", getAllUsers)

			//customer, agent, admin disable
			apiAdmin.DELETE("/users/:id", disableUser)

			//supplier
			apiAdmin.POST("/suppliers", newSupplier)
			apiAdmin.GET("/suppliers", getAllSuppliers)
			apiAdmin.GET("/suppliers/:id", getSupplier)
			apiAdmin.PUT("/suppliers/:id", updateSupplier)
			apiAdmin.DELETE("/suppliers/:id", disableSupplier)

			//price setting
			apiAdmin.POST("/pricesettings", addPriceRule)
			apiAdmin.GET("/pricesettings", getAllPriceRule)
			apiAdmin.GET("/pricesettings/:id", getPriceRule)
			apiAdmin.PUT("/pricesettings/:id", updatePriceRule)
			apiAdmin.DELETE("/pricesettings/:id", disablePriceRule)

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

		}
		//customer
		api.POST("/users", newUser)
		api.POST("/login", jwtMiddleware.LoginHandler)
		// api.POST("/login", userLogin)
		api.PATCH("/users", UserSessionMiddleWare(), jwtMiddleware.MiddlewareFunc(), updateUser)
		api.GET("/users", UserSessionMiddleWare(), getUser)
		api.POST("/logout", userLogout)
		api.GET("/users/:id/contactinfo", UserSessionMiddleWare(), jwtMiddleware.MiddlewareFunc(), agentContactInfo)

		//action- > add, delete
		api.POST("/shoppingList/:action", UserSessionMiddleWare(), jwtMiddleware.MiddlewareFunc(), toShoppingList)

		//products
		api.GET("/products", getAllProducts)
		api.GET("/products/diamonds", getAllDiamonds)
		api.GET("/products/diamonds/:id", getDiamond)
		api.GET("/products/jewelrys", getAllJewelrys)
		api.GET("/products/jewelrys/:id", getJewelry)
		api.GET("/products/small_diamonds", getAllSmallDiamonds)
		api.GET("/products/small_diamonds/:id", getSmallDiamond)
		api.GET("/products/gems", getAllGems)
		api.GET("/products/gems/:id", getGem)
		api.POST("/products/diamonds", newDiamond)
		api.POST("/products/small_diamonds", newSmallDiamond)
		api.POST("/products/jewelrys", newJewelry)

		//product search - diamond or jewelry by id or name
		api.POST("/products/search/:category", searchProducts)
		//product filter by fields value
		api.POST("/products/filter/:category", filterProducts)
		// api.POST("/products/diamonds/search", searchProducts)
		// api.POST("/products/jewelrys/search", searchProducts)

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
	//running dir
	// if err := chdir(); err != nil {
	// 	log.Fatal(err)
	// }
	var err error
	db, err = mysql.OpenDB()
	if db == nil && err != nil {
		util.Println(err.Error())
		os.Exit(1)
	}

	activeConfig = config{Rate: 0.01, CreatedBy: "system", CreatedAt: time.Now().Local()}
	val, err := redisClient.FlushAll().Result()
	if err != nil {
		fmt.Printf("fail to flush redis db. err: %s", err.Error())
	}
	fmt.Printf("flushed redis db. %s \n", val)

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
