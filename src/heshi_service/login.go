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

	q := fmt.Sprintf(`SELECT id, password FROM users where username='%s' or cellphone='%s' or email='%s'`,
		loginUser.Username, loginUser.Username, loginUser.Username)

	var id, password string
	if err := dbQueryRow(q).Scan(&id, &password); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, VEMSG_LOGIN_ERROR_USERNAME)
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if !util.IsPassOK(loginUser.Password, password) {
		c.JSON(http.StatusOK, VEMSG_LOGIN_ERROR_USERNAME)
		return
	}

	s := sessions.Default(c)
	s.Set(USER_SESSION_KEY, id)
	s.Save()

	c.JSON(http.StatusOK, "session")
}

func userLogout(c *gin.Context) {
	s := sessions.Default(c)
	s.Delete(USER_SESSION_KEY)
	c.JSON(http.StatusOK, "User logout!")
}
