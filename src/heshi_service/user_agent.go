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
	Level       int    `json:"level,omitempty"`
	LevelStr    string `json:"-"`
	Discount    int    `json:"discount,omitempty"`
	DiscountStr string `json:"-"`
	CreatedBy   string `json:"created_by,omitempty"`
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

func getAgent(uid string) (*Agent, error) {
	var level, discount int
	var createdBy string
	q := fmt.Sprintf("SELECT level, discount, created_by FROM agents WHERE user_id='%s'", uid)
	if err := dbQueryRow(q).Scan(&level, &discount, &createdBy); err != nil {
		return nil, err
	}
	a := &Agent{
		Level:     level,
		Discount:  discount,
		CreatedBy: createdBy,
	}
	fmt.Println(a)
	return a, nil
}

func updateAgent(c *gin.Context) {
	adminID := c.MustGet("id").(string)
	levelStr := c.PostForm("level")
	discountStr := c.PostForm("discount")
	if levelStr == "" && discountStr == "" {
		c.JSON(http.StatusOK, vemsgNotValid)
		return
	}
	agentID := c.Param("id")
	q := fmt.Sprintf(`UPDATE agents SET created_by='%s'`, agentID)

	if levelStr != "" {
		if !util.IsInArrayString(levelStr, VALID_AGENTLEVEL) {
			c.JSON(http.StatusOK, vemsgAgentLevelNotValid)
			return
		}
		level, _ := strconv.Atoi(levelStr)
		q = fmt.Sprintf("%s, level='%d'", q, level)
	}

	if discountStr != "" {
		discount, err := strconv.Atoi(discountStr)
		if err != nil {
			c.JSON(http.StatusOK, vemsgAgentLevelNotValid)
			return
		} else if discount < 0 || discount > 100 {
			c.JSON(http.StatusOK, vemsgAgentLevelNotValid)
			return
		} else {
			q = fmt.Sprintf("%s, discount='%d'", q, discount)
		}
	}
	q = fmt.Sprintf("%s WHERE user_id='%s'", q, adminID)
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (a *User) newAgent() error {
	q := fmt.Sprintf(`INSERT INTO agents (user_id, level, discount, created_by) VALUES (%s', '%d', '%d', '%s')`,
		a.ID, a.Agent.Level, a.Agent.Discount, a.Agent.CreatedBy)
	_, err := dbExec(q)
	return err
}

//TODO what kind of info should pass to agent
func getUsersRecommendedByAgent(c *gin.Context) {
	agentID := c.MustGet("id").(string)
	us, err := getUsersRecommendedBy(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, us)
}

func agentContactInfo(c *gin.Context) {
	id := c.MustGet("id").(string)

	q := fmt.Sprintf(`SELECT recommended_by from users where id=%s`, id)
	var recommendedBy string
	if err := dbQueryRow(q).Scan(&recommendedBy); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	ci, err := getUserContactInfoInvitationCode(recommendedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, *ci)
}

func getUserContactInfoInvitationCode(userID string) (*ContactInfo, error) {
	var cellphone, email, realName sql.NullString
	var wechatID, wechatName, wechatQR, address, additionalInfo sql.NullString
	q := fmt.Sprintf(`SELECT cellphone, email, realname, wechat_id, wechat_name, wechat_qr, address, additional_info 
	FROM users WHERE id='%s'`, userID)
	if err := dbQueryRow(q).Scan(&cellphone, &email, &realName, &wechatID, &wechatName, &wechatQR, &address, &additionalInfo); err != nil {
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

func (a *User) prevalidateNewAgent() ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage
	if !util.IsInArrayString(a.Agent.LevelStr, VALID_AGENTLEVEL) {
		vemsg = append(vemsg, vemsgAgentLevelNotValid)
	} else {
		level, err := strconv.Atoi(a.Agent.LevelStr)
		if err != nil {
			vemsg = append(vemsg, vemsgAgentLevelNotValid)
		} else if level < 0 || level > 10 {
			vemsg = append(vemsg, vemsgAgentLevelNotValid)
		} else {
			a.Agent.Level = level
		}
	}

	discount, err := strconv.Atoi(a.Agent.DiscountStr)
	if err != nil {
		vemsg = append(vemsg, vemsgAgentDiscountNotValid)
	} else if discount < 0 || discount > 100 {
		vemsg = append(vemsg, vemsgAgentDiscountNotValid)
	} else {
		a.Agent.Discount = discount
	}
	vmsg, err := a.validNewUser()
	if err != nil {
		return nil, err
	} else if len(vmsg) != 0 {
		vemsg = append(vemsg, vmsg...)
	}
	return vemsg, nil
}
