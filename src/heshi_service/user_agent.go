package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"
	"strconv"
	"util"

	"github.com/gin-gonic/gin"
)

type Agent struct {
	UserInfo    User   `json:"user"`
	Level       string `json:"level"`
	Discount    int    `json:"discount"`
	DiscountStr string `json:"-"`
	SetBy       string `json:"set_by"`
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

func configAgent(c *gin.Context) {
	adminID := c.MustGet("id").(string)
	level := c.PostForm("level")
	if !util.IsInArrayString(level, VALID_AGENTLEVEL) {
		c.JSON(http.StatusOK, VEMSG_AGENT_LEVEL_NOT_VALID)
		return
	}
	agentID := c.Param("id")
	q := fmt.Sprintf(`UPDATE Agents set (level, set_by) VALUES ('%s', '%s') WHERE user_id='%s'`, level, agentID, adminID)
	if _, err := db.Exec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (a *Agent) newAgent() error {
	q := fmt.Sprintf(`INSERT INTO Agents (user_id, level, discount, set_by) VALUES (%s', '%s', '%d', '%s')`,
		a.UserInfo.ID, a.Level, a.Discount, a.SetBy)
	_, err := db.Exec(q)
	return err
}

func agentContactInfo(c *gin.Context) {
	id := c.MustGet("id").(string)

	q := fmt.Sprintf(`SELECT recommanded_by from Users where id=%s`, id)
	var recommandedBy string
	if err := db.QueryRow(q).Scan(&recommandedBy); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	ci, err := getUserContactInfoInvitationCode(recommandedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
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

func (a *Agent) prevalidateNewAgent() ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage
	if !util.IsInArrayString(a.Level, VALID_AGENTLEVEL) {
		vemsg = append(vemsg, VEMSG_AGENT_LEVEL_NOT_VALID)
	}

	discount, err := strconv.Atoi(a.DiscountStr)
	if err != nil {
		vemsg = append(vemsg, VEMSG_AGENT_DISCOUNT_NOT_VALID)
	} else if discount < 0 || discount > 100 {
		vemsg = append(vemsg, VEMSG_AGENT_DISCOUNT_NOT_VALID)
	} else {
		a.Discount = discount
	}
	vmsg, err := a.UserInfo.validNewUser()
	if err != nil {
		return nil, err
	} else if len(vmsg) != 0 {
		vemsg = append(vemsg, vmsg...)
	}
	return vemsg, nil
}
