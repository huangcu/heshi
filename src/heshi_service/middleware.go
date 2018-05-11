package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"heshi/errors"
	"io"
	"io/ioutil"
	"jwt"
	"net/http"
	"os"
	"strings"
	"time"
	"util"

	detect "github.com/erizocosmico/detect.git"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func cORSMiddleware() gin.HandlerFunc {
	config := cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"http://localhost:8080"},
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
	if os.Getenv("STAGE") != "dev" {
		config.AllowAllOrigins = false
		config.AllowOrigins = []string{"http://localhost:8080", "https://localhost:8443"}
	}
	return cors.New(config)
	// cors.DefaultConfig()
}

func authenticateMiddleWare() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:            "HESHI",
		Key:              []byte("secret key"),
		Timeout:          30 * time.Minute,
		MaxRefresh:       30 * time.Minute,
		Authenticator:    jwtAuthenticator,
		Authorizator:     jwtAuthorizator,
		TokenLookup:      "header:Authorization",
		TokenHeadName:    "Bearer",
		PrivKeyFile:      "token.key",
		PubKeyFile:       "token_pk.pem",
		SigningAlgorithm: "RS512",
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}
}

func jwtAuthenticator(username, password1 string, c *gin.Context) (string, bool) {
	q := fmt.Sprintf(`SELECT id, password, user_type, status FROM users where username='%s' or cellphone='%s' or email='%s'`,
		username, username, username)
	usertype := c.PostForm("user_type")
	var id, password, userType, status string
	if err := dbQueryRow(q).Scan(&id, &password, &userType, &status); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Sprintf("%s not exists", username), false
		}
		return "System error, please try again later", false
	}

	if status != "ACTIVE" {
		return fmt.Sprintf("%s is not an active user!", username), false
	}

	if usertype != "" && usertype != userType {
		return fmt.Sprintf("%s is not %s", username, usertype), false
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
	// c.SetCookie(USER_SESSION_KEY, id, 10, "/", "localhost", true, false)

	if userType == ADMIN {
		s.Set(ADMIN_KEY, id)
		// c.SetCookie(ADMIN_KEY, id, 10, "/", "localhost", true, false)
	}
	if userType == AGENT {
		s.Set(AGENT_KEY, id)
		// c.SetCookie(AGENT_KEY, id, 10, "/", "localhost", true, false)
	}
	s.Save()
	return userProfile, true
}

func jwtAuthorizator(userID string, c *gin.Context) bool {
	userprofile := c.MustGet("userID").(string)
	var user User
	if err := json.Unmarshal([]byte(userprofile), &user); err != nil {
		return false
	}
	c.Set("id", user.ID)
	c.Set("user", user)
	if strings.HasPrefix(c.Request.RequestURI, "/api/admin") && user.UserType == ADMIN {
		return true
	}
	if strings.HasPrefix(c.Request.RequestURI, "/api/agent") && user.UserType == AGENT {
		return true
	}
	return false
}

func authMiddleWare() gin.HandlerFunc {
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

func sessionMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		if s.Get(USER_SESSION_KEY) != nil {
			c.Set("id", s.Get(USER_SESSION_KEY))
		} else {
			c.Set("id", "guest:"+c.Request.RemoteAddr)
		}
		c.Next()
	}
}

func userSessionMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("usertype", CUSTOMER)
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

func adminSessionMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("usertype", ADMIN)
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

func agentSessionMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("usertype", AGENT)
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
		if s.Get(AGENT_KEY) == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, "Login User is not a agent")
			return
		}
		c.Set("id", s.Get(USER_SESSION_KEY))
		c.Next()
	}
}
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.
		platform := detect.Platform(c.Request.Header.Get("User-Agent")).String()
		util.Tracef("=======GOT REQUEST - METHOD: %s; FROM: %s; URL %s=====", c.Request.Method, platform, c.Request.URL)
		util.Tracef("=======REQUEST HEADER: %v========", c.Request.Header)
		util.Tracef("=======REQUEST BODY: %s========", readBody(rdr1)) // Print request body
		s := sessions.Default(c)
		user := s.Get(USER_SESSION_KEY)
		if user == nil {
			user = "guest"
		}
		if err := userUsingRecord(c.Request.URL.Path, user.(string), platform, c.Request.RemoteAddr); err != nil {
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
	if len(s) > 2000 {
		return string(s[0:2000])
	}
	return s
}
