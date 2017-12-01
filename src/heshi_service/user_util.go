package main

import (
	"fmt"
	"strings"
	"util"
)

func composeUserQuery(nu User) (error, string) {
	params := userKV(nu)
	q := `INSERT INTO users (id, user_type, password`
	va := fmt.Sprintf(`VALUES ('%s', '%s', '%s'`, nu.ID, nu.UserType, util.Encrypt(nu.Password))
	for k, v := range params {
		q = fmt.Sprintf("%s, %s", q, k)
		va = fmt.Sprintf("%s, '%s'", va, v)
	}
	q = fmt.Sprintf("%s) %s)", q, va)
	return nil, q
}

func composeUserUpdate(uu User) (error, string) {
	params := userKV(uu)
	if uu.Password != "" {
		params["password"] = util.Encrypt(uu.Password)
	}
	q := `UPDATE users SET`
	for k, v := range params {
		q = fmt.Sprintf("%s %s='%s',", q, k, v)
	}

	q = fmt.Sprintf("%s WHERE id='%s'", strings.TrimSuffix(q, ","), uu.ID)
	return nil, q
}

func userKV(nu User) map[string]string {
	params := make(map[string]string)

	if nu.Username != "" {
		params["username"] = nu.Username
	}
	if nu.Cellphone != "" {
		params["cellphone"] = nu.Cellphone
	}
	if nu.Email != "" {
		params["email"] = nu.Email
	}
	if nu.RealName != "" {
		params["real_name"] = nu.RealName
	}
	if nu.WechatID != "" {
		params["wechat_id"] = nu.WechatID
	}
	if nu.WechatName != "" {
		params["wechat_name"] = nu.WechatName
	}
	if nu.WechatQR != "" {
		params["wechat_qr"] = nu.WechatQR
	}
	if nu.Address != "" {
		params["address"] = nu.Address
	}
	if nu.AdditionalInfo != "" {
		params["additional_info"] = nu.AdditionalInfo
	}
	if nu.Icon != "" {
		params["icon"] = nu.Icon
	}
	return params
}
