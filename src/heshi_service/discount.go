package main

import (
	"fmt"
	"heshi/errors"
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

type discount struct {
	DiscountCode string `json:"discount_code"`
	Discount     int    `json:"discount"`
	CreatedAt    int64  `json:"created_at"`
}

func newDiscount(c *gin.Context) {
	discountCode := c.PostForm("discount_code")
	discountNumber := c.PostForm("discount")
	id := uuid.NewV4().String()
	q := fmt.Sprintf(`INSERT INTO discounts (id, discount_code, discount) VALUES(%s, %s，%s)`, id, discountCode, discountNumber)
	fmt.Println(q)
	if _, err := db.Exec(q); err != nil {
		c.String(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.String(http.StatusOK, id)
}

func getDiscount(c *gin.Context) {
	id := c.Query("id")
	var discountCode string
	var discountNumber int
	var createdAt time.Time
	q := fmt.Sprintf(`SELECT discount_code, discount,created_at FROM discounts WHERE id = %s`, id)
	if err := db.QueryRow(q).Scan(&discountCode, &discountNumber, &createdAt); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	d := discount{
		DiscountCode: discountCode,
		Discount:     discountNumber,
		CreatedAt:    createdAt.Unix(),
	}
	c.JSON(http.StatusOK, d)
}

func getDiscounts(c *gin.Context) {
	q := `SELECT discount_code, discount,created_at FROM discounts ORDER BY created_at DESC`
	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var ds []discount
	if rows.Next() {
		var discountCode string
		var discountNumber int
		var createdAt time.Time
		if err := rows.Scan(&discountCode, &discountCode, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		d := discount{
			DiscountCode: discountCode,
			Discount:     discountNumber,
			CreatedAt:    createdAt.Unix(),
		}
		ds = append(ds, d)
	}
	c.JSON(http.StatusOK, ds)
}
