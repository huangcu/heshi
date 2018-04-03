package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"heshi/errors"
	"net/http"
	"strings"
	"util"

	"github.com/gin-gonic/gin"
)

type transaction struct {
	TransactionID string      `json:"transaction_id"`
	OrderItems    []orderItem `json:"order_items"`
}

type orderItem struct {
	ID            string  `json:"order_id"`
	ItemID        string  `json:"item_id"`
	ItemCategory  string  `json:"item_category"`
	ItemQuantity  int     `json:"item_quantity"`
	ItemPrice     float64 `json:"item_price"`
	SoldPriceUSD  float64 `json:"sold_price_usd"`
	SoldPriceCNY  float64 `json:"sold_price_cny"`
	SoldPriceEUR  float64 `json:"sold_price_eur"`
	ReturnPoint   float64 `json:"return_point"`
	ChosenBy      string  `json:"chosen_by"`
	ExtraInfo     string  `json:"extra_info"`
	SpecialNotice string  `json:"special_notice"`
	DownPayment   float64 `json:"downpayment"`
	BuyerID       string  `json:"buyer_id"`
	TransactionID string  `json:"transaction_id"`
	Status        string  `json:"status"`
	InStock       int     `json:"in_stock"`
}

func getOrderDetail(c *gin.Context) {
	oid := c.Param("id")
	uid := c.MustGet("id").(string)
	q := fmt.Sprintf(`SELECT item_id, item_price, item_category, item_quantity, transaction_id, downpayment,
		sold_price_usd, sold_price_cny, sold_price_eur,return_point, chosen_by,
	 status, extra_info, special_notice FROM orders where id='%s' AND buyer_id='%s'`, oid, uid)

	var itemID, itemCategory, transactionID, status string
	var chosenBy, extraInfo, specialNotice sql.NullString
	var itemPrice, downpayment float64
	var itemQuantity int
	var soldPriceUSD, soldPriceCNY, soldPriceEUR, returnPoint sql.NullFloat64

	if err := dbQueryRow(q).Scan(&itemID, &itemPrice, &itemCategory, &itemQuantity, &transactionID,
		&downpayment, &soldPriceUSD, &soldPriceCNY, &soldPriceEUR, &returnPoint, &chosenBy,
		&status, &extraInfo, &specialNotice); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	oi := orderItem{
		ID:            oid,
		ItemID:        itemID,
		ItemCategory:  itemCategory,
		ItemQuantity:  itemQuantity,
		ItemPrice:     itemPrice,
		SoldPriceUSD:  soldPriceUSD.Float64,
		SoldPriceCNY:  soldPriceCNY.Float64,
		SoldPriceEUR:  soldPriceEUR.Float64,
		ReturnPoint:   returnPoint.Float64,
		ChosenBy:      chosenBy.String,
		ExtraInfo:     extraInfo.String,
		SpecialNotice: specialNotice.String,
		DownPayment:   downpayment,
		BuyerID:       uid,
		TransactionID: transactionID,
		Status:        status,
	}
	c.JSON(http.StatusOK, oi)
}

func getTransactionDetail(c *gin.Context) {
	tid := c.Param("id")
	uid := c.MustGet("id").(string)
	q := fmt.Sprintf(`SELECT id, item_id, item_price, item_category, item_quantity, downpayment,
		sold_price_usd, sold_price_cny, sold_price_eur,return_point, chosen_by, 
		status, extra_info, special_notice FROM orders where transaction_id='%s' AND buyer_id='%s'`, tid, uid)

	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	var ois []orderItem
	for rows.Next() {
		var id, itemID, itemCategory, transactionID, status string
		var chosenBy, extraInfo, specialNotice sql.NullString
		var itemPrice, downpayment float64
		var itemQuantity int
		var soldPriceUSD, soldPriceCNY, soldPriceEUR, returnPoint sql.NullFloat64
		if err := dbQueryRow(q).Scan(&id, &itemID, &itemPrice, &itemCategory, &itemQuantity, &transactionID,
			&downpayment, &soldPriceUSD, &soldPriceCNY, &soldPriceEUR, &returnPoint, &chosenBy,
			&status, &extraInfo, &specialNotice); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		oi := orderItem{
			ID:            id,
			ItemID:        itemID,
			ItemCategory:  itemCategory,
			ItemQuantity:  itemQuantity,
			ItemPrice:     itemPrice,
			SoldPriceUSD:  soldPriceUSD.Float64,
			SoldPriceCNY:  soldPriceCNY.Float64,
			SoldPriceEUR:  soldPriceEUR.Float64,
			ReturnPoint:   returnPoint.Float64,
			ChosenBy:      chosenBy.String,
			ExtraInfo:     extraInfo.String,
			SpecialNotice: specialNotice.String,
			DownPayment:   downpayment,
			BuyerID:       uid,
			TransactionID: tid,
			Status:        status,
		}
		ois = append(ois, oi)
	}
	c.JSON(http.StatusOK, transaction{
		TransactionID: tid,
		OrderItems:    ois,
	})
}

func getAllTransactionsOfUser(c *gin.Context) {
	//for admin, pass userid in query, for user, from login session
	uid := c.Param("id")
	if uid == "" {
		uid = c.MustGet("id").(string)
	}
	ts, err := getAllTransactionsOfAUser(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ts)
}

func getAllTransactionsRecommendedByAgent(c *gin.Context) {
	agentID := c.MustGet("id").(string)
	ids, err := getUsersIDRecommendedBy(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	for _, id := range ids {
		// TODO instead of one by one,maybe user buyer_id in ("", ""), select all
		// 1 big query against lots of small query
		getAllTransactionsOfAUser(id)
	}
}

func getAllTransactionsOfAUser(buyerID string) ([]transaction, error) {
	q := fmt.Sprintf(`SELECT transaction_id FROM orders WHERE buyer_id='%s'`, buyerID)
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactionIDs []string
	for rows.Next() {
		var transactionID string
		if err := rows.Scan(&transactionID); err != nil {
			return nil, err
		}
		transactionIDs = append(transactionIDs, transactionID)
	}
	var ts []transaction
	for _, transactionID := range transactionIDs {
		q := fmt.Sprintf(`SELECT id, item_id, item_price, item_category, item_quantity, downpayment,
			sold_price_usd, sold_price_cny, sold_price_eur,return_point, chosen_by,
			status, extra_info, special_notice FROM orders WHERE buyer_id='%s' AND transcation_id='%s'`, buyerID, transactionID)
		rows, err := dbQuery(q)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var ois []orderItem
		for rows.Next() {
			var id, itemID, itemCategory, status string
			var chosenBy, extraInfo, specialNotice sql.NullString
			var itemPrice, downpayment float64
			var itemQuantity int
			var soldPriceUSD, soldPriceCNY, soldPriceEUR, returnPoint sql.NullFloat64

			if err := rows.Scan(&id, &itemID, &itemPrice, &itemCategory, &itemQuantity,
				&downpayment, &soldPriceUSD, &soldPriceCNY, &soldPriceEUR, &returnPoint, &chosenBy,
				&status, &extraInfo, &specialNotice); err != nil {
				return nil, err
			}
			oi := orderItem{
				ID:            id,
				ItemID:        itemID,
				ItemCategory:  itemCategory,
				ItemQuantity:  itemQuantity,
				ItemPrice:     itemPrice,
				SoldPriceUSD:  soldPriceUSD.Float64,
				SoldPriceCNY:  soldPriceCNY.Float64,
				SoldPriceEUR:  soldPriceEUR.Float64,
				ReturnPoint:   returnPoint.Float64,
				ChosenBy:      chosenBy.String,
				ExtraInfo:     extraInfo.String,
				SpecialNotice: specialNotice.String,
				DownPayment:   downpayment,
				BuyerID:       buyerID,
				TransactionID: transactionID,
				Status:        status,
			}
			ois = append(ois, oi)
		}
		t := transaction{
			TransactionID: transactionID,
			OrderItems:    ois,
		}
		ts = append(ts, t)
	}
	return ts, nil
}

//ADMIN only
func getAllTransactions(c *gin.Context) {
	rows, err := dbQuery("SELECT transaction_id FROM orders")
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	var transactionIDs []string
	for rows.Next() {
		var transactionID string
		if err := rows.Scan(&transactionID); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		transactionIDs = append(transactionIDs, transactionID)
	}
	var ts []transaction
	for _, transactionID := range transactionIDs {
		q := fmt.Sprintf(`SELECT id, item_id, item_price, item_category, item_quantity, downpayment,
			sold_price_usd, sold_price_cny, sold_price_eur,return_point, chosen_by,
			buyer_id, status, extra_info, special_notice FROM orders WHERE transaction_id='%s'`, transactionID)
		rows, err := dbQuery(q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		defer rows.Close()

		var ois []orderItem
		for rows.Next() {
			var id, itemID, itemCategory, buyerID, status string
			var chosenBy, extraInfo, specialNotice sql.NullString
			var itemPrice, downpayment float64
			var itemQuantity int
			var soldPriceUSD, soldPriceCNY, soldPriceEUR, returnPoint sql.NullFloat64

			if err := rows.Scan(&id, &itemID, &itemPrice, &itemCategory, &itemQuantity,
				&downpayment, &soldPriceUSD, &soldPriceCNY, &soldPriceEUR, &returnPoint, &chosenBy,
				&buyerID, &status, &extraInfo, &specialNotice); err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			oi := orderItem{
				ID:            id,
				ItemID:        itemID,
				ItemCategory:  itemCategory,
				ItemQuantity:  itemQuantity,
				ItemPrice:     itemPrice,
				SoldPriceUSD:  soldPriceUSD.Float64,
				SoldPriceCNY:  soldPriceCNY.Float64,
				SoldPriceEUR:  soldPriceEUR.Float64,
				ReturnPoint:   returnPoint.Float64,
				ExtraInfo:     extraInfo.String,
				SpecialNotice: specialNotice.String,
				DownPayment:   downpayment,
				BuyerID:       buyerID,
				TransactionID: transactionID,
				Status:        status,
			}
			ois = append(ois, oi)
		}
		t := transaction{
			TransactionID: transactionID,
			OrderItems:    ois,
		}
		ts = append(ts, t)
	}
	c.JSON(http.StatusOK, ts)
}

//AGENT & ADMIN ALLOW TO EDIT SOLD_PRICE_USD/CNY/EUR,SPECIALNOTICE,DOWNPAYMENT,STATUS ONLY
func updateOrder(c *gin.Context) {
	uid := c.MustGet("id").(string)
	oid := c.Param("id")
	if exist, err := isOrderExistByID(oid); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if !exist {
		c.JSON(http.StatusBadRequest, "Order doesn't exist")
		return
	}
	oi := orderItem{
		ID:            oid,
		ExtraInfo:     c.PostForm("extra_info"),
		SpecialNotice: c.PostForm("special_notice"),
		Status:        strings.ToUpper(c.PostForm("status")),
	}
	priceUSDStr := c.PostForm("sold_price_usd")
	if priceUSDStr != "" {
		cValue, err := util.StringToFloat(priceUSDStr)
		if err != nil {
			c.JSON(http.StatusOK, vemsgOrderPriceNotValid)
			return
		} else if cValue == 0 {
			c.JSON(http.StatusOK, vemsgOrderPriceNotValid)
			return
		} else {
			oi.SoldPriceCNY = cValue
		}
	}
	priceCNYStr := c.PostForm("sold_price_cny")
	if priceCNYStr != "" {
		cValue, err := util.StringToFloat(priceCNYStr)
		if err != nil {
			c.JSON(http.StatusOK, vemsgOrderPriceNotValid)
			return
		} else if cValue == 0 {
			c.JSON(http.StatusOK, vemsgOrderPriceNotValid)
			return
		} else {
			oi.SoldPriceCNY = cValue
		}
	}
	priceEURStr := c.PostForm("sold_price_eur")
	if priceEURStr != "" {
		cValue, err := util.StringToFloat(priceEURStr)
		if err != nil {
			c.JSON(http.StatusOK, vemsgOrderPriceNotValid)
			return
		} else if cValue == 0 {
			c.JSON(http.StatusOK, vemsgOrderPriceNotValid)
			return
		} else {
			oi.SoldPriceEUR = cValue
		}
	}
	downPaymentStr := c.PostForm("downpayment")
	if downPaymentStr != "" {
		cValue, err := util.StringToFloat(downPaymentStr)
		if err != nil {
			c.JSON(http.StatusOK, vemsgOrderDownPaymentNotValid)
			return
		} else if cValue == 0 {
			c.JSON(http.StatusOK, vemsgOrderDownPaymentNotValid)
			return
		} else {
			oi.DownPayment = cValue
		}
	}
	validStatus := []string{"ORDERED", "CANCELLED", "SOLD", "DOWNPAYMENT"}
	if !util.IsInArrayString(oi.Status, validStatus) {
		c.JSON(http.StatusOK, vemsgOrderStatusNotValid)
		return
	}
	q := oi.composeUpdateQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, oi.ID)
	go newHistoryRecords(uid, "orders", oi.ID, oi.parmsKV())
}

func createOrder(c *gin.Context) {
	//check if all items still available
	orderItems := make([]orderItem, 0)
	json.Unmarshal([]byte(c.PostForm("items")), &orderItems)
	if err := checkItems(orderItems); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	for _, item := range orderItems {
		if item.Status != "AVAILABLE" {
			c.JSON(http.StatusOK, orderItems)
			return
		}
	}

	//continue to order
	if len(orderItems) == 1 {
		t, err := orderSingleItem(orderItems[0])
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, t)
	} else {
		t, err := orderMultipleItems(orderItems)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, t)
	}
}

func cartItems(c *gin.Context) {
	items := make([]orderItem, 0)
	json.Unmarshal([]byte(c.PostForm("items")), &items)
	if err := checkItems(items); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, items)
}

func checkItems(items []orderItem) error {
	for _, item := range items {
		switch item.ItemCategory {
		case DIAMOND:
			if err := item.checkDiamondItem(); err != nil {
				return err
			}
		case JEWELRY:
			if err := item.checkJewelryItem(); err != nil {
				return err
			}
		case GEM:
			if err := item.checkGemItem(); err != nil {
				return err
			}
		case SMALLDIAMOND:
			if err := item.checkSmallDiamondItem(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (oi *orderItem) checkDiamondItem() error {
	var status string
	if err := dbQueryRow(fmt.Sprintf("SELECT status FROM diamonds WHERE id='%s' AND status='AVAILABLE'", oi.ItemID)).Scan(&status); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		oi.Status = "NOT AVAILABLE"
	}
	oi.Status = "AVAILABLE"
	return nil
}

func (oi *orderItem) checkJewelryItem() error {
	var quantity int
	if err := dbQueryRow(fmt.Sprintf("SELECT quantity FROM jewelrys WHERE id='%s' AND online='YES'", oi.ItemID)).Scan(&quantity); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		oi.Status = "NOT AVAILABLE"
	}
	oi.InStock = quantity
	if quantity > oi.ItemQuantity {
		oi.Status = "AVAILABLE"
	} else {
		oi.Status = "STOCK_NOT_ENOUGH"
	}
	return nil
}

func (oi *orderItem) checkGemItem() error {
	var quantity int
	if err := dbQueryRow(fmt.Sprintf("SELECT quantity FROM gems WHERE id='%s' AND online='YES'", oi.ItemID)).Scan(&quantity); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		oi.Status = "NOT AVAILABLE"
	}
	oi.InStock = quantity
	if quantity > oi.ItemQuantity {
		oi.Status = "AVAILABLE"
	} else {
		oi.Status = "STOCK_NOT_ENOUGH"
	}
	return nil
}

func (oi *orderItem) checkSmallDiamondItem() error {
	var quantity int
	if err := dbQueryRow(fmt.Sprintf("SELECT quantity FROM small_diamonds WHERE id='%s' AND online='YES'", oi.ItemID)).Scan(&quantity); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		oi.Status = "NOT AVAILABLE"
	}
	oi.InStock = quantity
	if quantity > oi.ItemQuantity {
		oi.Status = "AVAILABLE"
	} else {
		oi.Status = "STOCK_NOT_ENOUGH"
	}
	return nil
}

func orderSingleItem(item orderItem) (*transaction, error) {
	item.ID = newV4()
	item.TransactionID = item.ID
	var oq string
	switch item.ItemCategory {
	case DIAMOND:
		oq = fmt.Sprintf("UPDATE diamonds SET status='ORDERED', updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", item.ItemID)
	case JEWELRY:
		oq = fmt.Sprintf(`UPDATE jewelrys SET quantity='%d', updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND online='YES' AND quantity>='%d'`, item.InStock-item.ItemQuantity, item.ItemID, item.ItemQuantity)
	case GEM:
		oq = fmt.Sprintf(`UPDATE gems SET quantity='%d', updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND online='YES' AND quantity>='%d'`, item.InStock-item.ItemQuantity, item.ItemID, item.ItemQuantity)
	case SMALLDIAMOND:
		oq = fmt.Sprintf(`UPDATE small_diamonds SET quantity='%d', updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND online='YES' AND quantity>='%d'`, item.InStock-item.ItemQuantity, item.ItemID, item.ItemQuantity)
	}
	err := dbTransact(db, func(tx *sql.Tx) error {
		traceSQL(oq)
		result, err := tx.Exec(oq)
		if err != nil {
			return err
		}
		if r, err := result.RowsAffected(); err != nil {
			return err
		} else if r != 1 {
			return errors.Newf("Item %s not AVAILABLE any more", item.ItemID)
		}
		q := item.composeInsertQuery()
		traceSQL(q)
		if _, err := tx.Exec(q); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	t := &transaction{
		TransactionID: item.TransactionID,
		OrderItems:    []orderItem{item},
	}
	return t, nil
}

func orderMultipleItems(items []orderItem) (*transaction, error) {
	qs := make(map[string]orderItem)
	transactionID := newV4()
	for _, item := range items {
		var oq string
		item.TransactionID = transactionID
		item.ID = newV4()
		switch item.ItemCategory {
		case DIAMOND:
			oq = fmt.Sprintf("UPDATE diamonds SET status='ORDERED', updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", item.ItemID)
		case JEWELRY:
			oq = fmt.Sprintf(`UPDATE jewelrys SET quantity='%d', updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND online='YES' AND quantity>='%d'`, item.InStock-item.ItemQuantity, item.ItemID, item.ItemQuantity)
		case GEM:
			oq = fmt.Sprintf(`UPDATE gems SET quantity='%d', updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND online='YES' AND quantity>='%d'`, item.InStock-item.ItemQuantity, item.ItemID, item.ItemQuantity)
		case SMALLDIAMOND:
			oq = fmt.Sprintf(`UPDATE small_diamonds SET quantity='%d', updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND online='YES' AND quantity>='%d'`, item.InStock-item.ItemQuantity, item.ItemID, item.ItemQuantity)
		}
		qs[oq] = item
	}

	err := dbTransact(db, func(tx *sql.Tx) error {
		for oq, item := range qs {
			traceSQL(oq)
			result, err := tx.Exec(oq)
			if err != nil {
				return err
			}
			if r, err := result.RowsAffected(); err != nil {
				return err
			} else if r != 1 {
				return errors.Newf("Item %s not AVAILABLE any more", item.ItemID)
			}
			tq := item.composeInsertQuery()
			traceSQL(tq)
			if _, err := tx.Exec(tq); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	t := &transaction{
		TransactionID: transactionID,
		OrderItems:    items,
	}
	return t, nil
}
