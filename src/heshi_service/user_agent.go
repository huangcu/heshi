package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Agent struct {
	ID       string  `json"-"`
	UserInfo User    `json:"user"`
	Level    int     `json:"level"`
	Discount float64 `json:"discount"`
}

type ContactInfo struct {
	ID             string `json:"-"`
	Cellphone      string `json:"cellphone"`
	Email          string `json:"email"`
	RealName       string `json:"real_name"`
	WechatID       string `json:"wechat_id"`
	WechatName     string `json:"wechat_name"`
	WechatQR       string `json:"wechat_qr"`
	Address        string `json:"address"`
	AdditionalInfo string `json:"additional_info"`
}

func (a *Agent) newAgent() error {
	q := fmt.Sprintf(`INSERT INTO Agents (id, user_id, level, discount) VALUES(%s, %s, %d, %f)`, a.ID, a.UserInfo.ID, a.Level, a.Discount)
	_, err := db.Exec(q)
	return err
}

func agentContactInfo(c *gin.Context) {
	id := c.MustGet("id").(string)

	q := fmt.Sprintf(`SELECT recommanded_by from Users where id=%s`, id)
	var recommandedBy string
	if err := db.QueryRow(q).Scan(&recommandedBy); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ci, err := getUserContactInfoInvitationCode(recommandedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, *ci)
}

func getUserContactInfoInvitationCode(code string) (*ContactInfo, error) {
	var userID string
	if err := db.QueryRow("SELECT user_id from invitation_codes WHERE invitation_code=?", code).Scan(&userID); err != nil {
		return nil, err
	}
	var cellphone, email, realName sql.NullString
	var wechatID, wechatName, wechatQR, address, additionalInfo sql.NullString
	q := `SELECT cellphone, email, realname, wechat_id, wechat_name, wechat_qr, address, additional_info from users where id=?`
	if err := db.QueryRow(q, userID).Scan(&cellphone, &email, &realName, &wechatID, &wechatName, &wechatQR, &address, &additionalInfo); err != nil {
		return nil, err
	}
	return &ContactInfo{
		ID:             userID,
		Cellphone:      cellphone.String,
		Email:          email.String,
		RealName:       realName.String,
		WechatID:       wechatID.String,
		WechatName:     wechatName.String,
		WechatQR:       wechatQR.String,
		Address:        address.String,
		AdditionalInfo: additionalInfo.String,
	}, nil
}
