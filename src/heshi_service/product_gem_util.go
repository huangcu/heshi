package main

import (
	"fmt"
	"strings"
)

func (g *gem) composeInsertQuery() string {
	params := g.parmsKV()
	q := `INSERT INTO gems (id`
	va := fmt.Sprintf(`VALUES ('%s'`, g.ID)
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

func (g *gem) composeUpdateQuery() string {
	params := g.parmsKV()
	q := `UPDATE gems SET`
	for k, v := range params {
		switch v.(type) {
		case string:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(string))
		case float64:
			q = fmt.Sprintf("%s %s='%f',", q, k, v.(float64))
		case int:
			q = fmt.Sprintf("%s %s='%d',", q, k, v.(int))
		case int64:
			q = fmt.Sprintf("%s %s='%d',", q, k, v.(int64))
		}
	}
	q = fmt.Sprintf("%s WHERE id='%s'", strings.TrimSuffix(q, ","), g.ID)
	return q
}

// 	params := make(map[string]interface{})
//TODO validate input
func (g *gem) parmsKV() map[string]interface{} {
	params := make(map[string]interface{})
	if g.StockID != "" {
		params["stock_id"] = g.StockID
	}
	if g.Material != "" {
		params["material"] = g.Material
	}
	if g.Name != "" {
		params["name"] = g.Name
	}
	if g.Shape != "" {
		params["shape"] = g.Shape
	}
	if g.Size != 0 {
		params["size"] = g.Size
	}
	if g.Price != 0 {
		params["price"] = g.Price
	}
	if g.Featured != "" {
		params["featured"] = g.Featured
	}
	if g.Text != "" {
		params["text"] = g.Text
	}
	if g.Online != "" {
		params["online"] = g.Online
	}
	if g.Verified != "" {
		params["verified"] = g.Verified
	}
	if g.InStock != "" {
		params["in_stock"] = g.InStock
	}
	if g.Featured != "" {
		params["featured"] = g.Featured
	}
	if g.StockQuantity != 0 {
		params["stock_quantity"] = g.StockQuantity
	}
	if g.Profitable != "" {
		params["profitable"] = g.Profitable
	}
	if g.TotallyScanned != 0 {
		params["totally_scanned"] = g.TotallyScanned
	}
	if g.FreeAcc != "" {
		params["free_acc"] = g.FreeAcc
	}
	return params
}
