package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"heshi/errors"
	"net/http"

	"github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

// CREATE TABLE IF NOT EXISTS orders
// (
// 	id VARCHAR(225) PRIMARY KEY NOT NULL,
// 	transaction_id VARCHAR(225) NOT NULL,
// 	item_id INT NOT NULL,
// 	item_price FLOAT NOT NULL,
// 	item_category INT NOT NULL,
// 	buyer_id TINYINT(4) NOT NULL,
// 	downpayment FLOAT,
// 	status VARCHAR(20) NOT NULL DEFAULT 'ORDERED',
// 	extra_info VARCHAR(225),
// 	special_notice VARCHAR(225),
// 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
// 	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
// ) ENGINE=INNODB;

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

	} else {

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

func orderSingleItem(item orderItem) error {
	item.ID = uuid.NewV4().String()
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
	return err
}

func orderMultipleItems(items []orderItem) error {
	qs := make(map[string]orderItem)
	transactionID := uuid.NewV4().String()
	for _, item := range items {
		var oq string
		item.TransactionID = transactionID
		item.ID = uuid.NewV4().String()
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
	return err
}

func diamondOrder(item orderItem) error {
	r, err := dbExec(fmt.Sprintf("UPDATE diamonds SET status='ORDERED' WHERE id='%s'", item.ItemID))
	if err != nil {
		return err
	}
	if er, err := r.RowsAffected(); err != nil {
		return err
	} else if er == 0 {
		errors.Newf("Item %s not AVAILABLE any more", item.ItemID)
	}
	return nil
}

func insertIntoOrders() {

	// err := dbTransact(db, func(tx *sql.Tx) error {
	// 	q := nu.composeInsertQuery()
	// 	traceSQL(q)
	// 	if _, err := tx.Exec(q); err != nil {
	// 		return err
	// 	}
	// 	q = fmt.Sprintf(`INSERT INTO orders (user_id, level, discount, created_by) VALUES
	// 										(%s', '%d', '%d', '%s')`, a.ID, a.Level, a.Discount, a.CreatedBy)
	// 	traceSQL(q)
	// 	if _, err := tx.Exec("q"); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })
}
