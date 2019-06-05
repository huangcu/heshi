package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

//TODO search KEY: stock_id and certificate number???
func searchGems(c *gin.Context) ([]gem, error) {
	q := fmt.Sprintf(`SELECT gems.id, stock_id, shape, material, size, name, text, images, certificate, 
	gems.status, verified, featured, price, stock_quantity, profitable, 
	 totally_scanned, free_acc, last_scan_at,offline_at, 
	promotions.id, prom_type, prom_discount, prom_price, begin_at, end_at, promotions.status 
	FROM gems 
	LEFT JOIN promotions ON gems.promotion_id=promotions.id 
	WHERE gems.status IN ('AVAILABLE','OFFLINE') AND stock_quantity > 0 AND stock_id='%s' OR certificate='%s'`,
		strings.ToUpper(c.PostForm("ref")), strings.ToUpper(c.PostForm("ref")))

	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gs, err := composeGem(rows)
	if err != nil {
		return nil, err
	}
	return gs, nil
}

func filterGems(c *gin.Context) ([]gem, error) {
	direction := "DESC"
	if c.PostForm("order") == "UP" {
		direction = "ASC"
	}
	q := fmt.Sprintf(`SELECT gems.id, stock_id, shape, material, size, name, text, images, certificate, 
	gems.status, verified, featured, price, stock_quantity, profitable,
	 totally_scanned, free_acc, last_scan_at,offline_at, 
	promotions.id, prom_type, prom_discount, prom_price, begin_at, end_at, promotions.status 
	FROM gems 
	LEFT JOIN promotions ON gems.promotion_id=promotions.id 
	WHERE gems.status IN ('AVAILABLE','OFFLINE') AND stock_quantity > 0 ORDER BY price %s`, direction)

	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gs, err := composeGem(rows)
	if err != nil {
		return nil, err
	}
	return gs, nil
}
