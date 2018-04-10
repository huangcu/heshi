package main

import (
	"fmt"
	"time"
	"util"
)

func (oi *orderItem) composeInsertQuery() string {
	params := oi.parmsKV()
	q := `INSERT INTO orders (id`
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
		case time.Time:
			va = fmt.Sprintf("%s, '%s'", va, v.(time.Time).Format(timeFormat))
		}
	}
	q = fmt.Sprintf("%s) %s)", q, va)
	return q
}

func (oi *orderItem) composeUpdateQuery() string {
	params := oi.parmsKV()
	q := `UPDATE orders SET`
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
	q = fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, oi.ID)
	return q
}

func (oi *orderItem) parmsKV() map[string]interface{} {
	params := make(map[string]interface{})
	if oi.ItemID != "" {
		params["item_id"] = oi.ItemID
	}
	if oi.ItemCategory != "" {
		params["item_category"] = oi.ItemCategory
	}
	if oi.ItemQuantity != 0 {
		params["item_quantity"] = oi.ItemQuantity
	}
	if oi.ItemPrice != 0 {
		params["item_price"] = oi.ItemPrice
	}
	if oi.SoldPriceUSD != 0 {
		params["sold_price_usd"] = oi.SoldPriceUSD
	}
	if oi.SoldPriceCNY != 0 {
		params["sold_price_cny"] = oi.SoldPriceCNY
	}
	if oi.SoldPriceEUR != 0 {
		params["sold_price_eur"] = oi.SoldPriceEUR
	}
	if oi.ReturnPoint != 0 {
		params["return_point"] = oi.ReturnPoint
	}
	if oi.ChosenBy != "" {
		params["chosen_by"] = oi.ChosenBy
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

func isOrderExistByID(id string) (bool, error) {
	var count int
	if err := dbQueryRow("SELECT COUNT(*) FROM orders WHERE id=?", id).Scan(&count); err != nil {
		return false, err
	}
	return count == 1, nil
}
func isTransactionExistByID(tid string) (bool, error) {
	var count int
	if err := dbQueryRow("SELECT COUNT(*) FROM orders WHERE transaction_id=?", tid).Scan(&count); err != nil {
		return false, err
	}
	return count == 1, nil
}
func getOrderStatusByID(id string) (string, error) {
	var status string
	if err := dbQueryRow("SELECT status FROM orders WHERE id=?", id).Scan(&status); err != nil {
		return "", err
	}
	return status, nil
}

func getOrderByID(oid string) (*orderItem, error) {
	q := fmt.Sprintf(`SELECT id, item_id, item_price, item_category, item_quantity, downpayment, 
		buyer_id, transaction_id,	sold_price_usd, sold_price_cny, sold_price_eur,
		return_point, chosen_by, status, extra_info, special_notice 
		FROM orders 
		WHERE id='%s'`, oid)
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ois, err := composeOrders(rows)
	if err != nil {
		return nil, err
	}
	if len(ois) == 1 {
		return &ois[0], nil
	}
	return nil, nil
}
func isOrderStatusChangeAllowed(nowStatus, newStatus string) bool {
	// ADMIN can cancel anytime
	if newStatus == MCANCELLED {
		return true
	}
	switch nowStatus {
	case ORDERED:
		return newStatus == MDOWNPAYMENT
	case DOWNPAYMENT, MDOWNPAYMENT:
		return newStatus == MPAID
	case PAID, MPAID:
		return newStatus == MDELIVERED
	case DELIVERED, MDELIVERED:
		return newStatus == MRECEIVED
	default:
		return false
	}
}

//TODO check transaction, if no downpayment in 24 hrs of latest update, cancel it.
func longRunTransactionCheck() error {
	q := fmt.Sprintf(`SELECT id, item_id, item_category, item_quantity, transaction_id 
	FROM orders 
	WHERE updated_at < timestampadd(hour, -24, now()) 
	AND status='%s'`, ORDERED)
	rows, err := dbQuery(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ois []orderItem
	for rows.Next() {
		var id, transactionID, itemID, itemCategory string
		var itemQuantity int
		if err := rows.Scan(&id, &itemID, &itemCategory, &itemQuantity, &transactionID); err != nil {
			return err
		}
		oi := orderItem{
			ID:           id,
			ItemID:       itemID,
			ItemCategory: itemCategory,
			ItemQuantity: itemQuantity,
			Status:       CANCELLED,
		}
		ois = append(ois, oi)
	}

	transationOrderItemmap := make(map[string][]orderItem)
	for _, oItem := range ois {
		key := oItem.TransactionID
		transationOrderItemmap[key] = append(transationOrderItemmap[key], oItem)
	}

	var ts []transaction
	for transactionID, orderItems := range transationOrderItemmap {
		t := transaction{
			TransactionID: transactionID,
			OrderItems:    orderItems,
		}
		ts = append(ts, t)
	}

	//continue to cancel
	for _, t := range ts {
		go func(oisOfTransaction []orderItem) {
			if len(ts) == 0 {
				return
			}
			if len(ois) == 1 {
				_, err := cancelTransactionSingleOrder(ois[0])
				if err != nil {
					util.Printf("cancel transaction error: %#v", err)
					return
				}
				return
			}
			_, err = cancelTransactionMultipleOrders(ois)
			if err != nil {
				util.Printf("cancel transaction error: %#v", err)
				return
			}
		}(t.OrderItems)
	}
	return nil
}
