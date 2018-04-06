package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//TODO diamond accessory?????
type shoppingCartItem struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	ItemType      string  `json:"item_type"`
	ItemID        string  `json:"item_id"`
	ItemPrice     float64 `json:"item_price"`
	ItemQuantity  int     `json:"item_quantity"`
	StockQuantity int     `json:"stock_quantity"`
	Status        string  `json:"status"`
	ExtraInfo     string  `json:"extra_info"`
}

func getShoppingCartList(c *gin.Context) {
	scl, err := getUserShoppingCartList(c.MustGet("id").(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}

	for _, sc := range scl {
		switch sc.ItemType {
		case DIAMOND:
			if err := sc.getDiamondPriceQuantityStatus(); err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
		case JEWELRY:
			if err := sc.getJewelryPriceQuantityStatus(); err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
		case GEM:
			if err := sc.getGemPriceQuantityStatus(); err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
		default:
			// small_diamonds
		}
	}
	c.JSON(http.StatusOK, scl)
}

// add to or remove from shopping cart
func addToShoppingCart(c *gin.Context) {
	quantity, err := strconv.Atoi(c.PostForm("item_quantity"))
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s quantity is not right", c.PostForm("item_quantity")))
		return
	}

	sc := shoppingCartItem{
		UserID:       c.MustGet("id").(string),
		ItemID:       c.PostForm("item_id"),
		ItemType:     c.PostForm("item_type"),
		ItemQuantity: quantity,
	}

	items, err := getUserShoppingCartList(sc.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	existingItem := sc.isItemInShoppingCart(items)
	if existingItem != nil {
		if err := sc.addItemToShoppingCart(); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	}
	switch sc.ItemType {
	case DIAMOND:
		c.JSON(http.StatusOK, "it's already in your shopping cart")
		return
	case JEWELRY:
		if err := sc.addExistsJewelryItemToShoppingCart(existingItem); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusOK, "Not enought stock")
				return
			}
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	case GEM:
		if err := sc.addExistsJewelryItemToShoppingCart(existingItem); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	}
	c.JSON(http.StatusOK, sc)
}

func removeFromShoppingCart(c *gin.Context) {
	s := shoppingCartItem{
		UserID: c.MustGet("id").(string),
		ID:     c.Param("id"),
	}
	iq := c.PostForm("item_quantity")
	if iq != "" {
		decreasedQuantity, err := strconv.Atoi(iq)
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("%s quantity is not right", iq))
			return
		}
		s.ItemQuantity = decreasedQuantity
		if err := s.decreaseItemQuantityFromShoppingCartByID(); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	} else {
		if err := s.removeItemFromShoppingCartByID(); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	}

	// after remove, return existing shopping cart item
	items, err := getUserShoppingCartList(s.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, items)
}

func getUserShoppingCartList(uid string) ([]shoppingCartItem, error) {
	q := `SELECT id, item_type, item_id, item_quantity, extra_info FROM shopping_cart`
	rows, err := dbQuery(q)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var ss []shoppingCartItem
	for rows.Next() {
		var id, itemType, itemID, extraInfo string
		var itemQuantity int
		if err := rows.Scan(&id, &itemType, &itemID, &itemQuantity, &extraInfo); err != nil {
			return nil, err
		}
		s := shoppingCartItem{
			ID:           id,
			ItemType:     itemType,
			ItemID:       itemID,
			ItemQuantity: itemQuantity,
			ExtraInfo:    extraInfo,
		}
		ss = append(ss, s)
	}

	return ss, nil
}

func (s *shoppingCartItem) isItemInShoppingCart(itemList []shoppingCartItem) *shoppingCartItem {
	for _, item := range itemList {
		if s.ItemType == item.ItemType && s.ItemID == item.ID {
			return &item
		}
	}
	return nil
}

func (s *shoppingCartItem) addItemToShoppingCart() error {
	q := s.composeInsertQuery()
	if _, err := dbExec(q); err != nil {
		return err
	}
	return nil
}

func (s *shoppingCartItem) addExistsJewelryItemToShoppingCart(existingItem *shoppingCartItem) error {
	newItemQunatity := s.ItemQuantity + existingItem.ItemQuantity
	q := fmt.Sprintf(`SELECT stock_quantity 
		FROM jewelrys 
		WHERE id='%s' AND stock_quantity >'%d`, s.ItemID, newItemQunatity)
	var stockQuantity int
	if err := dbQueryRow(q).Scan(&stockQuantity); err != nil {
		return err
	}
	_, err := dbExec("UPDATE shopping_cart SET item_quantity=? WHERE id=?", newItemQunatity, existingItem.ID)
	if err != nil {
		return err
	}
	s.ID = existingItem.ID
	s.ItemQuantity = newItemQunatity
	return nil
}

func (s *shoppingCartItem) addExistsGemItemToShoppingCart(existingItem *shoppingCartItem) error {
	newItemQunatity := s.ItemQuantity + existingItem.ItemQuantity
	q := fmt.Sprintf(`SELECT stock_quantity 
		FROM gems 
		WHERE id='%s' AND stock_quantity >'%d`, s.ItemID, newItemQunatity)
	var stockQuantity int
	if err := dbQueryRow(q).Scan(&stockQuantity); err != nil {
		return err
	}
	_, err := dbExec("UPDATE shopping_cart SET item_quantity=? WHERE id=?", newItemQunatity, existingItem.ID)
	if err != nil {
		return err
	}
	s.ID = existingItem.ID
	s.ItemQuantity = newItemQunatity
	return nil
}

func (s *shoppingCartItem) removeItemFromShoppingCartByID() error {
	q := `DELETE FROM shopping_cart WHERE id=?`
	if _, err := dbExec(q, s.ID); err != nil {
		return err
	}
	return nil
}

func (s *shoppingCartItem) decreaseItemQuantityFromShoppingCartByID() error {
	q := fmt.Sprintf(`UPDATE shopping_cart 
		SET item_quantity=item_quantity-%d 
		WHERE id='%s' AND item_quantity>%d`, s.ItemQuantity, s.ID, s.ItemQuantity)
	if _, err := dbExec(q); err != nil {
		return err
	}
	return nil
}

func (s *shoppingCartItem) removeItemFromShoppingCartByItemProperties() error {
	q := `DELETE FROM shopping_cart WHERE user_id=? AND item_id=? AND item_type=?`
	if _, err := dbExec(q, s.UserID, s.ItemID, s.ItemType); err != nil {
		return err
	}
	return nil
}

func (s *shoppingCartItem) getDiamondPriceQuantityStatus() error {
	s.StockQuantity = 1
	var priceRetail float64
	if err := dbQueryRow(`SELECT price_retail FROM diamonds WHERE id='?'`, s.ItemID).Scan(&priceRetail); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		s.Status = "NOT AVAIABLE"
	}
	s.ItemPrice = priceRetail
	s.Status = "AVAIABLE"
	return nil
}

func (s *shoppingCartItem) getJewelryPriceQuantityStatus() error {
	s.StockQuantity = 1
	var price float64
	var stockQuantity int
	if err := dbQueryRow(`SELECT price,stock_quantity FROM jewelrys WHERE id='?'`, s.ItemID).Scan(&price, &stockQuantity); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		s.Status = "NOT AVAIABLE"
	}
	s.ItemPrice = price
	s.StockQuantity = stockQuantity
	if s.StockQuantity < s.ItemQuantity {
		s.Status = "STOCK NOT ENOUGH"
	}
	s.Status = "AVAIABLE"
	return nil
}

func (s *shoppingCartItem) getGemPriceQuantityStatus() error {
	s.StockQuantity = 1
	var price float64
	var stockQuantity int
	if err := dbQueryRow(`SELECT price,stock_quantity FROM gems WHERE id='?'`, s.ItemID).Scan(&price, &stockQuantity); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		s.Status = "NOT AVAIABLE"
	}
	s.ItemPrice = price
	s.StockQuantity = stockQuantity
	if s.StockQuantity < s.ItemQuantity {
		s.Status = "STOCK NOT ENOUGH"
	}
	s.Status = "AVAIABLE"
	return nil
}
