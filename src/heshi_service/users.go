package main

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username       string `json:"username" valid:"length(6|40),matches(^[a-zA-Z]*$),optional"`
	Cellphone      string `json:"cellphone" valid:"matches(^[0-9]*$),optional"`
	Email          string `json:"email" valid:"email,optional"`
	Password       string `json:"password" valid:"length(8|20),matches(^[a-zA-Z0-9_!@#$%^&.?()-=+]*$),required"`
	UserType       string `json:"user_type" valid:"in(admin|agent|customer),required"`
	RealName       string `json:"real_time" valid:"-"`
	WechatID       string `json:"wechat_id" valid:"-"`
	WechatName     string `json:"wechat_name" valid:"-"`
	WechatQR       string `json:"wechat_qr" valid:"-"`
	Address        string `json:"address" valid:"-"`
	AdditionalInfo string `json:"additional_info" valid:"-"`
	Icon           string `json:"icon" valid:"-"`
	// CreatedAt      time.Time `json:"created_at"`
	// UpdatedAt      time.Time `json:"updated_at"`
}

func newUser(c *gin.Context) {
	nu := User{
		Username:       c.PostForm("username"),
		Cellphone:      c.PostForm("cellphone"),
		Email:          c.PostForm("email"),
		Password:       c.PostForm("password"),
		UserType:       c.PostForm("user_type"),
		RealName:       c.PostForm("real_time"),
		WechatID:       c.PostForm("wechat_id"),
		WechatName:     c.PostForm("wechat_name"),
		WechatQR:       c.PostForm("wechat_qr"),
		Address:        c.PostForm("address"),
		AdditionalInfo: c.PostForm("additional_info"),
		Icon:           c.PostForm("icon"),
	}

	if _, errs := govalidator.ValidateStruct(nu); errs != nil {
		c.String(http.StatusOK, "Hello %s", errs.Error())
		return
	}
	fmt.Println(nu)
	c.String(http.StatusOK, nu.Username)

	// q := `INSERT INTO users (username, cellphone, email, password)
	// 	VALUES (?,?,?,?)`
	// if _, err := dbExec(q, c.PostForm("username"), c.PostForm("cellphone"), c.PostForm("email"),
	// 	c.PostForm("password")); err != nil {
	// 	c.String(status, errors.GetMessage(err))
	// }
}

func updateUser(c *gin.Context) {
	var usertype int
	switch usertype {
	case 0:
	case 1:
	case 2:
	}
}

func getUser(c *gin.Context) {

}

func getAllUsers(c *gin.Context) {

}
