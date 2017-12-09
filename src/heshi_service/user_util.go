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
		switch v.(type) {
		case string:
			va = fmt.Sprintf("%s, '%s'", va, v.(string))
		case float64:
			va = fmt.Sprintf("%s, '%f'", va, v.(float64))
		case int:
			va = fmt.Sprintf("%s, '%d'", va, v.(int))
		case int64:
			va = fmt.Sprintf("%s, '%d'", va, v.(int64))
		}
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
	for k, v := range params {
		q = fmt.Sprintf("%s, %s", q, k)
		switch v.(type) {
		case string:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(string))
		case float64:
			q = fmt.Sprintf("%s %s='%f',", q, k, v.(float64))
		case int:
			q = fmt.Sprintf("%s %s='%d',", q, k, v.(int))
		case int64:
			q = fmt.Sprintf("%s %s='%d',", q, k, v.(int64))
		}
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
	if u.AdditionalInfo != "" {
		params["recommended_by"] = u.RecommendedBy
	}
	if u.InvitationCode != "" {
		params["invitation_code"] = u.InvitationCode
	}
	if u.Icon != "" {
		params["icon"] = u.Icon
	}
	return params
}
