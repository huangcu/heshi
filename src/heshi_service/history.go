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
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ItemID     string    `json:"item_id"`
	TableName  string    `json:"table_name"`
	FieldName  string    `json:"field_name"`
	FieldValue string    `json:"field_value"`
	CreatedAt  time.Time `json:"created_at"`
}

func newHistoryRecords(userID, tableName, itemID string, fieldNameValue map[string]interface{}) {
	q := fmt.Sprintf(`INSERT INTO historys set user_id='%s', table_name='%s', item_id='%s',`, userID, strings.ToUpper(tableName), itemID)
	for k, v := range fieldNameValue {
		hq := fmt.Sprintf("%s id='%s', field_name='%s',", q, newV4(), k)
		switch v.(type) {
		case string:
			hq = fmt.Sprintf("%s field_value='%s'", hq, v.(string))
		case float64:
			hq = fmt.Sprintf("%s field_value='%f'", hq, v.(float64))
		case int:
			hq = fmt.Sprintf("%s field_value='%d'", hq, v.(int))
		case int64:
			hq = fmt.Sprintf("%s field_value='%d'", hq, v.(int64))
		case time.Time:
			hq = fmt.Sprintf("%s field_value='%s'", hq, v.(time.Time).Format(timeFormat))
		}
		if _, err := dbExec(hq); err != nil {
			util.Traceln("Fail to add to history: %s", hq)
		}
	}
}

func getHistory(c *gin.Context) {
	fieldName := strings.ToUpper(c.PostForm("field_name"))
	tableName := strings.ToUpper(c.PostForm("table_name"))
	itemID := c.PostForm("item_id")
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
	q := `SELECT id, user_id, item_id, table_name, field_name, field_value created_at FROM historys ORDER BY created_at DESC`
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hs []history
	for rows.Next() {
		var id, userID, itemID, tableName, fieldName, fieldValue string
		var createdAt time.Time
		if err := rows.Scan(&id, &userID, &itemID, &tableName, &fieldName, &fieldValue, &createdAt); err != nil {
			return nil, err
		}
		h := history{
			ID:         id,
			UserID:     userID,
			ItemID:     itemID,
			TableName:  tableName,
			FieldName:  fieldName,
			FieldValue: fieldValue,
			CreatedAt:  createdAt,
		}
		hs = append(hs, h)
	}
	return hs, nil
}

func getHistoryOfTable(tableName string) ([]history, error) {
	q := fmt.Sprintf(`SELECT id, user_id, item_id, field_name, field_value created_at FROM historys WHERE tableName='%s' ORDER BY created_at DESC`, tableName)
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hs []history
	for rows.Next() {
		var id, userID, itemID, fieldName, fieldValue string
		var createdAt time.Time
		if err := rows.Scan(&id, &userID, &itemID, &fieldName, &fieldValue, &createdAt); err != nil {
			return nil, err
		}
		h := history{
			ID:         id,
			UserID:     userID,
			ItemID:     itemID,
			TableName:  tableName,
			FieldName:  fieldName,
			FieldValue: fieldValue,
			CreatedAt:  createdAt,
		}
		hs = append(hs, h)
	}
	return hs, nil
}

func getHistoryOfTableField(tableName, fieldName string) ([]history, error) {
	q := fmt.Sprintf(`SELECT id, user_id, item_id, field_value created_at FROM historys WHERE table_name='%s' and field_name='%s' ORDER BY created_at DESC`, tableName, fieldName)
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hs []history
	for rows.Next() {
		var id, userID, itemID, fieldValue string
		var createdAt time.Time
		if err := rows.Scan(&id, &userID, &itemID, &fieldValue, &createdAt); err != nil {
			return nil, err
		}
		h := history{
			ID:         id,
			UserID:     userID,
			ItemID:     itemID,
			TableName:  tableName,
			FieldName:  fieldName,
			FieldValue: fieldValue,
			CreatedAt:  createdAt,
		}
		hs = append(hs, h)
	}
	return hs, nil
}

func getHistoryOfItem(itemID string) ([]history, error) {
	q := fmt.Sprintf(`SELECT id, user_id, table_name, field_name, field_value, created_at FROM historys WHERE item_id='%s' ORDER BY created_at DESC`, itemID)
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hs []history
	for rows.Next() {
		var id, userID, tableName, fieldName, fieldValue string
		var createdAt time.Time
		if err := rows.Scan(&id, &userID, &tableName, &fieldName, &fieldValue, &createdAt); err != nil {
			return nil, err
		}
		h := history{
			ID:         id,
			UserID:     userID,
			ItemID:     itemID,
			TableName:  tableName,
			FieldName:  fieldName,
			FieldValue: fieldValue,
			CreatedAt:  createdAt,
		}
		hs = append(hs, h)
	}
	return hs, nil
}
