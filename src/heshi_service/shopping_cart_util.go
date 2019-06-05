package main

import (
	"fmt"
	"strings"
	"time"
)

func (c *cartItemBase) composeInsertQuery() string {
	params := c.paramsKV()
	q := `INSERT INTO shopping_cart (id`
	va := fmt.Sprintf(`VALUES ('%s'`, c.ID)
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

func (c *cartItemBase) composeUpdateQuery() string {
	params := c.paramsKV()
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
		case time.Time:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(time.Time).Format(timeFormat))
		}
	}

	return fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, c.ID)
}

func (c *cartItemBase) paramsKV() map[string]interface{} {
	params := make(map[string]interface{})
	if c.UserID != "" {
		params["user_id"] = c.UserID
	}

	if c.ItemCategory != "" {
		params["item_category"] = strings.ToUpper(c.ItemCategory)
	}
	if c.ItemID != "" {
		params["item_id"] = c.ItemID
	}

	// if c.ItemPrice != 0 {
	// 	params["item_price"] = c.ItemPrice
	// }
	if c.ItemQuantity != 0 {
		params["item_quantity"] = c.ItemQuantity
	}
	if c.ExtraInfo != "" {
		params["extra_info"] = c.ExtraInfo
	}
	return params
}
