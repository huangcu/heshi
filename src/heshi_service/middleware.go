package main

import (
	"bytes"
	"io"
	"io/ioutil"
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
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("id", s.Get(USER_SESSION_KEY))
		c.Next()
	}
}

func AdminSessionMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("stage") == "dev" {
			return
		}
		s := sessions.Default(c)
		if s.Get(USER_SESSION_KEY) == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if s.Get(ADMIN_KEY) == nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Set("id", s.Get(USER_SESSION_KEY))
		c.Next()
	}
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		util.Println(c.Request.URL)
		util.Println(readBody(rdr1)) // Print request body

		c.Request.Body = rdr2
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
