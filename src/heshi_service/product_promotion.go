package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"heshi/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"util"

	"github.com/gin-gonic/gin"
)

type promotion struct {
	ID        string     `json:"id,omitempty"`
	PromType  string     `json:"prom_type,omitempty"`
	BeginAt   *time.Time `json:"begin_at,omitempty"`
	EndAt     *time.Time `json:"end_at,omitempty"`
	CreatedBy string     `json:"created_by,omitempty"`
	Status    string     `json:"status,omitempty"`
	PromotionDiscount
	PromotionPrice
}

type PromotionDiscount struct {
	PromDiscount int `json:"prom_discount,omitempty"`
}
type PromotionPrice struct {
	PromPrice float64 `json:"prom_price,omitempty"`
}

var PROM_TYPE = []string{"DISCOUNT", "FREE_ACCESSORY", "SPECIAL_OFFER"}

// promote|depromote - promotion_id = ""
type promotionProduct struct {
	ItemID       string `json:"item_id"`
	ItemCategory string `json:"item_category"`
	PromotionID  string `json:"promotion_id"`
}

func promoteProducts(c *gin.Context) {
	updatedBy := c.MustGet("id").(string)
	var promProducts []promotionProduct
	if err := json.Unmarshal([]byte(c.PostForm("proms")), &promProducts); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}
	for _, promProduct := range promProducts {
		promtionmap := make(map[string]interface{})
		promtionmap["promotion_id"] = promProduct.PromotionID
		switch strings.ToUpper(promProduct.ItemCategory) {
		case DIAMOND:
			q := fmt.Sprintf(`UPDATE diamonds SET promotion_id='%s' WHERE id='%s'`, promProduct.PromotionID, promProduct.ItemID)
			r, err := dbExec(q)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			rc, err := r.RowsAffected()
			if err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			if int(rc) == 1 {
				go newHistoryRecords(updatedBy, "diamonds", promProduct.ItemID, promtionmap)
			}
		case JEWELRY:
			q := fmt.Sprintf(`UPDATE jewelrys SET promotion_id='%s' WHERE id='%s'`, promProduct.PromotionID, promProduct.ItemID)
			r, err := dbExec(q)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			rc, err := r.RowsAffected()
			if err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			if int(rc) == 1 {
				go newHistoryRecords(updatedBy, "jewelrys", promProduct.ItemID, promtionmap)
			}
		case GEM:
			q := fmt.Sprintf(`UPDATE gems SET promotion_id='%s' WHERE id='%s'`, promProduct.PromotionID, promProduct.ItemID)
			r, err := dbExec(q)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			rc, err := r.RowsAffected()
			if err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			if int(rc) == 1 {
				go newHistoryRecords(updatedBy, "gems", promProduct.ItemID, promtionmap)
			}
		}
	}
	c.JSON(http.StatusOK, "SUCCESS")
}

func newPromotion(c *gin.Context) {
	createdBy := c.MustGet("id").(string)
	promType := strings.ToUpper(c.PostForm("prom_type"))
	if !util.IsInArrayString(promType, PROM_TYPE) {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s is not valid", promType))
		return
	}
	beginAt, err := time.Parse(timeFormat, c.PostForm("begin_at"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}
	endAt, err := time.Parse(timeFormat, c.PostForm("end_at"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}
	if !endAt.After(time.Now()) {
		c.JSON(http.StatusBadRequest, "already pass end time")
		return
	}
	if !endAt.After(beginAt) {
		c.JSON(http.StatusBadRequest, "end time must later than begin time")
		return
	}
	p := promotion{
		ID:        newV4(),
		PromType:  promType,
		BeginAt:   &beginAt,
		EndAt:     &endAt,
		CreatedBy: createdBy,
	}
	status := strings.ToUpper(c.PostForm("status"))
	if status != "" {
		if !util.IsInArrayString(status, []string{"ACTIVE", "INACTIVE"}) {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("%s is not valid", c.PostForm("status")))
			return
		}
		p.Status = status
	} else {
		p.Status = "INACTIVE"
	}

	switch promType {
	case "DISCOUNT":
		dis, err := strconv.Atoi(c.PostForm("prom_discount"))
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.GetMessage(err))
			return
		}
		p.PromDiscount = dis
	case "FREE_ACCESSORY":
	case "SPECIAL_OFFER":
		pp, err := util.StringToFloat(c.PostForm("prom_price"))
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.GetMessage(err))
			return
		}
		p.PromPrice = pp
	}

	q := p.composeInsertQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, p)
}

func updatePromotion(c *gin.Context) {
	pid := c.Param("id")
	p, err := getPromotionByID(pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}

	updatedBy := c.MustGet("id").(string)
	promType := strings.ToUpper(c.PostForm("prom_type"))
	if promType != "" {
		if !util.IsInArrayString(promType, PROM_TYPE) {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("%s is not valid", promType))
			return
		}
		p.PromType = promType
	}
	switch p.PromType {
	case "DISCOUNT":
		if c.PostForm("prom_discount") != "" {
			dis, err := strconv.Atoi(c.PostForm("prom_discount"))
			if err != nil {
				c.JSON(http.StatusBadRequest, errors.GetMessage(err))
				return
			}
			p.PromDiscount = dis
		}
	case "FREE_ACCESSORY":
	case "SPECIAL_OFFER":
		if c.PostForm("prom_price") != "" {
			pp, err := util.StringToFloat(c.PostForm("prom_price"))
			if err != nil {
				c.JSON(http.StatusBadRequest, errors.GetMessage(err))
				return
			}
			p.PromPrice = pp
		}
	}
	if promType == "" {
		p.PromType = ""
	}
	if c.PostForm("begin_at") != "" {
		beginAt, err := time.Parse(timeFormat, c.PostForm("begin_at"))
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.GetMessage(err))
			return
		}
		p.BeginAt = &beginAt
	}
	if c.PostForm("end_at") != "" {
		endAt, err := time.Parse(timeFormat, c.PostForm("end_at"))
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.GetMessage(err))
			return
		}
		p.EndAt = &endAt
	}
	if !p.EndAt.After(time.Now()) {
		c.JSON(http.StatusBadRequest, "already pass end time")
		return
	}
	if !p.EndAt.After(*p.BeginAt) {
		c.JSON(http.StatusBadRequest, "end time must later than begin time")
		return
	}

	status := strings.ToUpper(c.PostForm("status"))
	if status != "" {
		if !util.IsInArrayString(status, []string{"ACTIVE", "INACTIVE"}) {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("%s is not valid", c.PostForm("status")))
			return
		}
		p.Status = status
	} else {
		p.Status = ""
	}
	q := p.composeUpdateQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, p)
	go newHistoryRecords(updatedBy, "promotions", pid, p.paramsKV())
}

func getPromotion(c *gin.Context) {
	id := c.Param("id")
	p, err := getPromotionByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadGateway, "fail to find promotion")
			return
		}
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, p)
}

func getAllPromotions(c *gin.Context) {
	q := `SELECT id, prom_type, prom_discount, prom_price, status, begin_at, end_at, created_by
		 FROM promotions ORDER BY created_at DESC`
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	var ps []promotion
	for rows.Next() {
		var id, promType, createdBy, status string
		var promDiscount sql.NullInt64
		var promPrice sql.NullFloat64
		var beginAt, endAt time.Time

		if err := rows.Scan(&id, &promType, &promDiscount, &promPrice, &status, &beginAt, &endAt, &createdBy); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		b := beginAt
		e := endAt
		p := promotion{
			ID:        id,
			PromType:  promType,
			BeginAt:   &b,
			EndAt:     &e,
			Status:    status,
			CreatedBy: createdBy,
		}
		switch promType {
		case "DISCOUNT":
			p.PromDiscount = int(promDiscount.Int64)
		case "FREE_ACCESSORY":
		case "SPECIAL_OFFER":
			p.PromPrice = promPrice.Float64
		}
		ps = append(ps, p)
	}
	c.JSON(http.StatusOK, ps)
}

func getPromotionByID(pid string) (*promotion, error) {
	var promType, createdBy, status string
	var promDiscount sql.NullInt64
	var promPrice sql.NullFloat64
	var beginAt, endAt time.Time
	q := fmt.Sprintf(`SELECT prom_type, prom_discount, prom_price, status, begin_at, end_at, created_by
		 FROM promotions WHERE id = '%s'`, pid)
	if err := dbQueryRow(q).Scan(&promType, &promDiscount, &promPrice, &status, &beginAt, &endAt, &createdBy); err != nil {
		return nil, err
	}
	b := beginAt
	e := endAt
	p := promotion{
		ID:        pid,
		PromType:  promType,
		BeginAt:   &b,
		EndAt:     &e,
		Status:    status,
		CreatedBy: createdBy,
	}
	switch promType {
	case "DISCOUNT":
		p.PromDiscount = int(promDiscount.Int64)
	case "FREE_ACCESSORY":
	case "SPECIAL_OFFER":
		p.PromPrice = promPrice.Float64
	}
	return &p, nil
}

func (p *promotion) composeInsertQuery() string {
	params := p.paramsKV()
	q := `INSERT INTO promotions (id`
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
		case time.Time:
			va = fmt.Sprintf("%s, '%s'", va, v.(time.Time).Format(timeFormat))
		}
	}
	q = fmt.Sprintf("%s) %s)", q, va)
	return q
}

func (p *promotion) composeUpdateQuery() string {
	params := p.paramsKV()
	q := `UPDATE promotions SET`
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
		case time.Time:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(time.Time).Format(timeFormat))
		}
	}

	q = fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, p.ID)
	return q
}

func (p *promotion) paramsKV() map[string]interface{} {
	params := make(map[string]interface{})

	if p.PromType != "" {
		params["prom_type"] = p.PromType
	}
	if p.PromDiscount != 0 {
		params["prom_discount"] = p.PromDiscount
	}
	if p.PromPrice != 0 {
		params["prom_price"] = p.PromPrice
	}
	if !p.BeginAt.IsZero() {
		params["begin_at"] = *p.BeginAt
	}
	if !p.EndAt.IsZero() {
		params["end_at"] = *p.EndAt
	}
	if p.CreatedBy != "" {
		params["created_by"] = p.CreatedBy
	}
	if p.Status != "" {
		params["status"] = p.Status
	}
	return params
}
