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

	"github.com/gin-gonic/gin"
)

var (
	db                  *sql.DB
	ctx                 context.Context
	cancelFn            context.CancelFunc
	serverIsInterrupted bool
)

func main() {
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

	log.Fatal(startWebServer(":8443"))
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
	r.Use(gin.Recovery())
	configRoute(r)
	// webServer := &http.Server{Addr: port, Handler: r, TLSConfig: config}
	webServer := &http.Server{Addr: port, Handler: r}
	// return webServer.ListenAndServe()
	return webServer.ListenAndServeTLS("server.crt", "server.key")
}

func configRoute(r *gin.Engine) {
	api := r.Group("/api")
	{
		apiAdmin := api.Group("admin")
		{
			//admin user -
			apiAdmin.POST("/users", newUser)
			apiAdmin.GET("/users/:id", getUser)
			apiAdmin.GET("/users", getAllUsers)
			apiAdmin.PATCH("/users/:id", updateUser)
			apiAdmin.DELETE("/users/:id", removeUser)

			//currency rate
			apiAdmin.POST("/exchangerate", currencyRateReqValidator(newCurrencyRate))
		}
		//agent, customer
		api.POST("/users", newUser)
		api.PATCH("/users/:id", updateUser)
		api.GET("/users/:id", getUser)
		// store := sessions.NewCookieStore([]byte("secret"))
		// r.Use(sessions.Sessions("mysession", store))
		// store := sessions.NewMemcacheStore(memcache.New("localhost:11211"), "", []byte("secret"))
		// r.Use(sessions.Sessions("mysession", store))
		api.POST("/login", userLogin)
		// api.POST("/login", userLogin, sessions.Sessions("mysession", store))

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

		//get the latest in db
		api.GET("/exchangerate", getCurrencyRate)
		api.GET("/wechat/auth", wechatAuth)
		api.GET("/wechat/token", wechatToken)
	}
	api.Static("../webpage", "webpage")
}
