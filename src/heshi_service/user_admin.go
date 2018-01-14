package main

import (
	"fmt"
	"heshi/errors"
	"util"
)

type Admin struct {
	UserInfo   User   `json:"user"`
	Level      string `json:"level"`
	WechatKefu string `json:"wechat_kefu"`
	CreatedBy  string `json:"created_by"`
}

func (a *Admin) newAgent() error {
	q := fmt.Sprintf(`INSERT INTO admins (user_id, level, wechat_kefu, created_by) VALUES ('%s', '%s', '%s', '%s')`,
		a.UserInfo.ID, a.Level, a.WechatKefu, a.CreatedBy)
	_, err := db.Exec(q)
	return err
}

func (a *Admin) prevalidateNewAdmin() ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage
	//TODO admin level
	//TODO validate wechat_kefu???
	if !util.IsInArrayString(a.Level, VALID_AGENTLEVEL) {
		vemsg = append(vemsg, VEMSG_AGENT_LEVEL_NOT_VALID)
	}

	if vmsg, err := a.UserInfo.validNewUser(); err != nil {
		return nil, err
	} else if len(vmsg) != 0 {
		vemsg = append(vemsg, vmsg...)
	}
	return vemsg, nil
}
