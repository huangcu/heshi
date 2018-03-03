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
	Category      string  `json:"category"`
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"`
	ExtraInfo     string  `json:"extra_info"`
	SpecialNotice string  `json:"special_notice"`
	DownPayment   float64 `json:"downpayment"`
	BuyerID       string  `json:"buyer_id"`
	TransactionID string  `json:"-"`
	Status        string  `json:"status"`
	InStock       int     `json:"in_stock"`
}

//ALLOW TO EDIT PRICE,SPECIALNOTICE,DOWNPAYMENT,STATUS ONLY
func updateOrder(c *gin.Context) {
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
	priceStr := c.PostForm("price")
	if priceStr != "" {
		cValue, err := util.StringToFloat(priceStr)
		if err != nil {
			c.JSON(http.StatusOK, vemsgOrderPriceNotValid)
			return
		} else if cValue == 0 {
			c.JSON(http.StatusOK, vemsgOrderPriceNotValid)
			return
		} else {
			oi.Price = cValue
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
		switch item.Category {
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
	if quantity > oi.Quantity {
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
	if quantity > oi.Quantity {
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
	if quantity > oi.Quantity {
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
	switch item.Category {
	case DIAMOND:
		oq = fmt.Sprintf("UPDATE diamonds SET status='ORDERED', updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", item.ItemID)
	case JEWELRY:
		oq = fmt.Sprintf(`UPDATE jewelrys SET quantity='%d', updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND online='YES' AND quantity>='%d'`, item.InStock-item.Quantity, item.ItemID, item.Quantity)
	case GEM:
		oq = fmt.Sprintf(`UPDATE gems SET quantity='%d', updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND online='YES' AND quantity>='%d'`, item.InStock-item.Quantity, item.ItemID, item.Quantity)
	case SMALLDIAMOND:
		oq = fmt.Sprintf(`UPDATE small_diamonds SET quantity='%d', updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND online='YES' AND quantity>='%d'`, item.InStock-item.Quantity, item.ItemID, item.Quantity)
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
		switch item.Category {
		case DIAMOND:
			oq = fmt.Sprintf("UPDATE diamonds SET status='ORDERED', updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", item.ItemID)
		case JEWELRY:
			oq = fmt.Sprintf(`UPDATE jewelrys SET quantity='%d', updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND online='YES' AND quantity>='%d'`, item.InStock-item.Quantity, item.ItemID, item.Quantity)
		case GEM:
			oq = fmt.Sprintf(`UPDATE gems SET quantity='%d', updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND online='YES' AND quantity>='%d'`, item.InStock-item.Quantity, item.ItemID, item.Quantity)
		case SMALLDIAMOND:
			oq = fmt.Sprintf(`UPDATE small_diamonds SET quantity='%d', updated_at=(CURRENT_TIMESTAMP) 
		WHERE id='%s' AND online='YES' AND quantity>='%d'`, item.InStock-item.Quantity, item.ItemID, item.Quantity)
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
