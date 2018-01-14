package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID                  string  `json:"id"`
	Username            string  `json:"username"`
	Cellphone           string  `json:"cellphone"`
	Email               string  `json:"email"`
	Password            string  `json:"-"`
	UserType            string  `json:"user_type"`
	RealName            string  `json:"real_name"`
	WechatID            string  `json:"wechat_id"`
	WechatName          string  `json:"wechat_name"`
	WechatQR            string  `json:"wechat_qr"`
	Address             string  `json:"address"`
	AdditionalInfo      string  `json:"additional_info"`
	RecommendedBy       string  `json:"recommended_by"`
	InvitationCode      string  `json:"invitation_code"`
	Discount            float64 `json:"discount"`
	Point               int     `json:"point"`
	TotalPurchaseAmount float64 `json:"total_purchase_amount"`
	Icon                string  `json:"icon"`
	// CreatedAt      time.Time `json:"created_at"`
	// UpdatedAt      time.Time `json:"updated_at"`
}

func newAdminAgentUser(c *gin.Context) {
	adminID := c.MustGet("id").(string)
	userType := c.PostForm("user_type")
	if userType != AGENT && userType != ADMIN {
		VEMSG_USER_USERTYPE_NOT_VALID.Message = fmt.Sprintf("user type can only be %s or %s", ADMIN, AGENT)
		c.JSON(http.StatusOK, VEMSG_USER_USERTYPE_NOT_VALID)
		return
	}

	nu := User{
		ID:             uuid.NewV4().String(),
		Username:       c.PostForm("username"),
		Cellphone:      c.PostForm("cellphone"),
		Email:          c.PostForm("email"),
		Password:       c.PostForm("password"),
		UserType:       userType,
		RealName:       c.PostForm("real_name"),
		WechatID:       c.PostForm("wechat_id"),
		WechatName:     c.PostForm("wechat_name"),
		WechatQR:       c.PostForm("wechat_qr"),
		Address:        c.PostForm("address"),
		AdditionalInfo: c.PostForm("additional_info"),
		RecommendedBy:  c.PostForm("recommended_by"),
		Icon:           c.PostForm("icon"),
	}

	if userType == AGENT {
		a := Agent{
			UserInfo:    nu,
			Level:       c.PostForm("level"),
			DiscountStr: c.PostForm("discount"),
			SetBy:       adminID,
		}
		if vemsg, err := a.prevalidateNewAgent(); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		} else if len(vemsg) != 0 {
			c.JSON(http.StatusOK, vemsg)
			return
		}

		q := nu.composeInsertQuery()
		if _, err := db.Exec(q); err != nil {
			c.JSON(http.StatusBadRequest, errors.GetMessage(err))
			return
		}
		if err := a.newAgent(); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, a.UserInfo.ID)
	}

	if userType == ADMIN {
		a := Admin{
			UserInfo:   nu,
			Level:      c.PostForm("level"),
			WechatKefu: c.PostForm("wechat_kefu"),
			CreatedBy:  adminID,
		}
		if vemsg, err := a.prevalidateNewAdmin(); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		} else if len(vemsg) != 0 {
			c.JSON(http.StatusOK, vemsg)
			return
		}

		q := nu.composeInsertQuery()
		if _, err := db.Exec(q); err != nil {
			c.JSON(http.StatusBadRequest, errors.GetMessage(err))
			return
		}
		if err := a.newAgent(); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, a.UserInfo.ID)
	}
}

func newUser(c *gin.Context) {
	nu := User{
		ID:             uuid.NewV4().String(),
		Username:       c.PostForm("username"),
		Cellphone:      c.PostForm("cellphone"),
		Email:          c.PostForm("email"),
		Password:       c.PostForm("password"),
		UserType:       CUSTOMER,
		RealName:       c.PostForm("real_name"),
		WechatID:       c.PostForm("wechat_id"),
		WechatName:     c.PostForm("wechat_name"),
		WechatQR:       c.PostForm("wechat_qr"),
		Address:        c.PostForm("address"),
		AdditionalInfo: c.PostForm("additional_info"),
		RecommendedBy:  c.PostForm("recommended_by"),
		Icon:           c.PostForm("icon"),
	}

	if vemsg, err := nu.validNewUser(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsg) != 0 {
		c.JSON(http.StatusOK, vemsg)
		return
	}

	q := nu.composeInsertQuery()
	if _, err := db.Exec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	s := sessions.Default(c)
	s.Set(USER_SESSION_KEY, nu.ID)
	s.Save()

	c.JSON(http.StatusOK, nu.ID)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.MustGet("id").(string)
	}
	uu := User{
		ID:             id,
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

	//TODO validate updated user info too!!!
	//TODO what info can be updated!!
	q := uu.composeUpdateQuery()
	//TODO admin,agent update!!!!
	// var userType string
	// switch userType {
	// case "admin":
	// case "agent":
	// default:
	// }
	if _, err := db.Exec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, uu.ID)
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.MustGet("id").(string)
	}
	q := selectUserQuery(id)

	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	us, err := composeUser(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if us == nil {
		VEMSG_USER_NOT_EXIST.Message = fmt.Sprintf("Fail to find user with id: %s", c.Param("id"))
		c.JSON(http.StatusOK, VEMSG_USER_NOT_EXIST)
		return
	}
	c.JSON(http.StatusOK, us)
}

func getAllUsers(c *gin.Context) {
	q := selectUserQuery("")
	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	us, err := composeUser(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if us == nil {
		VEMSG_USER_NOT_EXIST.Message = "Fail to find users"
		c.JSON(http.StatusOK, VEMSG_USER_NOT_EXIST)
		return
	}
	c.JSON(http.StatusOK, us)
}

//TODO
func removeUser(c *gin.Context) {
	uid := c.Param("id")
	q := "SELECT user_type from users WHERE id=?"
	var userType string
	if err := db.QueryRow(q, uid).Scan(&userType); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, "wrong user")
			return
		}
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	switch userType {
	case "admin":
		q = `DELETE FROM admins WHERE user_id=?`
		if _, err := db.Exec(q, uid); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	case "agent":
		q = `DELETE FROM agents WHERE user_id=?`
		if _, err := db.Exec(q, uid); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	default:
		q = `DELETE FROM users WHERE id=?`
		if _, err := db.Exec(q, uid); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	}
}

func composeUser(rows *sql.Rows) ([]User, error) {
	var id, userType, icon, invitationCode string
	var username, cellphone, email, realName, recommandedBy sql.NullString
	var wechatID, wechatName, wechatQR, address, additionalInfo sql.NullString
	var discount, point int
	var totalPurchaseAmount float64

	var us []User
	for rows.Next() {
		if err := rows.Scan(&id, &username, &cellphone, &email, &realName, &userType, &wechatID,
			&wechatName, &wechatQR, &address, &additionalInfo, &recommandedBy, &invitationCode,
			&discount, &point, &totalPurchaseAmount, &icon); err != nil {
			return nil, err
		}
		u := User{
			ID:                  id,
			Username:            username.String,
			Cellphone:           cellphone.String,
			Email:               email.String,
			UserType:            userType,
			WechatID:            wechatID.String,
			WechatName:          wechatName.String,
			WechatQR:            wechatID.String,
			Address:             address.String,
			AdditionalInfo:      additionalInfo.String,
			RecommendedBy:       recommandedBy.String,
			InvitationCode:      invitationCode,
			Discount:            float64(discount) / 100,
			Point:               point,
			TotalPurchaseAmount: totalPurchaseAmount,
			Icon:                icon,
		}
		us = append(us, u)
	}
	return us, nil
}

func selectUserQuery(id string) string {
	q := `SELECT id,username,cellphone,email,real_name,user_type,wechat_id,
	wechat_name,wechat_qr,address,additional_info,recommended_by,invitation_code,
	discount,point,total_purchase_amount,icon FROM users`

	if id != "" {
		q = fmt.Sprintf("%s WHERE id='%s'", q, id)
	}
	return q
}
