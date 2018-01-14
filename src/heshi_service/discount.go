package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

type discount struct {
	DiscountCode string    `form:"discount_code" json:"discount_code" binding:"required"`
	Discount     int       `form:"discount" json:"discount" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
}

func newDiscount(c *gin.Context) {
	var nd discount
	if err := c.ShouldBind(&nd); err != nil {
		c.JSON(http.StatusOK, VEMSG_AGENT_DISCOUNT_NOT_VALID)
		return
	}
	id := uuid.NewV4().String()
	q := fmt.Sprintf(`INSERT INTO discounts (id, discount_code, discount) VALUES ('%s', '%s', '%d')`, id, nd.DiscountCode, nd.Discount)
	if _, err := db.Exec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, id)
}

func getDiscount(c *gin.Context) {
	id := c.Param("id")
	var discountCode string
	var discountNumber int
	var createdAt time.Time
	q := fmt.Sprintf(`SELECT discount_code, discount,created_at FROM discounts WHERE id = '%s'`, id)
	if err := db.QueryRow(q).Scan(&discountCode, &discountNumber, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, VEMSG_DISCOUNT_NOT_EXIST)
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	d := discount{
		DiscountCode: discountCode,
		Discount:     discountNumber,
		CreatedAt:    createdAt.Local(),
	}
	c.JSON(http.StatusOK, d)
}

func getAllDiscounts(c *gin.Context) {
	q := `SELECT discount_code, discount,created_at FROM discounts ORDER BY created_at DESC`
	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var ds []discount
	for rows.Next() {
		var discountCode string
		var discountNumber int
		var createdAt time.Time
		if err := rows.Scan(&discountCode, &discountNumber, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		d := discount{
			DiscountCode: discountCode,
			Discount:     discountNumber,
			CreatedAt:    createdAt.Local(),
		}
		ds = append(ds, d)
	}
	if ds == nil {
		c.JSON(http.StatusOK, VEMSG_DISCOUNT_NOT_EXIST)
		return
	}
	c.JSON(http.StatusOK, ds)
}
