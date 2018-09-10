package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"jwt"
	"net/http"
	"time"
	"util"

	"github.com/gin-gonic/gin"
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

	userProfile, err := getUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}
	// s := sessions.Default(c)
	// s.Set(userSessionKey, id)
	// // c.SetCookie(USER_SESSION_KEY, id, 10, "/", "localhost", true, false)
	// if userType == ADMIN {
	// 	s.Set(adminKey, id)
	// 	// c.SetCookie(ADMIN_KEY, id, 10, "/", "localhost", true, false)
	// }
	// if userType == AGENT {
	// 	s.Set(agentKey, id)
	// 	// c.SetCookie(AGENT_KEY, id, 10, "/", "localhost", true, false)
	// }
	// s.Save()
	c.JSON(http.StatusOK, gin.H{
		"code":        http.StatusOK,
		"token":       "faketoken",
		"userprofile": userProfile,
	})
}

func userLogout(c *gin.Context) {
	// s := sessions.Default(c)
	// fmt.Println(s.Get(userSessionKey).(string))
	// s.Delete(userSessionKey)
	// s.Delete(adminKey)
	// s.Delete(agentKey)
	// s.Save()
	token := jwt.GetToken(c)
	if token != "" {
		claims := jwt.ExtractClaims(c)
		var remaining time.Duration
		if validity, ok := claims["exp"].(int64); ok {
			tm := time.Unix(int64(validity), 0)
			remainer := tm.Sub(time.Now())
			if remainer > 0 {
				// TODO
				redisClient.Set(token, token, remaining)
			}
		}
		redisClient.Set(token, token, remaining)
	}
	c.JSON(http.StatusOK, "User logout!")
}

func isTokenInBlackList(token string) bool {
	token, err := redisClient.Get(token).Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	//Err, not in black list - false
	// no err, in black list - true
	return token != ""
}
