package main

import (
	"fmt"
	"strings"
)

func (d *diamond) composeInsertQuery() (string, error) {
	params := d.parmsKV()
	q := `INSERT INTO diamonds (id`
	va := fmt.Sprintf(`VALUES ('%s'`, d.ID)
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
		}
	}
	q = fmt.Sprintf("%s) %s)", q, va)
	return q, nil
}

func (d *diamond) composeUpdateQuery() (string, error) {
	params := d.parmsKV()
	q := `UPDATE diamonds SET`
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
	q = fmt.Sprintf("%s WHERE id='%s'", strings.TrimSuffix(q, ","), d.ID)
	return q, nil
}

// 	params := make(map[string]interface{})
func (d *diamond) parmsKV() map[string]interface{} {
	params := make(map[string]interface{})
	if d.StockRef != "" {
		params["stock_ref"] = d.StockRef
	}
	if d.Shape != "" {
		params["cellphone"] = d.Shape
	}
	if d.Carat != 0 {
		params["carat"] = d.Carat
	}
	if d.Color != "" {
		params["color"] = d.Color
	}
	if d.Clarity != "" {
		params["clarity"] = d.Clarity
	}
	if d.GradingLab != 0 {
		params["grading_lab"] = d.GradingLab
	}
	if d.CertificateNumber != "" {
		params["certificate_number"] = d.CertificateNumber
	}
	if d.CutGrade != "" {
		params["cut_grade"] = d.CutGrade
	}
	if d.Polish != "" {
		params["polish"] = d.Polish
	}
	if d.Symmetry != "" {
		params["symmetry"] = d.Symmetry
	}
	if d.FluorescenceIntensity != "" {
		params["fluorescence_intensity"] = d.FluorescenceIntensity
	}
	if d.Country != "" {
		params["country"] = d.Country
	}
	if d.Supplier != "" {
		params["supplier"] = d.Supplier
	}
	if d.PriceNoAddedValue != 0 {
		params["price_no_added_value"] = d.PriceNoAddedValue
	}
	if d.PriceRetail != 0 {
		params["price_retail"] = d.PriceRetail
	}
	if d.CertificateLink != "" {
		params["certificate_link"] = d.CertificateLink
	}
	if d.ClarityNumber != "" {
		params["clarity_number"] = d.ClarityNumber
	}
	if d.CutNumber != "" {
		params["cut_number"] = d.CutNumber
	}
	if d.Featured != "" {
		params["featured"] = d.Featured
	}
	if d.RecommandWords != "" {
		params["recommand_words"] = d.RecommandWords
	}
	if d.ExtraWords != "" {
		params["extra_words"] = d.ExtraWords
	}
	if d.Status != "" {
		params["status"] = d.Status
	}
	if d.OrderedBy != 0 {
		params["ordered_by"] = d.OrderedBy
	}
	if d.PickedUp != "" {
		params["picked_up"] = d.PickedUp
	}
	if d.Sold != "" {
		params["sold"] = d.Sold
	}
	if d.SoldPrice != 0 {
		params["sold_price"] = d.SoldPrice
	}
	if d.Profitable != "" {
		params["profitable"] = d.Profitable
	}
	return params
}
