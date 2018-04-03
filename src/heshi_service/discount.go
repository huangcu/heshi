package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type discount struct {
	DiscountCode string    `form:"discount_code" json:"discount_code" binding:"required"`
	Discount     int       `form:"discount" json:"discount" binding:"required"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

func newDiscount(c *gin.Context) {
	createdBy := c.MustGet("id").(string)
	var nd discount
	if err := c.ShouldBind(&nd); err != nil {
		vemsgAgentDiscountNotValid.Message = errors.GetMessage(err)
		c.JSON(http.StatusOK, vemsgAgentDiscountNotValid)
		return
	}
	id := newV4()
	q := fmt.Sprintf(`INSERT INTO discounts (id, discount_code, discount, created_by) 
	VALUES ('%s', '%s', '%d', '%s')`, id, nd.DiscountCode, nd.Discount, createdBy)
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, id)
}

func getDiscount(c *gin.Context) {
	id := c.Param("id")
	var discountCode, createdBy string
	var discountNumber int
	var createdAt time.Time
	q := fmt.Sprintf(`SELECT discount_code, discount, created_at, createdBy FROM discounts WHERE id = '%s'`, id)
	if err := dbQueryRow(q).Scan(&discountCode, &discountNumber, &createdAt, &createdBy); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, vemsgDiscountNotExist)
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	d := discount{
		DiscountCode: discountCode,
		Discount:     discountNumber,
		CreatedBy:    createdBy,
		CreatedAt:    createdAt.Local(),
	}
	c.JSON(http.StatusOK, d)
}

func getAllDiscounts(c *gin.Context) {
	q := `SELECT discount_code, discount,created_at FROM discounts ORDER BY created_at DESC`
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var ds []discount
	for rows.Next() {
		var discountCode, createdBy string
		var discountNumber int
		var createdAt time.Time
		if err := rows.Scan(&discountCode, &discountNumber, &createdAt, &createdBy); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		d := discount{
			DiscountCode: discountCode,
			Discount:     discountNumber,
			CreatedBy:    createdBy,
			CreatedAt:    createdAt.Local(),
		}
		ds = append(ds, d)
	}
	if ds == nil {
		c.JSON(http.StatusOK, vemsgDiscountNotExist)
		return
	}
	c.JSON(http.StatusOK, ds)
}
