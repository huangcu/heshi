package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"heshi/errors"
	"net/http"
	"strconv"
	"strings"
	"util"

	"github.com/gin-gonic/gin"
)

//TODO diamond accessory????? what extrainfo for??
type cartItemBase struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	ItemCategory string `json:"item_category"`
	ItemID       string `json:"item_id"`
	ItemQuantity int    `json:"item_quantity"`
	ExtraInfo    string `json:"extra_info"`
}
type shoppingCartItem struct {
	cartItemBase
	ItemPrice     float64 `json:"item_price"`
	StockQuantity int     `json:"stock_quantity"`
	Status        string  `json:"status"`
}

func getShoppingCartList(c *gin.Context) {
	scis, err := getUserShoppingCartList(c.MustGet("id").(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, scis)
}

// add to shopping cart, better only support add quantity=1;
func addToShoppingCart(c *gin.Context) {
	cib := cartItemBase{
		UserID:       c.MustGet("id").(string),
		ItemID:       c.PostForm("item_id"),
		ItemCategory: strings.ToUpper(c.PostForm("item_category")),
	}
	if c.PostForm("item_quantity") != "" {
		quantity, err := strconv.Atoi(c.PostForm("item_quantity"))
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("%s quantity is not right", c.PostForm("item_quantity")))
			return
		}
		cib.ItemQuantity = quantity
	}

	if !util.IsInArrayString(cib.ItemCategory, VALID_PRODUCTS) {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s category is not right", cib.ItemCategory))
		return
	}

	items, err := getUserShoppingCartItemBaseList(cib.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	existingItem := cib.isItemInShoppingCart(items)

	switch strings.ToUpper(cib.ItemCategory) {
	case DIAMOND:
		if existingItem != nil {
			c.JSON(http.StatusOK, "it's already in your shopping cart")
			return
		}
		cib.ItemQuantity = 1
		cib.ID = newV4()
		q := cib.composeInsertQuery()
		if _, err := dbExec(q); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	case JEWELRY:
		if err := cib.addJewelryItemToShoppingCart(existingItem); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusOK, "Not enought stock")
				return
			}
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	case GEM:
		if err := cib.addJewelryItemToShoppingCart(existingItem); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusOK, "Not enought stock")
				return
			}
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	default:
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Item Category: %s not right", c.PostForm("item_category")))
		return
	}

	c.JSON(http.StatusOK, cib)
}

func updateShoppingCart(c *gin.Context) {
	uid := c.MustGet("id").(string)
	var cibs []cartItemBase
	if err := json.Unmarshal([]byte(c.PostForm("items")), &cibs); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	items, err := getUserShoppingCartItemBaseList(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	for _, cib := range cibs {
		i, existingItem := cib.getItemInShoppingCartWithIndex(items)
		if existingItem != nil {
			// remove from existing array in, later will delete all left items as it is no long in users shopping cart
			items = append(items[:i], items[i+1:]...)
		}
		switch strings.ToUpper(cib.ItemCategory) {
		case DIAMOND:
			// for diamond, update, do nothing if already in, add when not exist
			if existingItem == nil {
				cib.ItemQuantity = 1
				cib.ID = newV4()
				q := cib.composeInsertQuery()
				if _, err := dbExec(q); err != nil {
					c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
					return
				}
			}
		case JEWELRY:
			if err := cib.addJewelryItemToShoppingCart(existingItem); err != nil {
				if err != sql.ErrNoRows {
					// for update, if not enought stock - do nothing, keep cart unchanaged
					c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
					return
				}
			}
		case GEM:
			if err := cib.addJewelryItemToShoppingCart(existingItem); err != nil {
				if err != sql.ErrNoRows {
					c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
					return
				}
			}
		}
	}

	// remove whatever left from old existing cart items, and delete them
	for _, cib := range items {
		if err := cib.removeItemFromShoppingCartByID(); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	}

	scis, err := getUserShoppingCartList(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, scis)
}

// remove from shopping cart
func removeFromShoppingCart(c *gin.Context) {
	dcib := cartItemBase{
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
		dcib.ItemQuantity = decreasedQuantity
	}
	cibs, err := getUserShoppingCartItemBaseList(dcib.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	for _, cib := range cibs {
		if cib.ID == dcib.ID && cib.UserID == dcib.UserID {
			switch cib.ItemCategory {
			case DIAMOND:
				if err := cib.removeItemFromShoppingCartByID(); err != nil {
					c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
					return
				}
			case JEWELRY, GEM:
				if cib.ItemQuantity <= dcib.ItemQuantity {
					if err := cib.removeItemFromShoppingCartByID(); err != nil {
						c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
						return
					}
				} else {
					if err := cib.decreaseItemQuantityFromShoppingCartByID(); err != nil {
						c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
						return
					}
				}
			}
		} else {
			fmt.Println(cib)
			fmt.Println(dcib)
		}
	}

	c.JSON(http.StatusOK, "SUCCESS")
}

func getUserShoppingCartList(uid string) ([]shoppingCartItem, error) {
	cibs, err := getUserShoppingCartItemBaseList(uid)
	if err != nil {
		return nil, err
	}

	var scis []shoppingCartItem
	for _, cib := range cibs {
		sci := shoppingCartItem{
			cartItemBase: cib,
		}
		switch strings.ToUpper(cib.ItemCategory) {
		case DIAMOND:
			if err := sci.getDiamondPriceQuantityStatus(); err != nil {
				return nil, err
			}
		case JEWELRY:
			if err := sci.getJewelryPriceQuantityStatus(); err != nil {
				return nil, err
			}
		case GEM:
			if err := sci.getGemPriceQuantityStatus(); err != nil {
				return nil, err
			}
			// case small_diamonds
		}
		scis = append(scis, sci)
	}
	return scis, nil
}

func getUserShoppingCartItemBaseList(uid string) ([]cartItemBase, error) {
	q := `SELECT id, item_category, item_id, item_quantity, extra_info FROM shopping_cart`
	rows, err := dbQuery(q)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var ss []cartItemBase
	for rows.Next() {
		var id, itemType, itemID string
		var itemQuantity int
		var extraInfo sql.NullString
		if err := rows.Scan(&id, &itemType, &itemID, &itemQuantity, &extraInfo); err != nil {
			return nil, err
		}
		s := cartItemBase{
			ID:           id,
			UserID:       uid,
			ItemCategory: itemType,
			ItemID:       itemID,
			ItemQuantity: itemQuantity,
			ExtraInfo:    extraInfo.String,
		}
		ss = append(ss, s)
	}

	return ss, nil
}

func (c *cartItemBase) getItemInShoppingCartWithIndex(itemList []cartItemBase) (int, *cartItemBase) {
	for k, item := range itemList {
		if c.ItemCategory == item.ItemCategory && c.ItemID == item.ItemID {
			return k, &item
		}
	}
	return 0, nil
}

func (c *cartItemBase) isItemInShoppingCart(itemList []cartItemBase) *cartItemBase {
	for _, item := range itemList {
		if c.ItemCategory == item.ItemCategory && c.ItemID == item.ItemID {
			return &item
		}
	}
	return nil
}

func (c *cartItemBase) addJewelryItemToShoppingCart(existingItem *cartItemBase) error {
	if existingItem != nil {
		c.ItemQuantity = c.ItemQuantity + existingItem.ItemQuantity
	}
	q := fmt.Sprintf(`SELECT stock_quantity 
		FROM jewelrys 
		WHERE id='%s' AND stock_quantity >=%d AND status='AVAILABLE'`, c.ItemID, c.ItemQuantity)
	var stockQuantity int
	if err := dbQueryRow(q).Scan(&stockQuantity); err != nil {
		return err
	}
	c.ID = newV4()
	q = c.composeInsertQuery()
	if existingItem != nil {
		c.ID = existingItem.ID
		q = c.composeUpdateQuery()
	}
	_, err := dbExec(q)
	if err != nil {
		return err
	}
	return nil
}

func (c *cartItemBase) addExistsGemItemToShoppingCart(existingItem *cartItemBase) error {
	if existingItem != nil {
		c.ItemQuantity = c.ItemQuantity + existingItem.ItemQuantity
	}
	q := fmt.Sprintf(`SELECT stock_quantity 
		FROM gems 
		WHERE id='%s' AND stock_quantity >=%d AND status='AVAILABLE'`, c.ItemID, c.ItemQuantity)
	var stockQuantity int
	if err := dbQueryRow(q).Scan(&stockQuantity); err != nil {
		return err
	}
	c.ID = newV4()
	q = c.composeInsertQuery()
	if existingItem != nil {
		c.ID = existingItem.ID
		q = c.composeUpdateQuery()
	}
	_, err := dbExec(q)
	if err != nil {
		return err
	}
	return nil
}

func (c *cartItemBase) removeItemFromShoppingCartByID() error {
	q := fmt.Sprintf(`DELETE FROM shopping_cart WHERE id='%s'`, c.ID)
	if _, err := dbExec(q); err != nil {
		return err
	}
	return nil
}

func (c *cartItemBase) decreaseItemQuantityFromShoppingCartByID() error {
	q := fmt.Sprintf(`UPDATE shopping_cart 
		SET item_quantity=item_quantity-%d 
		WHERE id='%s' AND item_quantity>%d`, c.ItemQuantity, c.ID, c.ItemQuantity)
	if _, err := dbExec(q); err != nil {
		return err
	}
	return nil
}

func (c *cartItemBase) removeItemFromShoppingCartByItemProperties() error {
	q := fmt.Sprintf(`DELETE FROM shopping_cart 
		WHERE user_id='%s' AND item_id='%s' AND item_category='%s'`,
		c.UserID, c.ItemID, c.ItemCategory)
	if _, err := dbExec(q); err != nil {
		return err
	}
	return nil
}

func (s *shoppingCartItem) getDiamondPriceQuantityStatus() error {
	s.StockQuantity = 1
	var priceRetail float64
	q := fmt.Sprintf(`SELECT price_retail FROM diamonds WHERE id='%s' AND status ='AVAILABLE'`, s.ItemID)
	if err := dbQueryRow(q).Scan(&priceRetail); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		s.Status = "NOT AVAILABLE"
	}
	s.ItemPrice = priceRetail
	s.Status = "AVAILABLE"
	return nil
}

func (s *shoppingCartItem) getJewelryPriceQuantityStatus() error {
	s.StockQuantity = 1
	var price float64
	var stockQuantity int
	q := fmt.Sprintf(`SELECT price,stock_quantity FROM jewelrys WHERE id='%s' AND status='AVAILABLE'`, s.ItemID)
	if err := dbQueryRow(q).Scan(&price, &stockQuantity); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		s.Status = "NOT AVAILABLE"
		return nil
	}
	s.ItemPrice = price
	s.StockQuantity = stockQuantity
	if s.StockQuantity < s.ItemQuantity {
		s.Status = "STOCK NOT ENOUGH"
		return nil

	}
	s.Status = "AVAILABLE"
	return nil
}

func (s *shoppingCartItem) getGemPriceQuantityStatus() error {
	s.StockQuantity = 1
	var price float64
	var stockQuantity int
	q := fmt.Sprintf(`SELECT price,stock_quantity FROM gems WHERE id='%s' AND status='AVAILABLE'`, s.ItemID)
	if err := dbQueryRow(q).Scan(&price, &stockQuantity); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		s.Status = "NOT AVAILABLE"
		return nil
	}
	s.ItemPrice = price
	s.StockQuantity = stockQuantity
	if s.StockQuantity < s.ItemQuantity {
		s.Status = "STOCK NOT ENOUGH"
		return nil
	}
	s.Status = "AVAILABLE"
	return nil
}
