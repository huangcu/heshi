package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"

	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
)

type WechatUserInfo struct {
	*mpoauth2.UserInfo
}

func (wu *WechatUserInfo) newWechatUser() error {
	q := fmt.Sprintf(`INSERT INTO wechat_users (openid, nickname,sex,city, province,contry, head_image_url,privilege, unionid) 
	VALUES (%s,%s,%d,%s,%s,%s,%s,%s,%s)`,
		wu.OpenId, wu.Nickname, wu.Sex, wu.City, wu.Province, wu.Country, wu.HeadImageURL, wu.Privilege, wu.UnionId)
	_, err := db.Exec(q)
	return err
}

func isWechatUserExist(openid string) (bool, error) {
	q := `SELECT COUNT(*) FROM wechat_users WHERE openid=?`
	var count int
	if err := db.QueryRow(q, openid).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if count == 1 {
		return true, nil
	}
	return true, errors.Newf("query returns :%d records", count)
}
