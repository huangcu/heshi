package main

import "fmt"

func (oi *orderItem) composeInsertQuery() string {
	params := oi.parmsKV()
	q := `INSERT INTO jewelrys (id`
	va := fmt.Sprintf(`VALUES ('%s'`, oi.ID)
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

func (oi *orderItem) composeUpdateQuery() string {
	params := oi.parmsKV()
	q := `UPDATE jewelrys SET`
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
	q = fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, oi.ID)
	return q
}

func (oi *orderItem) parmsKV() map[string]interface{} {
	params := make(map[string]interface{})
	if oi.ItemID != "" {
		params["item_id"] = oi.ItemID
	}
	if oi.Category != "" {
		params["item_category"] = oi.Category
	}
	if oi.Quantity != 0 {
		params["item_quantity"] = oi.Quantity
	}
	if oi.Price != 0 {
		params["item_price"] = oi.Price
	}
	if oi.ExtraInfo != "" {
		params["extra_info"] = oi.ExtraInfo
	}
	if oi.SpecialNotice != "" {
		params["special_notice"] = oi.SpecialNotice
	}
	if oi.DownPayment != 0 {
		params["downpayment"] = oi.DownPayment
	}
	if oi.BuyerID != "" {
		params["buyer_id"] = oi.BuyerID
	}
	if oi.TransactionID != "" {
		params["transaction_id"] = oi.TransactionID
	}
	if oi.Status != "" {
		params["status"] = oi.Status
	}
	return params
}
