package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func filterGems(c *gin.Context) ([]gem, error) {
	direction := "DESC"
	if c.PostForm("order") == "UP" {
		direction = "ASC"
	}
	q := fmt.Sprintf(`SELECT * FROM gems WHERE stock_quantity > 0 ORDER BY price %s`, direction)

	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	gs, err := composeGem(rows)
	if err != nil {
		return nil, err
	}
	return gs, nil
}
