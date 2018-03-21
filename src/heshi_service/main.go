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
	"runtime"
	"strings"
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

	port := ":8080"
	if os.Getenv("STAGE") != "dev" {
		port = ":8443"
	}

	log.Fatal(startWebServer(port))
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
	//CORS
	r.Use(CORSMiddleware())
	//session
	store = sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: int(30 * time.Minute), //30min
		Path:   "/",
	})
	r.Use(sessions.Sessions("SESSIONID", store))
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
		// auth - access service api
		api.Use(AuthMiddleWare())

		// Cross-Site Request Forgery (CSRF)
		// api.Use(csrf.Middleware(csrf.Options{
		// 	Secret: "secret",
		// 	ErrorFunc: func(c *gin.Context) {
		// 		c.String(400, "CSRF token mismatch")
		// 		c.Abort()
		// 	},
		// }))
	}
	//access api log
	api.Use(RequestLogger())

	apiCustomer := api.Group("customer")
	apiAdmin := api.Group("admin")
	apiWechat := api.Group("wechat")
	apiCustomer.Use(UserSessionMiddleWare())
	apiAdmin.Use(AdminSessionMiddleWare())
	//TODO wechat - > admin and customer
	apiWechat.Use(AdminSessionMiddleWare())

	if os.Getenv("STAGE") == "dev" {
		api.POST("/login", userLogin)
	} else {
		//jwt authentication(user login)
		jwtMiddleware := AuthenticateMiddleWare()
		apiCustomer.Use(jwtMiddleware.MiddlewareFunc())
		apiAdmin.Use(jwtMiddleware.MiddlewareFunc())
		api.POST("/login", jwtMiddleware.LoginHandler)
		api.GET("/refresh/token", jwtMiddleware.RefreshHandler)
	}

	{
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
			apiAdmin.GET("/configs/rate", getRateConfig)
			apiAdmin.GET("/configs/allrate", getAllRateConfigs)
			apiAdmin.POST("/configs/rate", newRateConfig)
			//account level, discount config
			apiAdmin.GET("/configs/level/:id", getLevelConfig)
			apiAdmin.PUT("/configs/level/:id", updateLevelConfig)
			apiAdmin.GET("/configs/level", getAllLevelConfigs)
			apiAdmin.POST("/configs/level", newLevelConfig)

			//products with customize header
			apiAdmin.POST("/upload", uploadAndGetFileHeaders)
			apiAdmin.POST("/process/diamond", processDiamonds)

			//upload products by csv file
			apiAdmin.POST("/products/upload", uploadAndProcessProducts)

			//manage products
			apiAdmin.POST("/products/diamonds", newDiamond)
			apiAdmin.PUT("/products/diamonds/:id", updateDiamond)
			apiAdmin.POST("/products/small_diamonds", newSmallDiamond)
			apiAdmin.POST("/products/jewelrys", newJewelry)
			apiAdmin.PUT("/products/jewelrys/:id", updateJewelry)
			apiAdmin.POST("/products/gems", newGems)
			apiAdmin.PUT("/products/gems/:id", updateGems)

			//manage orders
			apiAdmin.PUT("/orders/:id", updateOrder)
			apiAdmin.POST("/orders/:id", getOrderDetail)
			apiAdmin.POST("/transactions/detail/:id", getTransactionDetail)
			apiAdmin.POST("/transactions/all/:id", getAllTransctionsOfUser)
			apiAdmin.POST("/transactions/all", getAllTransctions)

			//wechat kf manage
			apiAdmin.POST("/wechat/kf", addKfAccount)
			apiAdmin.POST("/wechat/menu", createMenu)
		}
		//customer
		api.POST("/users", newUser)
		{
			apiCustomer.GET("/users", getUser)
			apiCustomer.PATCH("/users", updateUser)
			apiCustomer.GET("/users/:id/contactinfo", agentContactInfo)

			//ORDER
			apiCustomer.POST("/orders", createOrder)
			apiCustomer.POST("/orders/:id", getOrderDetail)
			apiCustomer.POST("/transactions/:id", getTransactionDetail)
			apiCustomer.POST("/transactions", getAllTransctionsOfUser)

			//action- > add, delete
			apiCustomer.POST("/shoppingList/:action", toShoppingList)

			apiCustomer.POST("/logout", userLogout)
			apiCustomer.POST("/password/change/:id", changePassword)
		}

		// exchange rate
		api.GET("/exchangerate", getCurrencyRate)

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
		//product search - diamond or jewelry by id or name
		api.POST("/products/search/:category", searchProducts)
		//product filter by fields value
		api.POST("/products/filter/:category", filterProducts)
		api.POST("/product/customize/:action", customizeProduct)

		//wechat - todo
		apiWechat.GET("/auth", wechatAuth)
		apiWechat.GET("/token", wechatToken)
		apiWechat.GET("/qrcode", wechatQrCode)

		api.GET("/wechat/temp_qrcode", wechatTempQrCode)
		api.GET("/wechat/status", wechatQrCodeStatus)
		api.GET("/wechat/callback", wechatCallback)
		api.POST("/wechat/callback", wechatCallback)

		api.Static("/image", ".image")
		api.Static("/video", ".video")
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
	if strings.ToUpper(runtime.GOOS) != "WINDOWS" {
		fmt.Println("OS: " + runtime.GOOS)
		val, err := redisClient.FlushAll().Result()
		if err != nil {
			log.Fatalf("fail to flush redis db. err: %s", err.Error())
		}
		fmt.Printf("flushed redis db. %s \n", val)
	}
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
