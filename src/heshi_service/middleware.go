package main

import (
	"net/http"
	"os"
	"time"
	"util"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"https://github.com"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	})
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := c.Request.Header.Get("X-Auth-Token")
		if t == "" {
			c.AbortWithStatus(401)
		} else if util.VerfiyToken(t) {
			c.Next()
		} else {
			c.AbortWithStatus(401)
		}
	}
}

func UserSessionMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("stage") == "dev" {
			return
		}
		s := sessions.Default(c)
		if s.Get(USER_SESSION_KEY) == nil {
			c.JSON(http.StatusForbidden, "must login first")
			return
		}
		c.Set("id", s.Get(USER_SESSION_KEY))
	}
}

func AdminSessionMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("stage") == "dev" {
			return
		}
		s := sessions.Default(c)
		if s.Get(USER_SESSION_KEY) == nil {
			c.JSON(http.StatusForbidden, "must login first")
			return
		}
		if s.Get(ADMIN_KEY) == nil {
			c.JSON(http.StatusForbidden, "not authorized")
			return
		}
		c.Set("id", s.Get(USER_SESSION_KEY))
	}
}
