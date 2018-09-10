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
	"log"
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

// give it 2 mins chance window to refresh token after current token expire
func jwtMiddleWare() *jwt.GinJWTMiddleware {
	mw, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "HESHI",
		Key:              []byte("secret key"),
		Timeout:          30 * time.Minute,
		MaxRefresh:       32 * time.Minute,
		Authenticator:    jwtAuthenticator,
		Authorizator:     jwtAuthorizator,
		TokenLookup:      "header:Authorization",
		TokenHeadName:    "Bearer",
		PrivKeyFile:      "token.key",
		PubKeyFile:       "token_pk.pem",
		SigningAlgorithm: "RS512",
		IdentityKey:      "userprofile",
		IdentityHandler:  identityHandler,
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}
	return mw
}

func jwtAuthenticator(c *gin.Context) (interface{}, error) {
	username := c.PostForm("username")
	password1 := c.PostForm("password")
	q := fmt.Sprintf(`SELECT id, password, user_type, status FROM users where username='%s' or cellphone='%s' or email='%s'`,
		username, username, username)
	usertype := c.PostForm("user_type")
	var id, password, userType, status string
	if err := dbQueryRow(q).Scan(&id, &password, &userType, &status); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errorLoginUserNamePassword)
		}
		return nil, errors.New("System error, please try again later")
	}

	if status != "ACTIVE" {
		return nil, errors.Newf("%s is not an active user!", username)
	}

	if usertype != "" && usertype != userType {
		return nil, errors.Newf("%s is not %s", username, usertype)
	}

	if !util.IsPassOK(password1, password) {
		return nil, errors.New(errorLoginUserNamePassword)
	}

	userProfile, err := getUserByID(id)
	if err != nil {
		return nil, err
	}
	s := sessions.Default(c)
	s.Set(userSessionKey, id)
	// c.SetCookie(USER_SESSION_KEY, id, 10, "/", "localhost", true, false)

	if userType == ADMIN {
		s.Set(adminKey, id)
		// c.SetCookie(ADMIN_KEY, id, 10, "/", "localhost", true, false)
	}
	if userType == AGENT {
		s.Set(agentKey, id)
		// c.SetCookie(AGENT_KEY, id, 10, "/", "localhost", true, false)
	}
	s.Save()
	return userProfile, nil
}

func jwtAuthorizator(data interface{}, c *gin.Context) bool {
	var user User
	userprofile := data.(string)
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

	if user.UserType == CUSTOMER {
		if strings.HasPrefix(c.Request.RequestURI, "/api/admin") || strings.HasPrefix(c.Request.RequestURI, "/api/agent") {
			return false
		}
	}
	return true
}

func identityHandler(c *gin.Context) interface{} {
	if isTokenInBlackList(jwt.GetToken(c)) {
		return ""
	}
	claims := jwt.ExtractClaims(c)
	ip := claims["ip"]
	userAgent := claims["user-agent"]

	// 	ipFormat := regexp.MustCompile(`^(\d+).(\d+).(\d+).(\d+):\d+$`)
	// previousIP := ipFormat.FindStringSubmatch(ip)
	// currentIP := ipFormat.FindStringSubmatch(request.RemoteAddr)
	// if len(previousIP) == 5 && len(currentIP) == 5 && AcceptRemoteIP <= 4 {
	// 	for i := 1; i < AcceptRemoteIP; i++ {
	// 		if previousIP[i] != currentIP[i] {
	// 			valid = false
	// 			break
	// 		}
	// 	}
	// }
	// For now, request must from same ip and same agent
	if ip != util.GetRequestIP(c.Request) && userAgent != c.Request.Header.Get("User-Agent") {
		return ""
	}
	return claims["userprofile"]
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

// func sessionMiddleWare() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		s := sessions.Default(c)
// 		if s.Get(userSessionKey) != nil {
// 			c.Set("id", s.Get(userSessionKey))
// 		} else {
// 			c.Set("id", "guest:"+c.Request.RemoteAddr)
// 		}
// 		c.Next()
// 	}
// }

// func userSessionMiddleWare() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Set("usertype", CUSTOMER)
// 		if os.Getenv("STAGE") == "dev" {
// 			c.Set("id", "system_dev_user")
// 			c.Next()
// 			return
// 		}
// 		s := sessions.Default(c)
// 		if s.Get(userSessionKey) == nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, "Must login first")
// 			return
// 		}
// 		c.Set("id", s.Get(userSessionKey))
// 		c.Next()
// 	}
// }

// func adminSessionMiddleWare() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Set("usertype", ADMIN)
// 		if os.Getenv("STAGE") == "dev" {
// 			c.Set("id", "system_dev_admin")
// 			c.Next()
// 			return
// 		}
// 		s := sessions.Default(c)
// 		if s.Get(userSessionKey) == nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, "Must login first")
// 			return
// 		}
// 		if s.Get(adminKey) == nil {
// 			c.AbortWithStatusJSON(http.StatusForbidden, "Login User is not admin")
// 			return
// 		}
// 		c.Set("id", s.Get(userSessionKey))
// 		c.Next()
// 	}
// }

// func agentSessionMiddleWare() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Set("usertype", AGENT)
// 		if os.Getenv("STAGE") == "dev" {
// 			c.Set("id", "system_dev_admin")
// 			c.Next()
// 			return
// 		}
// 		s := sessions.Default(c)
// 		if s.Get(userSessionKey) == nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, "Must login first")
// 			return
// 		}
// 		if s.Get(agentKey) == nil {
// 			c.AbortWithStatusJSON(http.StatusForbidden, "Login User is not a agent")
// 			return
// 		}
// 		c.Set("id", s.Get(userSessionKey))
// 		c.Next()
// 	}
// }
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
		user := s.Get(userSessionKey)
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
