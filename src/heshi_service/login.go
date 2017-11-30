package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"util"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginUser struct {
	ID        string `json:"id"`
	Username  string `json:"username" valid:"length(6|40),matches(^[a-zA-Z0-9]*$),optional"`
	Cellphone string `json:"cellphone" valid:"matches(^[0-9]*$),optional"`
	Email     string `json:"email" valid:"email,optional"`
	Password  string `json:"password" valid:"length(8|20),matches(^[a-zA-Z0-9_!@#$%^&.?()-=+]*$),required"`
	UserType  string `json:"user_type" valid:"in(admin|agent|customer),required"`
}

var store *sessions.CookieStore

func userLogin(c *gin.Context) {
	nu := LoginUser{
		Username:  c.PostForm("username"),
		Cellphone: c.PostForm("cellphone"),
		Email:     c.PostForm("email"),
		Password:  c.PostForm("password"),
		UserType:  c.PostForm("user_type"),
	}

	if _, errs := govalidator.ValidateStruct(nu); errs != nil {
		c.String(http.StatusOK, errs.Error())
		return
	}
	var q string
	if nu.Username != "" {
		q = fmt.Sprintf(`SELECT password FROM users where username=%s`, nu.Username)
	}
	if nu.Username != "" {
		q = fmt.Sprintf(`SELECT password FROM users where cellphone=%s`, nu.Cellphone)
	}
	if nu.Username != "" {
		q = fmt.Sprintf(`SELECT password FROM users where email=%s`, nu.Email)
	}

	var password string
	if err := dbQueryRow(q).Scan(&password); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, "login info not correct, wrong user or password")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if !util.IsPassOK(nu.Password, password) {
		c.JSON(http.StatusOK, "login info not correct, wrong user or password")
		return
	}
	session := sessions.Default(c)
	v := session.Get("mysession")
	if v == nil {
		session.Set("mysession", "count")
	}
	session.Save()

	c.JSON(http.StatusOK, session)
	// var usertype int
	// switch usertype {
	// case 0:
	// case 1:
	// case 2:
	// }

}
