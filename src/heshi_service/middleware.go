package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"heshi/errors"
	"io"
	"io/ioutil"
	"jwt"
	"net/http"
	"os"
	"time"
	"util"

	detect "github.com/erizocosmico/detect.git"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	config := cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	var allowHeaders []string
	allowHeaders = append(allowHeaders, []string{"Origin", "Content-Type", "Content-Length"}...)
	allowHeaders = append(allowHeaders, []string{"Accept", "Accept-Encoding", "Accept-Language"}...)
	allowHeaders = append(allowHeaders, []string{"Authorization", "Cookie", "Cache-Control", "X-Auth-Token"}...)
	allowHeaders = append(allowHeaders, []string{"Cache-Control", "Connection", "User-Agent"}...)
	config.AddAllowHeaders(allowHeaders...)
	return cors.New(config)
	// cors.DefaultConfig()
}

func AuthenticateMiddleWare() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:            "HESHI",
		Key:              []byte("secret key"),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour,
		Authenticator:    userLogin1,
		TokenLookup:      "header:Authorization",
		TokenHeadName:    "Bearer",
		PrivKeyFile:      "token.key",
		PubKeyFile:       "token_pk.pem",
		SigningAlgorithm: "RS512",
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}
}

func userLogin1(username string, password1 string, c *gin.Context) (string, bool) {
	q := fmt.Sprintf(`SELECT id, password, user_type FROM users where username='%s' or cellphone='%s' or email='%s'`,
		username, username, username)

	var id, password, userType string
	if err := dbQueryRow(q).Scan(&id, &password, &userType); err != nil {
		if err == sql.ErrNoRows {
			return "", false
		}
		return "", false
	}

	if !util.IsPassOK(password1, password) {
		return vemsgLoginErrorUserName.Message, false
	}

	userProfile, err := getUserByID(id)
	if err != nil {
		return errors.GetMessage(err), false
	}
	s := sessions.Default(c)
	s.Set(USER_SESSION_KEY, id)

	if userType == ADMIN {
		s.Set(ADMIN_KEY, id)
	}
	s.Save()
	return userProfile, true
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := c.Request.Header.Get("X-Auth-Token")
		if t == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Not allowed without auth token")
		} else if util.VerfiyToken(t) {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Auth token is not correct")
		}
	}
}

func UserSessionMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("STAGE") == "dev" {
			c.Set("id", "system_dev_user")
			c.Next()
			return
		}
		s := sessions.Default(c)
		if s.Get(USER_SESSION_KEY) == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Must login first")
			return
		}
		c.Set("id", s.Get(USER_SESSION_KEY))
		c.Next()
	}
}

func AdminSessionMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("STAGE") == "dev" {
			c.Set("id", "system_dev_admin")
			c.Next()
			return
		}
		s := sessions.Default(c)
		if s.Get(USER_SESSION_KEY) == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Must login first")
			return
		}
		if s.Get(ADMIN_KEY) == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, "Login User is not admin")
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
		platform := detect.Platform(c.Request.Header.Get("User-Agent")).String()
		util.Printf("=======GOT REQUEST - METHOD: %s; FROM: %s; URL %s=====", c.Request.Method, platform, c.Request.URL)
		util.Printf("=======REQUEST HEADER: %v========", c.Request.Header)
		util.Printf("=======REQUEST BODY: %s========", readBody(rdr1)) // Print request body
		s := sessions.Default(c)
		user := s.Get(USER_SESSION_KEY)
		if user == nil {
			user = "guest"
		}
		if err := userUsingRecord(c.Request.URL.Path, user.(string), platform); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
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
