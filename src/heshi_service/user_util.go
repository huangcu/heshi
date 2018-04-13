package main

import (
	"fmt"
	"time"
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
		case time.Time:
			va = fmt.Sprintf("%s, '%s'", va, v.(time.Time).Format(timeFormat))
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
		switch v.(type) {
		case string:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(string))
		case float64:
			q = fmt.Sprintf("%s %s='%f',", q, k, v.(float64))
		case int:
			q = fmt.Sprintf("%s %s='%d',", q, k, v.(int))
		case int64:
			q = fmt.Sprintf("%s %s='%d',", q, k, v.(int64))
		case time.Time:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(time.Time).Format(timeFormat))
		}
	}

	q = fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, u.ID)
	return q
}

func (u *User) paramsKV() map[string]interface{} {
	params := make(map[string]interface{})

	if u.Username != "" {
		params["username"] = u.Username
	} else {
		params["username"] = fmt.Sprintf("hs_%s", util.RandStringRunes(13))
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
	} else {
		params["invitation_code"] = util.NewUniqueId()
	}
	if u.Icon != "" {
		params["icon"] = u.Icon
	}
	return params
}

func isUserExistByID(id string) (bool, error) {
	var count int
	if err := dbQueryRow(fmt.Sprintf("SELECT COUNT(*) FROM users WHERE id='%s'", id)).Scan(&count); err != nil {
		return false, err
	}
	return count == 1, nil
}

func isUserExistByWechatOpenID(wecahtOpenID string) (bool, error) {
	var count int
	if err := dbQueryRow(fmt.Sprintf("SELECT COUNT(*) FROM users WHERE wechat_id='%s'", wecahtOpenID)).Scan(&count); err != nil {
		return false, err
	}
	return count == 1, nil
}
