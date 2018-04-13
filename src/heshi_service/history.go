package main

import (
	"fmt"
	"heshi/errors"
	"net/http"
	"strings"
	"time"
	"util"

	"github.com/gin-gonic/gin"
)

type history struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ItemID    string    `json:"item_id"`
	TableName string    `json:"table_name"`
	FieldName string    `json:"field_name"`
	NewValue  string    `json:"new_value"`
	OldValue  string    `json:"old_value"`
	CreatedAt time.Time `json:"created_at"`
}

func newHistoryRecords(userID, tableName, itemID string, fieldNameValue map[string]interface{}) {
	for fieldName, v := range fieldNameValue {
		hq := fmt.Sprintf(`INSERT INTO historys (user_id,table_name,item_id,field_name, old_value, new_value)
		VALUES ('%s','%s','%s','%s', (SELECT %s FROM %s WHERE id='%s')`,
			userID, strings.ToUpper(tableName), itemID, fieldName, fieldName, tableName, itemID)
		switch v.(type) {
		case string:
			hq = fmt.Sprintf("%s, '%s')", hq, v.(string))
		case float64:
			hq = fmt.Sprintf("%s, '%f')", hq, v.(float64))
		case int:
			hq = fmt.Sprintf("%s, '%d')", hq, v.(int))
		case int64:
			hq = fmt.Sprintf("%s, '%d')", hq, v.(int64))
		case time.Time:
			hq = fmt.Sprintf("%s, '%s')", hq, v.(time.Time).Format(timeFormat))
		}
		if _, err := dbExec(hq); err != nil {
			util.Tracef("Fail to add to history: %s", hq)
		}
	}
}

func deleteHistoryRecords(userID, tableName, itemID string, fieldNameValue map[string]interface{}) {
	for fieldName, v := range fieldNameValue {
		hq := fmt.Sprintf(`DELETE FROM historys 
			WHERE userID='%s' AND table_name='%s' AND item_id='%s' AND field_name='%s' AND new_value=`,
			userID, tableName, itemID, fieldName)
		switch v.(type) {
		case string:
			hq = fmt.Sprintf("%s'%s' ORDER BY created_at DESC LIMIT 1", hq, v.(string))
		case float64:
			hq = fmt.Sprintf("%s'%f' ORDER BY created_at DESC LIMIT 1", hq, v.(float64))
		case int:
			hq = fmt.Sprintf("%s'%d' ORDER BY created_at DESC LIMIT 1", hq, v.(int))
		case int64:
			hq = fmt.Sprintf("%s'%d' ORDER BY created_at DESC LIMIT 1", hq, v.(int64))
		case time.Time:
			hq = fmt.Sprintf("%s'%s' ORDER BY created_at DESC LIMIT 1", hq, v.(time.Time).Format(timeFormat))
		}
		if _, err := dbExec(hq); err != nil {
			util.Tracef("Fail to add to history: %s", hq)
		}
	}
}

func getHistory(c *gin.Context) {
	fieldName := strings.ToUpper(c.Query("field_name"))
	tableName := strings.ToUpper(c.Query("table_name"))
	itemID := c.Query("item_id")
	if itemID != "" {
		hs, err := getHistoryOfItem(itemID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, hs)
		return
	}
	if tableName != "" {
		if fieldName != "" {
			hs, err := getHistoryOfTableField(tableName, fieldName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			c.JSON(http.StatusOK, hs)
			return
		}
		hs, err := getHistoryOfTable(tableName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, hs)
		return
	}
	hs, err := getAllHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, hs)
}

func getAllHistory() ([]history, error) {
	q := `SELECT id, user_id, item_id, table_name, field_name, new_value, old_value, created_at 
	FROM historys ORDER BY created_at DESC`
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hs []history
	for rows.Next() {
		var id, userID, itemID, tableName, fieldName, newValue, oldValue string
		var createdAt time.Time
		if err := rows.Scan(&id, &userID, &itemID, &tableName, &fieldName, &newValue, &oldValue, &createdAt); err != nil {
			return nil, err
		}
		h := history{
			ID:        id,
			UserID:    userID,
			ItemID:    itemID,
			TableName: tableName,
			FieldName: fieldName,
			NewValue:  newValue,
			OldValue:  oldValue,
			CreatedAt: createdAt,
		}
		hs = append(hs, h)
	}
	return hs, nil
}

func getHistoryOfTable(tableName string) ([]history, error) {
	q := fmt.Sprintf(`SELECT id, user_id, item_id, field_name, new_value, old_value, created_at 
		FROM historys WHERE table_name='%s' ORDER BY created_at DESC`, tableName)
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hs []history
	for rows.Next() {
		var id, userID, itemID, fieldName, newValue, oldValue string
		var createdAt time.Time
		if err := rows.Scan(&id, &userID, &itemID, &fieldName, &newValue, &oldValue, &createdAt); err != nil {
			return nil, err
		}
		h := history{
			ID:        id,
			UserID:    userID,
			ItemID:    itemID,
			TableName: tableName,
			FieldName: fieldName,
			NewValue:  newValue,
			OldValue:  oldValue,
			CreatedAt: createdAt,
		}
		hs = append(hs, h)
	}
	return hs, nil
}

func getHistoryOfTableField(tableName, fieldName string) ([]history, error) {
	q := fmt.Sprintf(`SELECT id, user_id, item_id, new_value, old_value, created_at 
		FROM historys WHERE table_name='%s' and field_name='%s' ORDER BY created_at DESC`, tableName, fieldName)
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hs []history
	for rows.Next() {
		var id, userID, itemID, newValue, oldValue string
		var createdAt time.Time
		if err := rows.Scan(&id, &userID, &itemID, &newValue, &oldValue, &createdAt); err != nil {
			return nil, err
		}
		h := history{
			ID:        id,
			UserID:    userID,
			ItemID:    itemID,
			TableName: tableName,
			FieldName: fieldName,
			NewValue:  newValue,
			OldValue:  oldValue,
			CreatedAt: createdAt,
		}
		hs = append(hs, h)
	}
	return hs, nil
}

func getHistoryOfItem(itemID string) ([]history, error) {
	q := fmt.Sprintf(`SELECT id, user_id, table_name, field_name, new_value, old_value, created_at 
		FROM historys WHERE item_id='%s' ORDER BY created_at DESC`, itemID)
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hs []history
	for rows.Next() {
		var id, userID, tableName, fieldName, newValue, oldValue string
		var createdAt time.Time
		if err := rows.Scan(&id, &userID, &tableName, &fieldName, &newValue, &oldValue, &createdAt); err != nil {
			return nil, err
		}
		h := history{
			ID:        id,
			UserID:    userID,
			ItemID:    itemID,
			TableName: tableName,
			FieldName: fieldName,
			NewValue:  newValue,
			OldValue:  oldValue,
			CreatedAt: createdAt,
		}
		hs = append(hs, h)
	}
	return hs, nil
}
