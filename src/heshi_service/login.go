package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"jwt"
	"log"
	"net/http"
	"os"
	"time"
	"util"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// LoginUser ...
type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserType string `json:"user_type"`
}

func userLogin(c *gin.Context) {
	loginUser := LoginUser{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		UserType: c.PostForm("user_type"),
	}

	q := fmt.Sprintf(`SELECT id, password, user_type, status FROM users where username='%s' or cellphone='%s' or email='%s'`,
		loginUser.Username, loginUser.Username, loginUser.Username)

	var id, password, userType, status string
	if err := dbQueryRow(q).Scan(&id, &password, &userType, &status); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, errorLoginUserNamePassword)
			return
		}
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}

	if status != "ACTIVE" {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s is not an active user!", loginUser.Username))
		return
	}

	if loginUser.UserType != "" && loginUser.UserType != userType {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s is not %s", loginUser.Username, loginUser.UserType))
		return
	}

	if !util.IsPassOK(loginUser.Password, password) {
		c.JSON(http.StatusBadRequest, errorLoginUserNamePassword)
		return
	}

	userProfile, user, err := getUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  userProfile,
		"user":   user,
		"expire": time.Now().Add(time.Minute * 30).Format(time.RFC3339),
	})
}

func userLogout(c *gin.Context) {
	// s := sessions.Default(c)
	// fmt.Println(s.Get(userSessionKey).(string))
	// s.Delete(userSessionKey)
	// s.Delete(adminKey)
	// s.Delete(agentKey)
	// s.Save()
	fmt.Println("logout")
	if os.Getenv("STAGE") != "dev" {
		token := jwt.GetToken(c)
		if token != "" {
			redisClient.Del(token)
		}
	}
	c.JSON(http.StatusOK, "User logout!")
}

func isValidCacheToken(token string) bool {
	if _, err := redisClient.Get(token).Result(); err != nil {
		if err != redis.Nil {
			log.Println("redis error: " + err.Error())
		}
		return false
	}
	redisClient.Expire(token, time.Minute*30)
	return true
}
