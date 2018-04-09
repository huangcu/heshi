package main

import (
	"fmt"
	"time"
)

func (s *shoppingItem) composeInsertQuery() string {
	params := s.paramsKV()
	q := `INSERT INTO interested_items (id`
	va := fmt.Sprintf(`VALUES ('%s'`, s.ID)
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
	return fmt.Sprintf("%s) %s)", q, va)
}

func (s *shoppingItem) composeUpdateQuery() string {
	params := s.paramsKV()
	q := `UPDATE interested_items SET`
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

	return fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, s.ID)
}

func (s *shoppingItem) paramsKV() map[string]interface{} {
	params := make(map[string]interface{})
	if s.UserID != "" {
		params["user_id"] = s.UserID
	}

	if s.ItemType != "" {
		params["item_type"] = s.ItemType
	}
	if s.ItemID != "" {
		params["item_id"] = s.ItemID
	}

	if s.ItemAccessory != 0 {
		params["item_accessory"] = s.ItemAccessory
	}
	if s.ConfirmedForCheck != "" {
		params["confirmed_for_check"] = s.ConfirmedForCheck
	}
	if s.Available != "" {
		params["available"] = s.Available
	}

	if s.SpecialNotice != "" {
		params["special_notice"] = s.SpecialNotice
	}
	return params
}
