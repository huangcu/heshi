package main

import (
	"fmt"
	"heshi/errors"
	"net/http"
	"strconv"
	"util"

	"github.com/gin-gonic/gin"
)

type Admin struct {
	User
	// UserInfo   User   `json:"user"`
	Level      int    `json:"admin_level"`
	LevelStr   string `json:"-"`
	WechatKefu string `json:"wechat_kefu"`
	CreatedBy  string `json:"created_by"`
}

func getAdmin(uid string) (*Admin, error) {
	var level int
	var wechatKefu, createdBy string
	q := fmt.Sprintf("SELECT level, wechat_kefu, created_by FROM admins WHERE user_id='%s'", uid)
	if err := dbQueryRow(q).Scan(&level, &wechatKefu, &createdBy); err != nil {
		return nil, err
	}
	a := &Admin{
		Level:      level,
		WechatKefu: wechatKefu,
		CreatedBy:  createdBy,
	}
	return a, nil
}

func updateAdmin(c *gin.Context) {
	adminID := c.MustGet("id").(string)
	levelStr := c.PostForm("level")
	kefu := c.PostForm("wechat_kefu")
	if levelStr == "" && kefu == "" {
		c.JSON(http.StatusOK, VEMSG_NOT_VALID)
	}
	agentID := c.Param("id")
	q := fmt.Sprintf(`UPDATE admins SET created_by='%s'`, agentID)

	if levelStr != "" {
		if !util.IsInArrayString(levelStr, VALID_AGENTLEVEL) {
			c.JSON(http.StatusOK, VEMSG_AGENT_LEVEL_NOT_VALID)
			return
		}
		level, _ := strconv.Atoi(levelStr)
		q = fmt.Sprintf("%s, level='%d'", q, level)
	}

	if kefu != "" {
		q = fmt.Sprintf("%s, wechat_kefu='%s'", q, kefu)
	}
	q = fmt.Sprintf("%s WHERE user_id='%s'", q, adminID)
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (a *Admin) newAdmin() error {
	q := fmt.Sprintf(`INSERT INTO admins (user_id, level, wechat_kefu, created_by) VALUES ('%s', '%d', '%s', '%s')`,
		a.ID, a.Level, a.WechatKefu, a.CreatedBy)
	_, err := dbExec(q)
	return err
}

func (a *Admin) prevalidateNewAdmin() ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage
	//TODO admin level
	//TODO validate wechat_kefu???
	if !util.IsInArrayString(a.LevelStr, VALID_ADMINLEVEL) {
		vemsg = append(vemsg, VEMSG_ADMIN_LEVEL_NOT_VALID)
	} else {
		level, err := strconv.Atoi(a.LevelStr)
		if err != nil {
			vemsg = append(vemsg, VEMSG_ADMIN_LEVEL_NOT_VALID)
		} else if level < 0 || level > 10 {
			vemsg = append(vemsg, VEMSG_ADMIN_LEVEL_NOT_VALID)
		} else {
			a.Level = level
		}
	}

	if vmsg, err := a.validNewUser(); err != nil {
		return nil, err
	} else if len(vmsg) != 0 {
		vemsg = append(vemsg, vmsg...)
	}
	return vemsg, nil
}
