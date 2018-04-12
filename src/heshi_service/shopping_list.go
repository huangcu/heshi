package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

//TODO diamond accessory
//if the user has the right to make order	????
type shoppingItem struct {
	ID                string `json:"id"`
	UserID            string `json:"user_id"`
	ItemCategory      string `json:"item_category"`
	ItemID            string `json:"item_id"`
	ItemAccessory     int    `json:"item_accessory"`
	ConfirmedForCheck string `json:"confirmed_for_check"`
	Available         string `json:"available"`
	SpecialNotice     string `json:"special_notice"`
}

func toShoppingList(c *gin.Context) {
	action := strings.ToLower(c.Param("action"))
	s := shoppingItem{
		UserID:       c.MustGet("id").(string),
		ItemID:       c.PostForm("item_id"),
		ItemCategory: strings.ToUpper(c.PostForm("item_category")),
	}
	handleShoppingList(action, s)
}

func handleShoppingList(action string, item shoppingItem) error {
	items, err := getUserShoppingList(item.UserID)
	if err != nil {
		return err
	}
	//TODO diamond accessory
	if item.ItemCategory != "ACCESSORY" {
		switch action {
		case "add":
		case "delete":
		}
	} else {
		switch action {
		case "add":
			if !item.isItemInShoppingItemList(items) {
				if err := item.addItemToInterestedItems(); err != nil {
					return err
				}
			}
		case "delete":
			if err := item.removeItemFromInterestedItemsByItemProperties(); err != nil {
				return err
			}
		}
	}
	return nil
}

func getUserShoppingList(uid string) ([]shoppingItem, error) {
	q := `SELECT id, item_category, item_id, item_accessory, confirmed_for_check,
	available, special_notice FROM interested_items`
	rows, err := dbQuery(q)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var ss []shoppingItem
	for rows.Next() {
		var id, itemType, itemID, confirmedForCheck, available, specialNotice string
		var itemAccessory int
		if err := rows.Scan(&id, &itemType, &itemID, &itemAccessory,
			&confirmedForCheck, &available, &specialNotice); err != nil {
			return nil, err
		}
		s := shoppingItem{
			ID:                id,
			ItemCategory:      itemType,
			ItemID:            itemID,
			ItemAccessory:     itemAccessory,
			ConfirmedForCheck: confirmedForCheck,
			Available:         available,
			SpecialNotice:     specialNotice,
		}
		ss = append(ss, s)
	}

	return ss, nil
}

func (s *shoppingItem) isItemInShoppingItemList(itemList []shoppingItem) bool {
	for _, item := range itemList {
		if s.ItemCategory == item.ItemCategory && s.ItemID == item.ID {
			return true
		}
	}
	return false
}

func (s *shoppingItem) addItemToInterestedItems() error {
	q := s.composeInsertQuery()
	if _, err := dbExec(q); err != nil {
		return err
	}
	return nil
}

func (s *shoppingItem) removeItemFromInterestedItemsByID() error {
	q := `DELETE FROM interested_items WHERE id=?`
	if _, err := dbExec(q, s.ID); err != nil {
		return err
	}
	return nil
}

func (s *shoppingItem) removeItemFromInterestedItemsByItemProperties() error {
	q := `DELETE FROM interested_items WHERE user_id=? AND item_id=? AND item_category=?`
	if _, err := dbExec(q, s.UserID, s.ItemID, s.ItemCategory); err != nil {
		return err
	}
	return nil
}
