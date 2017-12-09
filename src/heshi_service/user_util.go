package main

import (
	"fmt"
	"strings"
	"util"
)

func (u *User) composeInsertQuery() string {
	params := u.paramsKV()
	q := `INSERT INTO users (id, user_type, password`
	va := fmt.Sprintf(`VALUES ('%s', '%s', '%s'`, u.ID, u.UserType, util.Encrypt(u.Password))
	for k, v := range params {
		q = fmt.Sprintf("%s, %s", q, k)
		va = fmt.Sprintf("%s, '%s'", va, v)
	}
	q = fmt.Sprintf("%s) %s)", q, va)
	return q
}

func (u *User) composeUpdateQuery() string {
	params := u.paramsKV()
	if u.Password != "" {
		params["password"] = util.Encrypt(u.Password)
	}
	q := `UPDATE users SET`
	for k, v := range params {
		q = fmt.Sprintf("%s %s='%s',", q, k, v)
	}

	q = fmt.Sprintf("%s WHERE id='%s'", strings.TrimSuffix(q, ","), u.ID)
	return q
}

func (u *User) paramsKV() map[string]interface{} {
	params := make(map[string]interface{})

	if u.Username != "" {
		params["username"] = u.Username
	}
	if u.Cellphone != "" {
		params["cellphone"] = u.Cellphone
	}
	if u.Email != "" {
		params["email"] = u.Email
	}
	if u.RealName != "" {
		params["real_name"] = u.RealName
	}
	if u.WechatID != "" {
		params["wechat_id"] = u.WechatID
	}
	if u.WechatName != "" {
		params["wechat_name"] = u.WechatName
	}
	if u.WechatQR != "" {
		params["wechat_qr"] = u.WechatQR
	}
	if u.Address != "" {
		params["address"] = u.Address
	}
	if u.AdditionalInfo != "" {
		params["additional_info"] = u.AdditionalInfo
	}
	if u.Icon != "" {
		params["icon"] = u.Icon
	}
	return params
}
