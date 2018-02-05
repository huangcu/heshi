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
	User
	// UserInfo    User   `json:"user"`
	Level       int    `json:"agent_level"`
	LevelStr    string `json:"-"`
	Discount    int    `json:"agent_discount"`
	DiscountStr string `json:"-"`
	CreatedBy   string `json:"created_by"`
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
	return a, nil
}

func updateAgent(c *gin.Context) {
	adminID := c.MustGet("id").(string)
	levelStr := c.PostForm("level")
	discountStr := c.PostForm("discount")
	if levelStr == "" && discountStr == "" {
		c.JSON(http.StatusOK, vemsgNotValid)
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

func (a *Agent) newAgent() error {
	q := fmt.Sprintf(`INSERT INTO agents (user_id, level, discount, created_by) VALUES (%s', '%d', '%d', '%s')`,
		a.ID, a.Level, a.Discount, a.CreatedBy)
	_, err := dbExec(q)
	return err
}

//find user is recommanded by - user_id
//from invitation code, get which user recommanded this
// if the recommand is agent ???
func agentContactInfo(c *gin.Context) {
	id := c.MustGet("id").(string)

	q := fmt.Sprintf(`SELECT recommanded_by from users where id=%s`, id)
	var recommandedBy string
	if err := dbQueryRow(q).Scan(&recommandedBy); err != nil {
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

// func userRecommandedByAgent(uid string) (string, error) {
// 	q := fmt.Sprintf(`SELECT recommanded_by from users where id=%s`, uid)
// 	var recommandedBy string
// 	if err := dbQueryRow(q).Scan(&recommandedBy); err != nil {
// 		return "", err
// 	}

// }

func getUserContactInfoInvitationCode(userID string) (*ContactInfo, error) {
	// var userID string
	// if err := dbQueryRow("SELECT user_id from invitation_codes WHERE invitation_code=?", code).Scan(&userID); err != nil {
	// 	return nil, err
	// }
	var cellphone, email, realName sql.NullString
	var wechatID, wechatName, wechatQR, address, additionalInfo sql.NullString
	q := `SELECT cellphone, email, realname, wechat_id, wechat_name, wechat_qr, address, additional_info from users where id=?`
	if err := dbQueryRow(q, userID).Scan(&cellphone, &email, &realName, &wechatID, &wechatName, &wechatQR, &address, &additionalInfo); err != nil {
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
	if !util.IsInArrayString(a.LevelStr, VALID_AGENTLEVEL) {
		vemsg = append(vemsg, vemsgAgentLevelNotValid)
	} else {
		level, err := strconv.Atoi(a.LevelStr)
		if err != nil {
			vemsg = append(vemsg, vemsgAgentLevelNotValid)
		} else if level < 0 || level > 10 {
			vemsg = append(vemsg, vemsgAgentLevelNotValid)
		} else {
			a.Level = level
		}
	}

	discount, err := strconv.Atoi(a.DiscountStr)
	if err != nil {
		vemsg = append(vemsg, vemsgAgentDiscountNotValid)
	} else if discount < 0 || discount > 100 {
		vemsg = append(vemsg, vemsgAgentDiscountNotValid)
	} else {
		a.Discount = discount
	}
	vmsg, err := a.validNewUser()
	if err != nil {
		return nil, err
	} else if len(vmsg) != 0 {
		vemsg = append(vemsg, vmsg...)
	}
	return vemsg, nil
}
