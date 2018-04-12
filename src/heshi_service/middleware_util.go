package main

import (
	"fmt"
	"strings"
)

func userUsingRecord(URLPath, user, platform string) error {
	switch strings.ToLower(URLPath) {
	// no need to track
	case "/api/exchangerate", "/api/wechat/status", "/refresh/token":
		return nil
	default:
		//track which product viewed by customer
		if strings.HasPrefix(URLPath, "/api/products/diamonds/") {
			did := strings.Trim(URLPath, "/api/products/diamonds/")
			if err := addUserUsingRecord(user, did, "diamond", platform); err != nil {
				return err
			}
		}
		if strings.HasPrefix(URLPath, "/api/products/jewelrys/") {
			jid := strings.Trim(URLPath, "/api/products/jewelrys/")
			if err := addUserUsingRecord(user, jid, "jewelry", platform); err != nil {
				return err
			}
		}
		if strings.HasPrefix(URLPath, "/api/products/gems/") {
			jid := strings.Trim(URLPath, "/api/products/gems/")
			if err := addUserUsingRecord(user, jid, "gem", platform); err != nil {
				return err
			}
		}
		return addUserActiveRecord(user, URLPath)
	}
}

//TODO better define
func addUserUsingRecord(user, itemID, itemType, device string) error {
	q := fmt.Sprintf(`INSERT INTO user_using_records (id, user_id, item_id, item_type, device) 
	VALUES ('%s','%s','%s','%s','%s')`, newV4(), user, itemID, itemType, device)
	if _, err := dbExec(q); err != nil {
		return err
	}
	return nil
}

//TODO better define
func addUserActiveRecord(user, URLPath string) error {
	q := fmt.Sprintf(`INSERT INTO user_active_records (id, user_id, page) 
	VALUES ('%s','%s','%s')`, newV4(), user, URLPath)
	if _, err := dbExec(q); err != nil {
		return err
	}
	return nil
}
