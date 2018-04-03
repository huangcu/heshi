package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"
	"time"
	"util"

	"github.com/gin-gonic/gin"
)

type productStockHandleRecord struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Category       string    `json:"category"`
	Action         string    `json:"action"`
	Filename       string    `json:"file_name_uploaded"`
	FileNameOnDisk string    `json:"file_name_on_disk"`
	CreatedAt      time.Time `json:"created_at"`
}

func (p *productStockHandleRecord) newProductStockHanldeRecords() {
	q := p.composeInsertQuery()
	if _, err := dbExec(q); err != nil {
		util.Printf("Fail to add to product_stock_handle_records: %s", q)
	}
}

func getProductStockHanldeRecordsOfUser(c *gin.Context) {
	userID := c.Param("id")
	q := fmt.Sprintf(`SELECT id, category, action, filename, filenameondisk, created_at 
		FROM product_stock_handle_records 
		WHERE user_id='%s' 
		ORDER BY created_at DESC`, userID)
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	var ps []productStockHandleRecord
	for rows.Next() {
		var id, userID, category, action string
		var filename, filenameOnDisk sql.NullString
		var createdAt time.Time
		if err := rows.Scan(&id, &userID, &category, &action, &filename, &filenameOnDisk, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		pshr := productStockHandleRecord{
			ID:             id,
			UserID:         userID,
			Category:       category,
			Action:         action,
			Filename:       filename.String,
			FileNameOnDisk: filenameOnDisk.String,
			CreatedAt:      createdAt,
		}
		ps = append(ps, pshr)
	}
	c.JSON(http.StatusOK, ps)
}

func getAllProductStockHanldeRecords(c *gin.Context) {
	q := `SELECT id, user_id, category, action, filename, filenameondisk, created_at 
	FROM product_stock_handle_records 
	ORDER BY created_at DESC`
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	var ps []productStockHandleRecord
	for rows.Next() {
		var id, userID, category, action string
		var filename, filenameOnDisk sql.NullString
		var createdAt time.Time
		if err := rows.Scan(&id, &userID, &category, &action, &filename, &filenameOnDisk, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		pshr := productStockHandleRecord{
			ID:             id,
			UserID:         userID,
			Category:       category,
			Action:         action,
			Filename:       filename.String,
			FileNameOnDisk: filenameOnDisk.String,
			CreatedAt:      createdAt,
		}
		ps = append(ps, pshr)
	}
	c.JSON(http.StatusOK, ps)
}

func (p *productStockHandleRecord) composeInsertQuery() string {
	params := p.paramsKV()
	q := `INSERT INTO product_stock_handle_records (id, `
	va := fmt.Sprintf(`VALUES ('%s'`, p.ID)
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
		}
	}
	q = fmt.Sprintf("%s) %s)", q, va)
	return q
}

func (p *productStockHandleRecord) paramsKV() map[string]interface{} {
	params := make(map[string]interface{})

	if p.UserID != "" {
		params["userID"] = p.UserID
	}
	if p.Category != "" {
		params["category"] = p.Category
	}
	if p.Action != "" {
		params["action"] = p.Action
	}
	if p.Filename != "" {
		params["filename"] = p.Filename
	}
	if p.FileNameOnDisk != "" {
		params["filenameondisk"] = p.FileNameOnDisk
	}

	return params
}
