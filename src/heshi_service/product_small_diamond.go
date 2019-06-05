package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type smallDiamond struct {
	ID       string  `json:"id"`
	SizeFrom float64 `json:"size_from"`
	SizeTo   float64 `json:"size_to"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func getAllSmallDiamonds(c *gin.Context) {
	q := selectSmallDiamondQuery("")
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ds, err := composeSmallDiamond(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ds)
}

func getSmallDiamond(c *gin.Context) {
	q := selectSmallDiamondQuery(c.Param("id"))
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ds, err := composeSmallDiamond(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if ds == nil {
		vemsgNotExist.Message = fmt.Sprintf("Fail to find small diamond with id: %s", c.Param("id"))
		c.JSON(http.StatusOK, vemsgNotExist)
		return
	}
	c.JSON(http.StatusOK, ds)
}

//TODO
func newSmallDiamond(c *gin.Context) {
	q := selectSmallDiamondQuery(c.Param("id"))
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ds, err := composeSmallDiamond(rows)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, fmt.Sprintf("Fail to find product with id: %s", c.Param("id")))
			return
		}
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ds)
}

func composeSmallDiamond(rows *sql.Rows) ([]smallDiamond, error) {
	var id string
	var sizeFrom, sizeTo, price float64
	var quantity int

	var ds []smallDiamond
	for rows.Next() {
		if err := rows.Scan(&id, &sizeFrom, &sizeTo, &price, &quantity); err != nil {
			return nil, err
		}
		d := smallDiamond{
			ID:       id,
			SizeFrom: sizeFrom,
			SizeTo:   sizeTo,
			Price:    price,
			Quantity: quantity}
		ds = append(ds, d)
	}
	return ds, nil
}

func selectSmallDiamondQuery(id string) string {
	q := `SELECT id, size_from, size_to, price, stock_quantity FROM small_diamonds`

	if id != "" {
		q = fmt.Sprintf("%s WHERE id='%s'", q, id)
	}
	return q
}
