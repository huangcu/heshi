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
	TransactionID string       `json:"transaction_id"`
	OrderItems    []*orderItem `json:"order_items"`
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

type orderResult struct {
	ItemID  string `json:"item_id"`
	Message string `json:"message"`
}

func getOrderDetail(c *gin.Context) {
	oid := c.Param("id")
	q := fmt.Sprintf(`SELECT id, item_id, item_price, item_category, item_quantity, downpayment, 
		buyer_id, transaction_id,	sold_price_usd, sold_price_cny, sold_price_eur,
		return_point, chosen_by, status, extra_info, special_notice 
		FROM orders 
		WHERE id='%s'`, oid)
	userType := c.MustGet("usertype").(string)
	if userType == CUSTOMER {
		buyerID := c.MustGet("id").(string)
		q = fmt.Sprintf("%s AND buyer_id='%s'", q, buyerID)
	}
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ois, err := composeOrders(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if len(ois) == 1 {
		c.JSON(http.StatusOK, ois[0])
		return
	}

	c.JSON(http.StatusOK, ois[0])
}

func getOrderDetailOfUserRecommendedByAgent(c *gin.Context) {
	oid := c.Param("id")
	agentID := c.MustGet("id").(string)
	q := fmt.Sprintf(`SELECT id, item_id, item_price, item_category, item_quantity, downpayment, 
		buyer_id, transaction_id,	sold_price_usd, sold_price_cny, sold_price_eur,
		return_point, chosen_by, status, extra_info, special_notice 
	  FROM orders 
	  WHERE id='%s' 
	  AND buyer_id IN (SELECT id FROM users WHERE status='ACTIVE' AND recommended_by='%s')`, oid, agentID)
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	ois, err := composeOrders(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if len(ois) == 1 {
		c.JSON(http.StatusOK, ois[0])
		return
	}
	c.JSON(http.StatusBadRequest, "NOT EXIST OR NOT ALLOWED")
}

func getTransactionDetail(c *gin.Context) {
	tid := c.Param("id")
	q := fmt.Sprintf(`SELECT id, item_id, item_price, item_category, item_quantity, downpayment, 
		buyer_id, transaction_id,	sold_price_usd, sold_price_cny, sold_price_eur,
		return_point, chosen_by, status, extra_info, special_notice 
		FROM orders 
		WHERE transaction_id='%s'`, tid)
	userType := c.MustGet("usertype").(string)
	if userType == CUSTOMER {
		buyerID := c.MustGet("id").(string)
		q = fmt.Sprintf("%s AND buyer_id='%s'", q, buyerID)
	}
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()
	ois, err := composeOrders(rows)
	c.JSON(http.StatusOK, transaction{
		TransactionID: tid,
		OrderItems:    ois,
	})
}

func getTransactionDetailOfUserRecommendedByAgent(c *gin.Context) {
	tid := c.Param("id")
	agentID := c.MustGet("id").(string)
	q := fmt.Sprintf(`SELECT id, item_id, item_price, item_category, item_quantity, downpayment, 
		buyer_id, transaction_id,	sold_price_usd, sold_price_cny, sold_price_eur,
		return_point, chosen_by, status, extra_info, special_notice 
		FROM orders 
		WHERE transaction_id='%s' 
		AND buyer_id IN (SELECT id FROM users WHERE status='ACTIVE' AND recommended_by='%s')`, tid, agentID)

	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()
	ois, err := composeOrders(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, transaction{
		TransactionID: tid,
		OrderItems:    ois,
	})
}

func getAllTransactionsOfAUser(c *gin.Context) {
	//for admin, pass userid in query, for user, from login session
	buyerID := c.Param("id")
	if buyerID == "" {
		buyerID = c.MustGet("id").(string)
	}
	q := fmt.Sprintf(`SELECT id, item_id, item_price, item_category, item_quantity, downpayment,
			buyer_id, transaction_id, sold_price_usd, sold_price_cny, sold_price_eur,return_point, 
			chosen_by, status, extra_info, special_notice 
			FROM orders 
			WHERE buyer_id='%s'`, buyerID)
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ts, err := composeTransactions(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ts)
}

func getAllTransactionsOfAUserRecommendedByAgent(c *gin.Context) {
	userID := c.Param("id")
	agentID := c.MustGet("id").(string)
	q := fmt.Sprintf(`SELECT id, item_id, item_price, item_category, item_quantity, downpayment,
			buyer_id, transaction_id, sold_price_usd, sold_price_cny, sold_price_eur,return_point, 
			chosen_by, status, extra_info, special_notice 
			FROM orders 
			WHERE buyer_id='%s'
			AND buyer_id IN (SELECT id FROM users WHERE status='ACTIVE' AND recommended_by='%s')`, userID, agentID)
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ts, err := composeTransactions(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ts)
}

func getAllTransactionsOfUserRecommendedByAgent(c *gin.Context) {
	agentID := c.MustGet("id").(string)
	q := fmt.Sprintf(`SELECT id, item_id, item_price, item_category, item_quantity, downpayment,
			buyer_id, transaction_id, sold_price_usd, sold_price_cny, sold_price_eur,return_point, 
			chosen_by, status, extra_info, special_notice 
			FROM orders 
			WHERE buyer_id IN (SELECT id FROM users WHERE status='ACTIVE' AND recommended_by='%s')`, agentID)
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ts, err := composeTransactions(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ts)
}

//ADMIN only
func getAllTransactions(c *gin.Context) {
	q := `SELECT id, item_id, item_price, item_category, item_quantity, downpayment,
			buyer_id, transaction_id, sold_price_usd, sold_price_cny, sold_price_eur, return_point, 
			chosen_by, status, extra_info, special_notice 
			FROM orders`
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ts, err := composeTransactions(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ts)
}

// TODO
// 1. what can be updated
// 2. input data format
// 3. sperate api??? status, price etc (extra info, return_point)
func updateTransaction(c *gin.Context) {
	// updatedBy := c.MustGet("id").(string)
	orderItems := make([]*orderItem, 0)
	if err := json.Unmarshal([]byte(c.PostForm("items")), &orderItems); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

}

//TODO agent discount 9.9/9 --- base on current sold_price????, only allowed to change sold_price???
//ADMIN ALLOW TO EDIT SOLD_PRICE_USD/CNY/EUR,SPECIALNOTICE,DOWNPAYMENT,STATUS ONLY
// TODO order status change?? one by one if downpayment ---
func updateOrder(c *gin.Context) {
	uid := c.MustGet("id").(string)
	oid := c.Param("id")
	oiInDB, err := getOrderByID(oid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if oiInDB == nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Order %s doesn't exist", oid))
		return
	}
	oi := orderItem{
		ID:            oid,
		ExtraInfo:     c.PostForm("extra_info"),
		SpecialNotice: c.PostForm("special_notice"),
		Status:        orderStatusAToOrderStatusM(strings.ToUpper(c.PostForm("status"))),
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

	if oi.Status != "" {
		if !util.IsInArrayString(oi.Status, VALID_ORDER_STATUS_M) {
			c.JSON(http.StatusOK, vemsgOrderStatusNotValid)
			return
		}

		if oi.Status == CANCELLED || oi.Status == MCANCELLED {
			vemsgOrderStatusNotValid.Message = "Not allowed to cancel order, please call cancel API"
			c.JSON(http.StatusOK, vemsgOrderStatusNotValid)
			return
		}
		if !isOrderStatusChangeAllowed(oiInDB.Status, oi.Status) {
			vemsgOrderStatusNotValid.Message = fmt.Sprintf("Cannot change status from %s to %s", oiInDB.Status, oi.Status)
			c.JSON(http.StatusOK, vemsgOrderStatusNotValid)
			return
		}
	}
	if oi.Status != "" && strings.ToUpper(oiInDB.ItemCategory) == DIAMOND {
		err := dbTransact(db, func(tx *sql.Tx) error {
			q := oi.composeUpdateQuery()
			traceSQL(q)
			result, err := tx.Exec(q)
			if err != nil {
				return err
			}
			if r, err := result.RowsAffected(); err != nil {
				return err
			} else if r != 1 {
				return nil
			}
			tq := fmt.Sprintf(`UPDATE diamonds SET status='%s' WHERE id='%s'`, oi.Status, oiInDB.ItemID)
			traceSQL(tq)
			result, err = tx.Exec(tq)
			if err != nil {
				return err
			}
			r, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if r == 1 {
				oStateMap := make(map[string]interface{})
				oStateMap["status"] = oi.Status + ",Due to Order: " + oid + " Status Change"
				//diamonds status changed
				go newHistoryRecords(uid, "diamonds", oiInDB.ItemID, oStateMap)
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	} else {
		q := oi.composeUpdateQuery()
		if _, err := dbExec(q); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	}

	c.JSON(http.StatusOK, "SUCCESS")
	go newHistoryRecords(uid, "orders", oi.ID, oi.parmsKV())
}

func createOrder(c *gin.Context) {
	//check if all items still available
	buyerID := c.MustGet("id").(string)
	orderItems := make([]*orderItem, 0)
	if err := json.Unmarshal([]byte(c.PostForm("items")), &orderItems); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	if err := checkItems(orderItems); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	for _, item := range orderItems {
		fmt.Println(item.Status)
		if item.Status != AVAILABLE {
			c.JSON(http.StatusBadRequest, orderResult{item.ItemID, item.Status})
			return
		}
	}

	//continue to order
	for _, item := range orderItems {
		item.BuyerID = buyerID
	}
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

// TODO only allowed cancel whole transacton, not allowed to cancel single
// WHO is allowed to cancel order? agent???
// WHAT ORDERS ARE ALLOWED TO CANCEL --
func cancelTransaction(c *gin.Context) {
	tid := c.Param("id")
	if exist, err := isTransactionExistByID(tid); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if !exist {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("No such transaction %s", tid))
		return
	}

	q := fmt.Sprintf(`SELECT id, item_id, item_category, item_quantity  
		FROM orders 
		WHERE transaction_id='%s' 
		AND status NOT IN ('%s','%s','%s','%s','%s','%s')`,
		tid, CANCELLED, MCANCELLED, DELIVERED, MDELIVERED, RECEIVED, MRECEIVED)
	adminOrUserID := c.MustGet("id").(string)
	userType := c.MustGet("usertype").(string)
	status := CANCELLED
	if userType == ADMIN {
		status = MCANCELLED
	} else {
		q = fmt.Sprintf("%s AND buyer_id='%s'", q, adminOrUserID)
	}
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	var ois []*orderItem
	for rows.Next() {
		var id, itemID, itemCategory string
		var itemQuantity int
		if err := rows.Scan(&id, &itemID, &itemCategory, &itemQuantity); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		oi := &orderItem{
			ID:           id,
			ItemID:       itemID,
			ItemCategory: itemCategory,
			ItemQuantity: itemQuantity,
			Status:       status,
		}
		ois = append(ois, oi)
	}
	if len(ois) == 0 {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Transaction %s is not allowed to cancel Or not exist", tid))
		return
	}

	//continue to cancel
	if len(ois) == 1 {
		t, err := cancelTransactionSingleOrder(adminOrUserID, ois[0])
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, t)
	} else {
		t, err := cancelTransactionMultipleOrders(adminOrUserID, ois)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, t)
	}

	if userType == ADMIN {
		go func(uid string) {
			for _, oi := range ois {
				kv := make(map[string]interface{})
				kv["status"] = CANCELLED
				go newHistoryRecords(uid, "orders", oi.ID, kv)
			}
		}(adminOrUserID)
	}
}

func cartItems(c *gin.Context) {
	items := make([]*orderItem, 0)
	if err := json.Unmarshal([]byte(c.PostForm("items")), &items); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}
	if err := checkItems(items); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, items)
}

func checkItems(items []*orderItem) error {
	for _, item := range items {
		switch strings.ToUpper(item.ItemCategory) {
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
	var priceRetail float64
	q := fmt.Sprintf("SELECT status,price_retail FROM diamonds WHERE id='%s' AND status='AVAILABLE'", oi.ItemID)
	if err := dbQueryRow(q).Scan(&status, &priceRetail); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		oi.Status = "NOT AVAILABLE"
		return nil
	}
	oi.ItemPrice = priceRetail
	oi.Status = AVAILABLE
	return nil
}

func (oi *orderItem) checkJewelryItem() error {
	var quantity int
	var price float64
	q := fmt.Sprintf("SELECT stock_quantity,price FROM jewelrys WHERE id='%s' AND status='AVAILABLE'", oi.ItemID)
	if err := dbQueryRow(q).Scan(&quantity, &price); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		oi.Status = "NOT AVAILABLE"
		return nil
	}
	oi.InStock = quantity
	if quantity > oi.ItemQuantity {
		oi.Status = AVAILABLE
	} else {
		oi.Status = "STOCK_NOT_ENOUGH"
	}
	oi.ItemPrice = price
	return nil
}

func (oi *orderItem) checkGemItem() error {
	var quantity int
	var price float64

	q := fmt.Sprintf("SELECT stock_quantity,price FROM gems WHERE id='%s' AND status='AVAILABLE'", oi.ItemID)
	if err := dbQueryRow(q).Scan(&quantity, &price); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		oi.Status = "NOT AVAILABLE"
		return nil
	}
	oi.InStock = quantity
	if quantity > oi.ItemQuantity {
		oi.Status = AVAILABLE
	} else {
		oi.Status = "STOCK_NOT_ENOUGH"
	}
	oi.ItemPrice = price
	return nil
}

func (oi *orderItem) checkSmallDiamondItem() error {
	var quantity int
	if err := dbQueryRow(fmt.Sprintf("SELECT stock_quantity FROM small_diamonds WHERE id='%s' AND status='AVAILABLE'", oi.ItemID)).Scan(&quantity); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		oi.Status = "NOT AVAILABLE"
		return nil
	}
	oi.InStock = quantity
	if quantity > oi.ItemQuantity {
		oi.Status = AVAILABLE
	} else {
		oi.Status = "STOCK_NOT_ENOUGH"
	}
	return nil
}

func orderSingleItem(item *orderItem) (*transaction, error) {
	item.ID = newV4()
	item.TransactionID = item.ID
	var oq string
	switch strings.ToUpper(item.ItemCategory) {
	case DIAMOND:
		oq = fmt.Sprintf("UPDATE diamonds SET status='ORDERED', updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", item.ItemID)
	case JEWELRY:
		oq = fmt.Sprintf(`UPDATE jewelrys SET stock_quantity=stock_quantity-%d, updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND status='AVAILABLE' AND stock_quantity>='%d'`, item.ItemQuantity, item.ItemID, item.ItemQuantity)
	case GEM:
		oq = fmt.Sprintf(`UPDATE gems SET stock_quantity=stock_quantity-%d, updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND status='AVAILABLE' AND stock_quantity>='%d'`, item.ItemQuantity, item.ItemID, item.ItemQuantity)
	case SMALLDIAMOND:
		oq = fmt.Sprintf(`UPDATE small_diamonds SET stock_quantity=stock_quantity-%d, updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND stock_quantity>='%d'`, item.ItemQuantity, item.ItemID, item.ItemQuantity)
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
		item.Status = ORDERED
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
		OrderItems:    []*orderItem{item},
	}
	return t, nil
}

func orderMultipleItems(items []*orderItem) (*transaction, error) {
	qs := make(map[string]*orderItem)
	transactionID := newV4()
	for _, item := range items {
		var oq string
		item.TransactionID = transactionID
		item.ID = newV4()
		switch strings.ToUpper(item.ItemCategory) {
		case DIAMOND:
			oq = fmt.Sprintf("UPDATE diamonds SET status='ORDERED', updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", item.ItemID)
		case JEWELRY:
			oq = fmt.Sprintf(`UPDATE jewelrys SET stock_quantity=stock_quantity-%d, updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND status='AVAILABLE' AND stock_quantity>='%d'`, item.ItemQuantity, item.ItemID, item.ItemQuantity)
		case GEM:
			oq = fmt.Sprintf(`UPDATE gems SET stock_quantity=stock_quantity-%d, updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND status='AVAILABLE' AND stock_quantity>='%d'`, item.ItemQuantity, item.ItemID, item.ItemQuantity)
		case SMALLDIAMOND:
			oq = fmt.Sprintf(`UPDATE small_diamonds SET stock_quantity=stock_quantity-%d, updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND stock_quantity>='%d'`, item.ItemQuantity, item.ItemID, item.ItemQuantity)
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
			item.Status = ORDERED
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

func cancelTransactionSingleOrder(uid string, item *orderItem) (*transaction, error) {
	var oq string
	var table string
	changeMap := make(map[string]interface{})
	switch strings.ToUpper(item.ItemCategory) {
	case DIAMOND:
		oq = fmt.Sprintf("UPDATE diamonds SET status='AVAILABLE', updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", item.ItemID)
		table = "diamonds"
		changeMap["status"] = CANCELLED
	case JEWELRY:
		oq = fmt.Sprintf(`UPDATE jewelrys SET stock_quantity= stock_quantity + %d, updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s'`, item.ItemQuantity, item.ItemID)
		table = "jewelrys"
		changeMap["stock_quantity"] = fmt.Sprintf("+%d, Due to Order: '%s' cancel", item.ItemQuantity, item.ID)
	case GEM:
		oq = fmt.Sprintf(`UPDATE gems SET stock_quantity= stock_quantity+ %d, updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s'`, item.ItemQuantity, item.ItemID)
		table = "gems"
		changeMap["stock_quantity"] = fmt.Sprintf("+%d, Due to Order: '%s' cancel", item.ItemQuantity, item.ID)
	case SMALLDIAMOND:
		oq = fmt.Sprintf(`UPDATE small_diamonds SET stock_quantity=stock_quantity+ %d, updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s'`, item.ItemQuantity, item.ItemID)
		table = "small_diamonds"
		changeMap["stock_quantity"] = fmt.Sprintf("+%d, Due to Order: '%s' cancel", item.ItemQuantity, item.ID)
	}
	err := dbTransact(db, func(tx *sql.Tx) error {
		traceSQL(oq)
		result, err := tx.Exec(oq)
		if err != nil {
			return err
		}
		r, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if r != 1 {
			// should ingore, item to be canceled shouldn't be not AVAILABLE, if it is, return cancelled
			fmt.Printf("Cancel Transaction: Item %s not AVAILABLE any more", item.ItemID)
		} else {
			// TODO should track product change due to order cancel???
			go newHistoryRecords(uid, table, item.ItemID, changeMap)
		}

		q := item.composeUpdateQuery()
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
		OrderItems:    []*orderItem{item},
	}
	return t, nil
}

func cancelTransactionMultipleOrders(uid string, items []*orderItem) (*transaction, error) {
	qs := make(map[string]*orderItem)
	for _, item := range items {
		var oq string
		switch strings.ToUpper(item.ItemCategory) {
		case DIAMOND:
			oq = fmt.Sprintf("UPDATE diamonds SET status='AVAILABLE', updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", item.ItemID)
		case JEWELRY:
			oq = fmt.Sprintf(`UPDATE jewelrys SET stock_quantity=stock_quantity + %d, updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s'`, item.ItemQuantity, item.ItemID)
		case GEM:
			oq = fmt.Sprintf(`UPDATE gems SET stock_quantity=stock_quantity + %d, updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s'`, item.ItemQuantity, item.ItemID)
		case SMALLDIAMOND:
			oq = fmt.Sprintf(`UPDATE small_diamonds SET stock_quantity=stock_quantity + %d, updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s'`, item.ItemQuantity, item.ItemID)
		}
		qs[oq] = item
	}
	err := dbTransact(db, func(tx *sql.Tx) error {
		for oq, item := range qs {
			var table string
			changeMap := make(map[string]interface{})
			switch strings.ToUpper(item.ItemCategory) {
			case DIAMOND:
				table = "diamonds"
				changeMap["status"] = CANCELLED
			case JEWELRY:
				table = "jewelrys"
				changeMap["stock_quantity"] = fmt.Sprintf("+%d, Due to Order: '%s' cancel", item.ItemQuantity, item.ID)
			case GEM:
				table = "gems"
				changeMap["stock_quantity"] = fmt.Sprintf("+%d, Due to Order: '%s' cancel", item.ItemQuantity, item.ID)
			case SMALLDIAMOND:
				table = "small_diamonds"
				changeMap["stock_quantity"] = fmt.Sprintf("+%d, Due to Order: '%s' cancel", item.ItemQuantity, item.ID)
			}
			traceSQL(oq)
			result, err := tx.Exec(oq)
			if err != nil {
				return err
			}
			r, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if r != 1 {
				// should ingore, item to be canceled shouldn't be not AVAILABLE, if it is, return cancelled
				fmt.Printf("Cancel Transaction: Item %s not AVAILABLE any more", item.ItemID)
			} else {
				// TODO should track product change due to order cancel???
				go newHistoryRecords(uid, table, item.ItemID, changeMap)
			}
			tq := item.composeUpdateQuery()
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
		TransactionID: items[0].TransactionID,
		OrderItems:    items,
	}
	return t, nil
}

func composeOrders(rows *sql.Rows) ([]*orderItem, error) {
	var ois []*orderItem
	for rows.Next() {
		var id, itemID, itemCategory, status, transactionID, buyerID string
		var chosenBy, extraInfo, specialNotice sql.NullString
		var itemPrice float64
		var itemQuantity int
		var downpayment, soldPriceUSD, soldPriceCNY, soldPriceEUR, returnPoint sql.NullFloat64

		if err := rows.Scan(&id, &itemID, &itemPrice, &itemCategory, &itemQuantity,
			&downpayment, &buyerID, &transactionID, &soldPriceUSD, &soldPriceCNY, &soldPriceEUR,
			&returnPoint, &chosenBy, &status, &extraInfo, &specialNotice); err != nil {
			return nil, err
		}
		oi := &orderItem{
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
			DownPayment:   downpayment.Float64,
			BuyerID:       buyerID,
			TransactionID: transactionID,
			Status:        status,
		}
		ois = append(ois, oi)
	}
	return ois, nil
}

func composeTransactions(rows *sql.Rows) ([]transaction, error) {
	var ois []*orderItem
	for rows.Next() {
		var id, itemID, itemCategory, status, transactionID, buyerID string
		var chosenBy, extraInfo, specialNotice sql.NullString
		var itemPrice float64
		var itemQuantity int
		var downpayment, soldPriceUSD, soldPriceCNY, soldPriceEUR, returnPoint sql.NullFloat64

		if err := rows.Scan(&id, &itemID, &itemPrice, &itemCategory, &itemQuantity,
			&downpayment, &buyerID, &transactionID, &soldPriceUSD, &soldPriceCNY, &soldPriceEUR,
			&returnPoint, &chosenBy, &status, &extraInfo, &specialNotice); err != nil {
			return nil, err
		}
		oi := &orderItem{
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
			DownPayment:   downpayment.Float64,
			BuyerID:       buyerID,
			TransactionID: transactionID,
			Status:        status,
		}
		ois = append(ois, oi)
	}

	transationOrderItemmap := make(map[string][]*orderItem)
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
	return ts, nil
}
