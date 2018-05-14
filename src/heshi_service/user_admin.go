package main

import (
	"fmt"
	"heshi/errors"
	"net/http"
	"strconv"
	"util"

	"github.com/gin-gonic/gin"
)

// Admin ...
type Admin struct {
	Level      int    `json:"admin_level,omitempty"`
	LevelStr   string `json:"-"`
	WechatKefu string `json:"wechat_kefu,omitempty"`
	CreatedBy  string `json:"created_by,omitempty"`
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
		c.JSON(http.StatusOK, vemsgNotValid)
		return
	}
	agentID := c.Param("id")
	q := fmt.Sprintf(`UPDATE admins SET created_by='%s'`, agentID)

	if levelStr != "" {
		if !util.IsInArrayString(levelStr, validAgentLevel) {
			c.JSON(http.StatusOK, vemsgAgentLevelNotValid)
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

func (a *User) newAdmin() error {
	q := fmt.Sprintf(`INSERT INTO admins (user_id, level, wechat_kefu, created_by) VALUES ('%s', '%d', '%s', '%s')`,
		a.ID, a.Admin.Level, a.Admin.WechatKefu, a.Admin.CreatedBy)
	_, err := dbExec(q)
	return err
}

func (a *User) prevalidateNewAdmin() ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage
	//TODO admin level
	//TODO validate wechat_kefu???
	if !util.IsInArrayString(a.Admin.LevelStr, validAdminLevel) {
		vemsg = append(vemsg, vemsgAdminLevelNotValid)
	} else {
		level, err := strconv.Atoi(a.Admin.LevelStr)
		if err != nil {
			vemsg = append(vemsg, vemsgAdminLevelNotValid)
		} else if level < 0 || level > 10 {
			vemsg = append(vemsg, vemsgAdminLevelNotValid)
		} else {
			a.Admin.Level = level
		}
	}

	if vmsg, err := a.validNewUser(); err != nil {
		return nil, err
	} else if len(vmsg) != 0 {
		vemsg = append(vemsg, vmsg...)
	}
	return vemsg, nil
}
