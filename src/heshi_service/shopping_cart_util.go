package main

import (
	"fmt"
)

func (s *shoppingCartItem) composeInsertQuery() string {
	params := s.paramsKV()
	q := `INSERT INTO shopping_cart (id`
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
		}
	}
	return fmt.Sprintf("%s) %s)", q, va)
}

func (s *shoppingCartItem) composeUpdateQuery() string {
	params := s.paramsKV()
	q := `UPDATE shopping_cart SET`
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
		}
	}

	return fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, s.ID)
}

func (s *shoppingCartItem) paramsKV() map[string]interface{} {
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

	if s.ItemPrice != 0 {
		params["item_price"] = s.ItemPrice
	}
	if s.ItemQuantity != 0 {
		params["item_quantity"] = s.ItemQuantity
	}
	if s.ExtraInfo != "" {
		params["extra_info"] = s.ExtraInfo
	}
	return params
}
