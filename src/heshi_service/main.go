package main

import (
	"context"
	"database/sql"
	"db/mysql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
)
var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func main() {
	ticker := time.NewTicker(time.Hour * 8)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := getLatestRates(); err != nil {
					util.FailToGetCurrencyExchangeAlert()
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
		r.Use(AuthMiddleWare())
	}

	r.Use(sessions.Sessions("sessionid", store))
	r.Use(gin.Recovery())
	configRoute(r)
	// webServer := &http.Server{Addr: port, Handler: r, TLSConfig: config}
	webServer := &http.Server{Addr: port, Handler: r}
	return webServer.ListenAndServe()
	// return webServer.ListenAndServeTLS("server.crt", "server.key")
}

func configRoute(r *gin.Engine) {
	api := r.Group("/api")
	{
		apiAdmin := api.Group("admin")
		{
			//admin user -
			apiAdmin.POST("/users", newUser)
			apiAdmin.GET("/users/:id", AdminSessionMiddleWare(), getUser)
			apiAdmin.GET("/users", AdminSessionMiddleWare(), getAllUsers)
			apiAdmin.PATCH("/users/:id", AdminSessionMiddleWare(), updateUser)
			apiAdmin.DELETE("/users/:id", AdminSessionMiddleWare(), removeUser)

			//currency rate
			apiAdmin.GET("/exchangerate", AdminSessionMiddleWare(), getCurrencyRate)
			apiAdmin.POST("/exchangerate", AdminSessionMiddleWare(), currencyRateReqValidator(newCurrencyRate))

			//discount
			apiAdmin.GET("/discount/:id", AdminSessionMiddleWare(), getDiscount)
			apiAdmin.GET("/discount", AdminSessionMiddleWare(), getDiscounts)
			apiAdmin.POST("/discount", AdminSessionMiddleWare(), newDiscount)

			//config
			apiAdmin.GET("/config", AdminSessionMiddleWare(), getConfig)
			apiAdmin.GET("/configs", AdminSessionMiddleWare(), getConfigs)
			apiAdmin.POST("/config", AdminSessionMiddleWare(), newConfig)

			//products
			apiAdmin.POST("/upload", AdminSessionMiddleWare(), uploadProducts)
			apiAdmin.POST("/process/diamond", AdminSessionMiddleWare(), processDiamonds)
		}
		//agent, customer
		api.POST("/users", newUser)
		api.PATCH("/users/:id", UserSessionMiddleWare(), updateUser)
		api.GET("/users/:id", UserSessionMiddleWare(), getUser)
		api.POST("/login", userLogin)
		api.POST("/logout/:id", userLogout)

		api.GET("/users/:id/contactinfo", UserSessionMiddleWare(), agentContactInfo)

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

		//wechat
		api.GET("/wechat/auth", wechatAuth)
		api.GET("/wechat/token", wechatToken)
		api.GET("/wechat/qrcode", wechatQrCode)
		api.GET("/wechat/temp_qrcode", wechatTempQrCode)
		api.POST("/wechat/status", wechatQrCodeStatus)
		api.GET("/wechat/callback", wechatCallback)
		api.POST("/wechat/callback", wechatCallback)
	}
	api.Static("../webpage", "webpage")
}

func init() {
	// u, err := qrCodePic()
	// if err != nil {
	// 	log.Fatalf("qr code pic error %s", err.Error())
	// }
	// log.Println(u)
	os.Setenv("stage", "dev")
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
		fmt.Println(err.Error())
	}

	store = sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: int(30 * time.Minute), //30min
		Path:   "/",
	})
	// if err := getLatestRates(); err != nil {
	// 	log.Fatalf("init fail. err: %s;", err.Error())
	// }
	activeConfig = config{Rate: 0.01, CreatedBy: "system", CreatedAt: time.Now().Local()}
	val, err := redisClient.FlushAll().Result()
	if err != nil {
		log.Printf("fail to flush redis db. err: %s", err.Error())
	}
	log.Printf("flushed redis db. %s", val)
}
