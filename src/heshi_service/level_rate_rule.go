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

	"github.com/gin-gonic/gin"
)

// customer, agent 升级标准 以及 level和 discount的map关系
// Rule type: CUSTOMER. AGENT, RATE
// CUSTOMER: LEVEL, DISCOUNT, AMOUNT
// AGENT: LEVEL, DISCOUNT, AMOUNT, PIECES, return_point_percent
// TODO limit LEVEL values to predefined
type levelRule struct {
	ID                    string    `json:"id"`
	Level                 int       `json:"level"`
	LevelStr              string    `json:"-"`
	Discount              int       `json:"discount"`
	DiscountStr           string    `json:"-"`
	Pieces                int       `json:"pieces"`
	PiecesStr             string    `json:"-"`
	Amount                int       `json:"amount"`
	AmountStr             string    `json:"-"`
	ReturnPointPercent    int       `json:"return_point_percent"`
	ReturnPointPercentStr string    `json:"-"`
	RuleType              string    `json:"rule_type"`
	CreatedBy             string    `json:"created_by"`
	CreatedAt             time.Time `json:"created_at"`
}

// RATE: EXCHANGE RATE FLOAT
type exchangeRateFloat struct {
	ID                string    `json:"id"`
	ExchangeRateFloat float64   `json:"rate"`
	RuleType          string    `json:"rule_type"`
	CreatedBy         string    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
}

func getLevelConfig(c *gin.Context) {
	var discount, amount, pieces, returnPointPercent sql.NullInt64
	var level sql.NullInt64
	var id, ruleType, createdBy string
	var createdAt time.Time
	q := fmt.Sprintf(`SELECT id, discount, level, amount, pieces, return_point_percent, 
		rule_type, created_by,created_at 
		FROM level_rate_rules 
		WHERE id = '%s' AND rule_type!='RATE'`, c.Param("id"))
	if err := dbQueryRow(q).Scan(&id, &discount, &level, &amount, &pieces, &returnPointPercent,
		&ruleType, &createdBy, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, "rule not exist")
			return
		}
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	conf := levelRule{
		ID:                 id,
		Discount:           int(discount.Int64),
		Level:              int(level.Int64),
		Amount:             int(amount.Int64),
		Pieces:             int(pieces.Int64),
		RuleType:           ruleType,
		ReturnPointPercent: int(returnPointPercent.Int64),
		CreatedBy:          createdBy,
		CreatedAt:          createdAt,
	}
	c.JSON(http.StatusOK, conf)
}

func newLevelConfig(c *gin.Context) {
	if c.PostForm("discount") == "" {
		c.JSON(http.StatusBadRequest, "discount can not be empty")
		return
	}
	if c.PostForm("level") == "" {
		c.JSON(http.StatusBadRequest, "level can not be empty")
		return
	}
	ruleType := strings.ToUpper(c.PostForm("rule_type"))
	if ruleType == "" {
		c.JSON(http.StatusBadRequest, "rule type can not be empty")
		return
	}
	conf := levelRule{
		DiscountStr: c.PostForm("discount"),
		LevelStr:    c.PostForm("level"),
		RuleType:    c.PostForm("rule_type"),
	}
	switch ruleType {
	case AGENT:
		if c.PostForm("pieces") == "" && c.PostForm("amount") == "" {
			c.JSON(http.StatusBadRequest, "must specifiy pieces or amount")
			return
		}
		conf.AmountStr = c.PostForm("amount")
		conf.PiecesStr = c.PostForm("pieces")
		conf.ReturnPointPercentStr = c.PostForm("return_point_percent")
	case CUSTOMER:
		if c.PostForm("amount") == "" {
			c.JSON(http.StatusBadRequest, "amount can not be empty")
			return
		}
		conf.AmountStr = c.PostForm("amount")
	default:
		c.JSON(http.StatusBadRequest, ruleType+"is not valid")
		return
	}

	conf.CreatedBy = c.MustGet("id").(string)
	conf.ID = newV4()
	if vemsgs, err := conf.validateReq(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsgs) != 0 {
		c.JSON(http.StatusOK, vemsgs)
		return
	}

	q := conf.composeInsertQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, conf)
}

// better, should limit update fields base on ruletype
func updateLevelConfig(c *gin.Context) {
	createdBy := c.MustGet("id").(string)
	conf := levelRule{
		ID:                    c.Param("id"),
		DiscountStr:           c.PostForm("discount"),
		LevelStr:              c.PostForm("level"),
		AmountStr:             c.PostForm("amount"),
		PiecesStr:             c.PostForm("pieces"),
		ReturnPointPercentStr: c.PostForm("return_point_percent"),
		// RuleType:              c.PostForm("rule_type"), rule type not allowed to change once created
	}
	if vemsgs, err := conf.validateReq(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsgs) != 0 {
		c.JSON(http.StatusOK, vemsgs)
		return
	}
	q := conf.composeUpdateQueryTrack(createdBy)
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, conf.ID)
}

func getAllLevelConfigs(c *gin.Context) {
	rows, err := dbQuery(`SELECT id, discount, level, amount, rule_type, pieces, return_point_percent, 
		created_by,created_at FROM level_rate_rules WHERE rule_type!='RATE' ORDER BY created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	var cs []levelRule
	for rows.Next() {
		var discount, amount, pieces, returnPointPercent sql.NullInt64
		var id, ruleType, createdBy string
		var level sql.NullInt64
		var createdAt time.Time
		if err := rows.Scan(&id, &discount, &level, &amount, &ruleType, &pieces, &returnPointPercent,
			&createdBy, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		conf := levelRule{
			ID:        id,
			Discount:  int(discount.Int64),
			Level:     int(level.Int64),
			Amount:    int(amount.Int64),
			Pieces:    int(pieces.Int64),
			RuleType:  ruleType,
			CreatedBy: createdBy,
			CreatedAt: createdAt,
		}
		cs = append(cs, conf)
	}
	c.JSON(http.StatusOK, cs)
}

func newRateConfig(c *gin.Context) {
	createdBy := c.MustGet("id").(string)
	id := newV4()
	rate, err := util.StringToFloat(c.PostForm("rate"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}
	q := fmt.Sprintf(`INSERT INTO level_rate_rules (id, exchange_rate_float, rule_type, created_by) 
	VALUES ('%s','%f','%s','%s')`, id, rate, "RATE", createdBy)
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	rc := exchangeRateFloat{
		ID:                id,
		ExchangeRateFloat: rate,
		RuleType:          "RATE",
		CreatedBy:         createdBy,
	}
	c.JSON(http.StatusOK, rc)
}

func getRateConfig(c *gin.Context) {
	var rate float64
	var id, createdBy string
	var createdAt time.Time
	q := `SELECT id, exchange_rate_float,created_by,created_at 
	FROM level_rate_rules WHERE rule_type = 'RATE' ORDER BY created_at DESC LIMIT 1`
	if err := dbQueryRow(q).Scan(&id, &rate, &createdBy, &createdAt); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	rc := exchangeRateFloat{
		ID:                id,
		ExchangeRateFloat: rate,
		RuleType:          "RATE",
		CreatedBy:         createdBy,
		CreatedAt:         createdAt,
	}
	c.JSON(http.StatusOK, rc)
}

func getAllRateConfigs(c *gin.Context) {
	rows, err := dbQuery(`SELECT id, exchange_rate_float, created_by,created_at 
		FROM level_rate_rules 
		WHERE rule_type = 'RATE' ORDER BY created_at`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	var cs []exchangeRateFloat
	for rows.Next() {
		var rate float64
		var id, createdBy string
		var createdAt time.Time
		rows.Scan(&id, &rate, &createdBy, &createdAt)
		rc := exchangeRateFloat{
			ID:                id,
			ExchangeRateFloat: rate,
			RuleType:          "RATE",
			CreatedBy:         createdBy,
			CreatedAt:         createdAt,
		}
		cs = append(cs, rc)
	}
	c.JSON(http.StatusOK, cs)
}

func (ac *exchangeRateFloat) getActiveRateConfig() {
	var rate float64
	var id, createdBy string
	var createdAt time.Time
	q := `SELECT id, exchange_rate_float,created_by,created_at 
	FROM level_rate_rules WHERE rule_type = 'RATE' ORDER BY created_at DESC LIMIT 1`
	if err := dbQueryRow(q).Scan(&id, &rate, &createdBy, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			util.Println("fail to get active level_rate_rules, use default level_rate_rules")
			return
		}
	}
	ac.ID = id
	ac.ExchangeRateFloat = rate
	ac.CreatedBy = createdBy
	ac.CreatedAt = createdAt
}

func (ac *levelRule) validateReq() ([]errors.HSMessage, error) {
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
	if ac.LevelStr != "" {
		level, err := strconv.Atoi(ac.LevelStr)
		if err != nil {
			return nil, err
		}
		// LIMIT level to no more than 10
		if level > 10 {
			vemsgAgentDiscountNotValid.Message = "level can not be more that 10"
			vmsgs = append(vmsgs, vemsgAgentDiscountNotValid)
		}
		ac.Level = level
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
	if ac.ReturnPointPercentStr != "" {
		rpp, err := strconv.Atoi(ac.ReturnPointPercentStr)
		if err != nil {
			return nil, err
		}
		ac.ReturnPointPercent = rpp
	}

	if ac.RuleType != "" {
		if !util.IsInArrayString(strings.ToUpper(ac.RuleType), []string{CUSTOMER, AGENT}) {
			vemsgNotValid.Message = "level_rate_rules type not valid"
			vmsgs = append(vmsgs, vemsgNotValid)
		}
		ac.RuleType = strings.ToUpper(ac.RuleType)
	}
	return vmsgs, nil
}

func (ac *levelRule) composeInsertQuery() string {
	params := ac.paramsKV()
	q := `INSERT INTO level_rate_rules (id`
	va := fmt.Sprintf(`VALUES ('%s'`, ac.ID)
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

func (ac *levelRule) composeUpdateQuery() string {
	params := ac.paramsKV()
	q := `UPDATE level_rate_rules SET`
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
	q = fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, ac.ID)
	return q
}

func (ac *levelRule) composeUpdateQueryTrack(updatedBy string) string {
	params := ac.paramsKV()
	q := `UPDATE level_rate_rules SET`
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
	newHistoryRecords(updatedBy, "level_rate_rules", ac.ID, params)
	q = fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, ac.ID)
	return q
}

func (ac *levelRule) paramsKV() map[string]interface{} {
	params := make(map[string]interface{})

	if ac.Level != 0 {
		params["level"] = ac.Level
	}
	if ac.Discount != 0 {
		params["discount"] = ac.Discount
	}
	if ac.Amount != 0 {
		params["amount"] = ac.Amount
	}
	if ac.RuleType != "" {
		params["rule_type"] = ac.RuleType
	}
	if ac.Pieces != 0 {
		params["pieces"] = ac.Pieces
	}
	if ac.ReturnPointPercent != 0 {
		params["return_point_percent"] = ac.ReturnPointPercent
	}
	if ac.CreatedBy != "" {
		params["created_by"] = ac.CreatedBy
	}
	return params
}
