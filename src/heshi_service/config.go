package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"util"

	"github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

var activeConfig config

type config struct {
	ID          string    `json:"id"`
	Discount    int       `json:"discount"`
	DiscountStr string    `json:"-"`
	Pieces      int       `json:"pieces"`
	PiecesStr   string    `json:"-"`
	Level       string    `json:"level"`
	Amount      int       `json:"amount"`
	AmountStr   string    `json:"-"`
	Type        string    `json:"type"`
	Rate        float64   `json:"rate"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

func getLevelConfig(c *gin.Context) {
	var discount, amount, pieces int
	var id, level, levelType, createdBy string
	var createdAt time.Time
	q := fmt.Sprintf("SELECT id, discount, level, amount, pieces, type, created_by,created_at FROM configs WHERE id = '%s'", c.Param("id"))
	if err := dbQueryRow(q).Scan(&id, &discount, &level, &amount, &pieces, &levelType, &createdBy, &createdAt); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	conf := config{
		ID:        id,
		Discount:  discount,
		Level:     level,
		Amount:    amount,
		Pieces:    pieces,
		Type:      levelType,
		CreatedBy: createdBy,
		CreatedAt: createdAt.Local(),
	}
	c.JSON(http.StatusOK, conf)
}

func newLevelConfig(c *gin.Context) {
	if c.PostForm("discount") == "" {
		c.JSON(http.StatusOK, "discount can not be empty")
		return
	}
	if c.PostForm("level") == "" {
		c.JSON(http.StatusOK, "level can not be empty")
		return
	}
	if c.PostForm("type") == "" {
		c.JSON(http.StatusOK, "type can not be empty")
		return
	} else if strings.ToUpper(c.PostForm("type")) == AGENT {
		if c.PostForm("pieces") == "" {
			c.JSON(http.StatusOK, "pieces can not be empty")
			return
		}
	}
	if c.PostForm("amount") == "" {
		c.JSON(http.StatusOK, "amount can not be empty")
		return
	}

	createdBy := c.MustGet("id").(string)
	id := uuid.NewV4().String()
	conf := config{
		DiscountStr: c.PostForm("discount"),
		Level:       c.PostForm("level"),
		AmountStr:   c.PostForm("amount"),
		PiecesStr:   c.PostForm("pieces"),
		Type:        c.PostForm("type"),
	}
	if vemsgs, err := conf.validateReq(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsgs) != 0 {
		c.JSON(http.StatusOK, vemsgs)
		return
	}
	q := fmt.Sprintf(`INSERT INTO configs (id, discount, level, amount, pieces, type, created_by) 
	VALUES ('%s','%d','%s', '%d','%d' ,'%s','%s')`,
		id, conf.Discount, conf.Level, conf.Amount, conf.Pieces, conf.Type, createdBy)
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, id)
}

func updateLevelConfig(c *gin.Context) {
	createdBy := c.MustGet("id").(string)
	conf := config{
		ID:          c.Param("id"),
		DiscountStr: c.PostForm("discount"),
		Level:       c.PostForm("level"),
		AmountStr:   c.PostForm("amount"),
		PiecesStr:   c.PostForm("pieces"),
		Type:        c.PostForm("type"),
		CreatedBy:   createdBy,
	}
	if vemsgs, err := conf.validateReq(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsgs) != 0 {
		c.JSON(http.StatusOK, vemsgs)
		return
	}
	q := conf.composeUpdateQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, conf.ID)
}

func getAllLevelConfigs(c *gin.Context) {
	rows, err := dbQuery("SELECT id, discount, level, amount, pieces, type,created_by,created_at FROM configs WHERE (type = '?' OR type = '?') ORDER BY created_at",
		CUSTOMER, AGENT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var cs []config
	for rows.Next() {
		var discount, amount, pieces int
		var id, level, levelType, createdBy string
		var createdAt time.Time
		rows.Scan(&id, &discount, &level, &levelType, &amount, &pieces, &createdBy, &createdAt)
		conf := config{
			ID:        id,
			Discount:  discount,
			Level:     level,
			Amount:    amount,
			Pieces:    pieces,
			Type:      levelType,
			CreatedBy: createdBy,
			CreatedAt: createdAt.Local(),
		}
		cs = append(cs, conf)
	}
	c.JSON(http.StatusOK, cs)
}

func getRateConfig(c *gin.Context) {
	var rate float64
	var id, createdBy string
	var createdAt time.Time
	q := "SELECT id, rate,created_by,created_at FROM configs WHERE type = 'RATE' ORDER BY created_at DESC LIMIT 1"
	if err := dbQueryRow(q).Scan(id, &rate, &createdBy, &createdAt); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, config{ID: id, Rate: rate, CreatedBy: createdBy, CreatedAt: createdAt.Local()})
}

func newRateConfig(c *gin.Context) {
	createdBy := c.MustGet("id").(string)
	id := uuid.NewV4().String()
	rate, err := util.StringToFloat(c.PostForm("rate"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	q := fmt.Sprintf("INSERT INTO configs (id, rate, type, created_by) VALUES ('%s','%f','%s','%s')",
		id, rate, "RATE", createdBy)
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, id)
}

func getAllRateConfigs(c *gin.Context) {
	rows, err := dbQuery("SELECT id, rate,created_by,created_at FROM configs WHERE type = 'RATE' ORDER BY created_at")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var cs []config
	for rows.Next() {
		var rate float64
		var id, createdBy string
		var createdAt time.Time
		rows.Scan(&id, &rate, &createdBy, &createdAt)
		cs = append(cs, config{ID: id, Rate: rate, CreatedBy: createdBy, CreatedAt: createdAt.Local()})
	}
	c.JSON(http.StatusOK, cs)
}

func (ac *config) getActiveRateConfig() {
	var rate float64
	var id, createdBy string
	var createdAt time.Time
	q := "SELECT id, rate,created_by,created_at FROM configs WHERE type = 'RATE' ORDER BY created_at DESC LIMIT 1"
	if err := dbQueryRow(q).Scan(&id, &rate, &createdBy, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			util.Println("fail to get active config, use default config")
		}
	}
	ac.ID = id
	ac.Rate = rate
	ac.CreatedBy = createdBy
	ac.CreatedAt = createdAt
}

func (ac *config) validateReq() ([]errors.HSMessage, error) {
	var vmsgs []errors.HSMessage
	if ac.DiscountStr != "" {
		discount, err := strconv.Atoi(ac.DiscountStr)
		if err != nil {
			return nil, err
		}
		if discount > 100 {
			vmsgs = append(vmsgs, vemsgAgentDiscountNotValid)
		}
		ac.Discount = discount
	}
	if ac.PiecesStr != "" {
		p, err := strconv.Atoi(ac.PiecesStr)
		if err != nil {
			return nil, err
		}
		ac.Pieces = p
	}
	if ac.AmountStr != "" {
		amount, err := strconv.Atoi(ac.AmountStr)
		if err != nil {
			return nil, err
		}
		ac.Amount = amount
	}

	if ac.Level != "" {
		ac.Level = "LEVEL" + strings.TrimSpace(ac.Level)
	}

	if ac.Type != "" {
		if !util.IsInArrayString(strings.ToUpper(ac.Type), []string{CUSTOMER, AGENT}) {
			vemsgNotValid.Message = "config type not valid"
			vmsgs = append(vmsgs, vemsgNotValid)
		}
		ac.Type = strings.ToUpper(ac.Type)
	}
	return vmsgs, nil
}

func (ac *config) composeUpdateQuery() string {
	params := ac.paramsKV()
	q := `UPDATE configs SET`
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

	q = fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, ac.ID)
	return q
}

func (ac *config) paramsKV() map[string]interface{} {
	params := make(map[string]interface{})

	if ac.Discount != 0 {
		params["discount"] = ac.Discount
	}
	if ac.Level != "" {
		params["level"] = ac.Level
	}
	if ac.Type != "" {
		params["type"] = ac.Type
	}
	if ac.Amount != 0 {
		params["amount"] = ac.Amount
	}
	if ac.Pieces != 0 {
		params["pieces"] = ac.Pieces
	}
	if ac.CreatedBy != "" {
		params["created_by"] = ac.CreatedBy
	}
	return params
}
