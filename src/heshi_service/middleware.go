package main

import (
	"util"

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
