package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"db/mysql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	db                  *sql.DB
	ctx                 context.Context
	cancelFn            context.CancelFunc
	serverIsInterrupted bool
)

func main() {
	var err error
	db, err = mysql.OpenDB()
	if db == nil && err != nil {
		fmt.Println(err.Error())
	}
	log.Fatal(startWebServer(":8443"))
}

func startWebServer(port string) error {
	log.SetFlags(log.Lshortfile)

	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	if err != nil {
		log.Println(err)
		return err
	}
	r := gin.New()
	if os.Getenv("stage") == "dev" {
		gin.SetMode(gin.DebugMode)
		r.Use(gin.Logger())
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Use(gin.Recovery())
	configRoute(r)
	webServer := &http.Server{Addr: port, Handler: r, TLSConfig: config}
	return webServer.ListenAndServe()
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
			apiAdmin.POST("/users/:id", updateUser)
			apiAdmin.DELETE("/users/:id", getUser)
		}
		//agent, customer
		api.POST("/users", newUser)
		api.POST("/users/:id", updateUser)
		api.GET("/users/:id", getUser)
	}
}
