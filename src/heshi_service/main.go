package main

import (
	"db/mysql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := mysql.OpenDB()

}

func startWebServer(port string) error {
	r := gin.New()
	if os.Getenv("stage") == "dev" {
		gin.SetMode(gin.DebugMode)
		r.Use(gin.Logger())
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Use(gin.Recovery())
	configRoute(r)
	webServer := &http.Server{Addr: port, Handler: r}
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
