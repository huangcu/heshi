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
	Username       string `json:"username"`
	Cellphone      string `json:"cellphone"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	UserType       string `json:"user_type"`
	RealName       string `json:"real_name"`
	WechatID       string `json:"wechat_id"`
	WechatName     string `json:"wechat_name"`
	WechatQR       string `json:"wechat_qr"`
	Address        string `json:"address"`
	AdditionalInfo string `json:"additional_info"`
	RecommendedBy  string `json:"recommended_by"`
	Icon           string `json:"icon"`
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
		RecommendedBy:  c.PostForm("recommended_by"),
		Icon:           c.PostForm("icon"),
	}

	if vemsg := nu.preValidateNewUser(); vemsg != "" {
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
	if err, q = nu.composeInsertQuery(); err != nil {
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
		RecommendedBy:  c.PostForm("recommended_by"),
		Icon:           c.PostForm("icon"),
	}

	var q string
	var err error
	//TODO validate updated user info too!!!
	if err, q = uu.composeUpdateQuery(); err != nil {
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
				wechat_qr,address,additional_info,recommended_by,icon FROM users WHERE id=?`
	var userType, icon string
	var username, cellphone, email, realName, recommandedBy sql.NullString
	var wechatID, wechatName, wechatQR, address, additionalInfo sql.NullString
	if err := dbQueryRow(q, id).Scan(&username, &cellphone, &email, &realName, &userType, &wechatID,
		&wechatName, &wechatQR, &address, &additionalInfo, &recommandedBy, &icon); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	u := User{
		ID:             id,
		Username:       username.String,
		Cellphone:      cellphone.String,
		Email:          email.String,
		UserType:       userType,
		WechatID:       wechatID.String,
		WechatName:     wechatName.String,
		WechatQR:       wechatID.String,
		Address:        address.String,
		AdditionalInfo: additionalInfo.String,
		RecommendedBy:  recommandedBy.String,
		Icon:           icon,
	}
	c.JSON(http.StatusOK, u)
}

func getAllUsers(c *gin.Context) {
	q := `SELECT id, username,cellphone,email,real_name,user_type,wechat_id,wechat_name,
				wechat_qr,address,additional_info,recommended_by,icon FROM users`
	rows, err := dbQuery(q)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var us []User
	for rows.Next() {
		var id, userType, icon string
		var username, cellphone, email, realName, recommandedBy sql.NullString
		var wechatID, wechatName, wechatQR, address, additionalInfo sql.NullString
		if err := rows.Scan(&id, &username, &cellphone, &email, &realName, &userType, &wechatID,
			&wechatName, &wechatQR, &address, &additionalInfo, &recommandedBy, &icon); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		u := User{
			ID:             id,
			Username:       username.String,
			Cellphone:      cellphone.String,
			Email:          email.String,
			UserType:       userType,
			WechatID:       wechatID.String,
			WechatName:     wechatName.String,
			WechatQR:       wechatID.String,
			Address:        address.String,
			AdditionalInfo: additionalInfo.String,
			RecommendedBy:  recommandedBy.String,
			Icon:           icon,
		}
		us = append(us, u)
	}
	c.JSON(http.StatusOK, us)
}
