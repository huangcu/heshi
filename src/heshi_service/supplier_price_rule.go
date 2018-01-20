package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type PriceSetting struct {
	ID              string  `json:"id"`
	SupplierID      string  `json:"supplier_id"`
	CaratFrom       float64 `json:"carat_from"`
	CaratFromStr    string  `json:"-"`
	CaratTo         float64 `json:"carat_to"`
	CaratToStr      string  `json:"-"`
	Colors          string  `json:"color"`
	Clarities       string  `json:"clarity"`
	Cuts            string  `json:"cut"`
	Symmetries      string  `json:"symmetry"`
	Polishs         string  `json:"polish"`
	Fluos           string  `json:"fluo"`
	Certificates    string  `json:"certificate"`
	TheParaValue    float64 `json:"the_para_value"`
	TheParaValueStr string  `json:"-"`
	Priority        int     `json:"priority"`
	PriorityStr     string  `json:"-"`
	Status          string  `json:"status"`
}

func addPriceRule(c *gin.Context) {
	ps := PriceSetting{
		SupplierID:      c.PostForm("supplier_id"),
		CaratFromStr:    c.PostForm("carat_from"),
		CaratToStr:      c.PostForm("carat_to"),
		Colors:          c.PostForm("color"),
		Clarities:       c.PostForm("clarity"),
		Cuts:            c.PostForm("cut"),
		Symmetries:      c.PostForm("symmetry"),
		Polishs:         c.PostForm("polish"),
		Fluos:           c.PostForm("fluo"),
		Certificates:    c.PostForm("certificate"),
		TheParaValueStr: c.PostForm("the_para_value"),
		PriorityStr:     c.PostForm("priority"),
	}

	if vemsg, err := ps.validatePriceSetting(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsg) != 0 {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	ps.ID = uuid.NewV4().String()
	q := ps.composeInsertQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, ps.ID)
}

func updatePriceRule(c *gin.Context) {
	s := PriceSetting{
		ID:              c.Param("id"),
		CaratFromStr:    c.PostForm("carat_from"),
		CaratToStr:      c.PostForm("carat_to"),
		Colors:          c.PostForm("color"),
		Clarities:       c.PostForm("clarity"),
		Cuts:            c.PostForm("cut"),
		Symmetries:      c.PostForm("symmetry"),
		Polishs:         c.PostForm("polish"),
		Fluos:           c.PostForm("fluo"),
		Certificates:    c.PostForm("certificate"),
		TheParaValueStr: c.PostForm("the_para_value"),
		PriorityStr:     c.PostForm("priority"),
	}

	if vemsg, err := s.validatePriceSetting(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsg) != 0 {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	q := s.composeUpdateQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, s.ID)
}

func disablePriceRule(c *gin.Context) {
	q := fmt.Sprintf("UPDATE price_settings_universal SET status='disabled' WHERE id='%s'", c.Param("id"))
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, "SUCCESS")
}

func getPriceRule(c *gin.Context) {
	q := selectPriceRulesQuery(c.Param("id"))
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ds, err := composePriceSetting(rows)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, fmt.Sprintf("Fail to find price rule with id: %s", c.Param("id")))
			return
		}
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ds)
}

func getAllPriceRule(c *gin.Context) {
	q := selectPriceRulesQuery("")
	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ds, err := composePriceSetting(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ds)
}

func selectPriceRulesQuery(id string) string {
	q := `SELECT id,supplier_id,carat_from,carat_to,color,clarity,cut,symmetry,polish,fluo,
	certificate,the_para_value,priority,status FROM price_settings_universal`

	if id != "" {
		q = fmt.Sprintf("%s WHERE id='%s'", q, id)
	}
	return q
}

func composePriceSetting(rows *sql.Rows) ([]PriceSetting, error) {
	var id, supplierID, color, clarity, cut, symmetry, polish, fluo, certificate, status string
	var caratFrom, caratTo, theParaValue float64
	var priority int

	var ps []PriceSetting
	for rows.Next() {
		if err := rows.Scan(&id, &supplierID, &caratFrom, &caratTo, &color, &clarity, &cut, &symmetry,
			&polish, &fluo, &certificate, &theParaValue, &priority, &status); err != nil {
			return nil, err
		}
		p := PriceSetting{
			ID:           id,
			SupplierID:   supplierID,
			CaratFrom:    caratFrom,
			CaratTo:      caratTo,
			Colors:       color,
			Clarities:    clarity,
			Cuts:         cut,
			Symmetries:   symmetry,
			Polishs:      polish,
			Fluos:        fluo,
			Certificates: certificate,
			TheParaValue: theParaValue,
			Priority:     priority,
			Status:       status,
		}
		ps = append(ps, p)
	}
	return ps, nil
}
