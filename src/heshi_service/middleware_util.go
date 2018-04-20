package main

import (
	"fmt"
	"strings"
)

func userUsingRecord(URLPath, user, platform, remoteAddr string) error {
	switch strings.ToLower(URLPath) {
	// no need to track
	case "/api/exchangerate", "/api/wechat/status", "/refresh/token":
		return nil
	default:
		//track which product viewed by customer
		if strings.HasPrefix(URLPath, "/api/products/diamonds/") {
			did := strings.Trim(URLPath, "/api/products/diamonds/")
			if err := addUserUsingRecord(user, did, "diamond", platform, remoteAddr); err != nil {
				return err
			}
		}
		if strings.HasPrefix(URLPath, "/api/products/jewelrys/") {
			jid := strings.Trim(URLPath, "/api/products/jewelrys/")
			if err := addUserUsingRecord(user, jid, "jewelry", platform, remoteAddr); err != nil {
				return err
			}
		}
		if strings.HasPrefix(URLPath, "/api/products/gems/") {
			jid := strings.Trim(URLPath, "/api/products/gems/")
			if err := addUserUsingRecord(user, jid, "gem", platform, remoteAddr); err != nil {
				return err
			}
		}
		return addUserActiveRecord(user, URLPath, platform, remoteAddr)
	}
}

//TODO better define
func addUserUsingRecord(user, itemID, itemType, device, remoteAddr string) error {
	q := fmt.Sprintf(`INSERT INTO user_using_records (id, user_id, item_id, item_type, device, remote_addr) 
	VALUES ('%s','%s','%s','%s','%s','%s')`, newV4(), user, itemID, itemType, device, remoteAddr)
	if _, err := dbExec(q); err != nil {
		return err
	}
	return nil
}

//TODO better define
func addUserActiveRecord(user, URLPath, platform, remoteAddr string) error {
	q := fmt.Sprintf(`INSERT INTO user_active_records (id, user_id, page, device, remote_addr) 
	VALUES ('%s','%s','%s','%s','%s')`, newV4(), user, URLPath, platform, remoteAddr)
	if _, err := dbExec(q); err != nil {
		return err
	}
	return nil
}
