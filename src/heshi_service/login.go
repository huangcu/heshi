package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func userLogin(c *gin.Context) {
	loginUser := LoginUser{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}

	q := fmt.Sprintf(`SELECT id, password, user_type FROM users where username='%s' or cellphone='%s' or email='%s'`,
		loginUser.Username, loginUser.Username, loginUser.Username)

	var id, password, userType string
	if err := dbQueryRow(q).Scan(&id, &password, &userType); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, vemsgLoginErrorUserName)
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if !util.IsPassOK(loginUser.Password, password) {
		c.JSON(http.StatusOK, vemsgLoginErrorUserName)
		return
	}

	s := sessions.Default(c)
	s.Set(USER_SESSION_KEY, id)
	if userType == ADMIN {
		s.Set(ADMIN_KEY, id)
	}
	s.Save()

	c.JSON(http.StatusOK, "success")
}

func userLogout(c *gin.Context) {
	s := sessions.Default(c)
	s.Delete(USER_SESSION_KEY)
	c.JSON(http.StatusOK, "User logout!")
}
