package main

import (
	"fmt"
	"heshi/errors"
	"net/http"
	"util"

	"github.com/asaskevich/govalidator"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username       string `json:"username" valid:"length(6|40),matches(^[a-zA-Z0-9]*$),optional"`
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
		RealName:       c.PostForm("real_name"),
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
	if nu.Username == "" && nu.Cellphone == "" && nu.Email == "" {
		c.String(http.StatusBadRequest, "username, cellphone, email mustn't be empty")
		return
	}
	q := `INSERT INTO users (password,user_type`
	v := `VALUES (?,?`

	p := []string{util.Encrypt(nu.Password), nu.UserType}
	if nu.Username != "" {
		q = fmt.Sprintf("%s, username", q)
		v = fmt.Sprintf("%s, ?", v)
		p = append(p, nu.Username)
	}
	if nu.Cellphone != "" {
		q = fmt.Sprintf("%s, cellphone", q)
		v = fmt.Sprintf("%s, ?", v)
		p = append(p, nu.Cellphone)
	}
	if nu.Email != "" {
		q = fmt.Sprintf("%s, email", q)
		v = fmt.Sprintf("%s, ?", v)
		p = append(p, nu.Email)
	}
	if nu.RealName != "" {
		q = fmt.Sprintf("%s, real_name", q)
		v = fmt.Sprintf("%s, ?", v)
	}
	if nu.WechatID != "" {
		q = fmt.Sprintf("%s, wechat_id", q)
		v = fmt.Sprintf("%s, ?", v)
	}
	if nu.WechatName != "" {
		q = fmt.Sprintf("%s, wechat_name", q)
		v = fmt.Sprintf("%s, ?", v)
	}
	if nu.WechatQR != "" {
		q = fmt.Sprintf("%s, wechat_qr", q)
		v = fmt.Sprintf("%s, ?", v)
	}
	if nu.Address != "" {
		q = fmt.Sprintf("%s, address", q)
		v = fmt.Sprintf("%s, ?", v)
	}
	if nu.AdditionalInfo != "" {
		q = fmt.Sprintf("%s, additional_info", q)
		v = fmt.Sprintf("%s, ?", v)
	}
	if nu.Icon != "" {
		q = fmt.Sprintf("%s, icon", q)
		v = fmt.Sprintf("%s, ?", v)
	}
	q = fmt.Sprintf("%s) %s)", q, v)
	pv := make([]interface{}, len(p))
	for i, v := range p {
		pv[i] = v
	}
	if _, err := dbExec(q, pv...); err != nil {
		c.String(http.StatusBadRequest, errors.GetMessage(err))
	}
	c.String(http.StatusOK, nu.Username)
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
