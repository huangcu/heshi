package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID             string `json:"id"`
	Username       string `json:"username" valid:"length(6|40),matches(^[a-zA-Z0-9]*$),optional"`
	Cellphone      string `json:"cellphone" valid:"matches(^[0-9]*$),optional"`
	Email          string `json:"email" valid:"email,optional"`
	Password       string `json:"password" valid:"length(8|20),matches(^[a-zA-Z0-9_!@#$%^&.?()-=+]*$),required"`
	UserType       string `json:"user_type" valid:"in(admin|agent|customer),required"`
	RealName       string `json:"real_name" valid:"-"`
	WechatID       string `json:"wechat_id" valid:"-"`
	WechatName     string `json:"wechat_name" valid:"-"`
	WechatQR       string `json:"wechat_qr" valid:"-"`
	Address        string `json:"address" valid:"-"`
	AdditionalInfo string `json:"additional_info" valid:"-"`
	Icon           string `json:"icon" valid:"-"`
	// CreatedAt      time.Time `json:"created_at"`
	// UpdatedAt      time.Time `json:"updated_at"`
}

type Admin struct {
	UserInfo   User   `json:"user"`
	Level      int    `json:"level"`
	WechatKefu string `json:"wechat_kefu"`
}

type Agent struct {
	UserInfo User    `json:"user"`
	Level    int     `json:"level"`
	Discount float64 `json:"discount"`
}

func newUser(c *gin.Context) {
	nu := User{
		ID:             uuid.NewV4().String(),
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

	if vemsg := preValidateNewUser(nu); vemsg != "" {
		c.String(http.StatusOK, vemsg)
		return
	}
	if nu.Username == "" {
		var count int
		q := "SELECT count(*) FROM users"
		if err := dbQueryRow(q).Scan(&count); err != nil {
			c.String(http.StatusOK, errors.GetMessage(err))
			return
		}
		nu.Username = fmt.Sprintf("heshi_%d%d", rand.Intn(3), count)
	}

	var q string
	var err error
	if err, q = composeUserQuery(nu); err != nil {
		c.String(http.StatusOK, errors.GetMessage(err))
		return
	}
	fmt.Println(q)
	if _, err := dbExec(q); err != nil {
		c.String(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.String(http.StatusOK, nu.ID)
}

func removeUser(c *gin.Context) {
	uid := c.Param("id")
	q := "SELECT user_type from users WHERE id=?"
	var userType string
	if err := dbQueryRow(q, uid).Scan(&userType); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, "wrong user")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	switch userType {
	case "admin":
		q = `DELETE FROM admins WHERE user_id=?`
		if _, err := dbExec(q, uid); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	case "agent":
		q = `DELETE FROM agents WHERE user_id=?`
		if _, err := dbExec(q, uid); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	default:
		q = `DELETE FROM users WHERE id=?`
		if _, err := dbExec(q, uid); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func updateUser(c *gin.Context) {
	uu := User{
		ID:             c.Param("id"),
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

	var q string
	var err error
	//TODO validate updated user info too!!!
	if err, q = composeUserUpdate(uu); err != nil {
		c.String(http.StatusOK, errors.GetMessage(err))
		return
	}

	fmt.Println(q)
	//TODO admin,agent update!!!!
	// var userType string
	// switch userType {
	// case "admin":
	// case "agent":
	// default:
	// }
	if _, err := dbExec(q); err != nil {
		c.String(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.String(http.StatusOK, uu.ID)
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	q := `SELECT username,cellphone,email,real_name, user_type,wechat_id,wechat_name,
				wechat_qr,address,additional_info,icon FROM users WHERE id=?`
	var userType, icon string
	var username, cellphone, email, realName sql.NullString
	var wechatID, wechatName, wechatQR, address, additionalInfo sql.NullString
	if err := dbQueryRow(q, id).Scan(&username, &cellphone, &email, &realName, &userType, &wechatID,
		&wechatName, &wechatQR, &address, &additionalInfo, &icon); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, User{id, username.String, cellphone.String, email.String, "", userType, realName.String,
		wechatID.String, wechatName.String, wechatQR.String, address.String, additionalInfo.String, icon})
}

func getAllUsers(c *gin.Context) {
	q := `SELECT id, username,cellphone,email,real_name, user_type,wechat_id,wechat_name,
				wechat_qr,address,additional_info,icon FROM users`
	rows, err := dbQuery(q)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var us []User
	for rows.Next() {
		var id, userType, icon string
		var username, cellphone, email, realName sql.NullString
		var wechatID, wechatName, wechatQR, address, additionalInfo sql.NullString
		if err := rows.Scan(&id, &username, &cellphone, &email, &realName, &userType, &wechatID,
			&wechatName, &wechatQR, &address, &additionalInfo, &icon); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		u := User{id, username.String, cellphone.String, email.String, "", userType, realName.String,
			wechatID.String, wechatName.String, wechatQR.String, address.String, additionalInfo.String, icon}
		us = append(us, u)
	}
	c.JSON(http.StatusOK, us)
}
