package main

import (
	"net/http"
	"util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		// if c.Req.Method == "OPTIONS" {
		// 	fmt.Println("options")
		// 	c.Abort(200)
		// 	return
		// }
		c.Next()
	}
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
