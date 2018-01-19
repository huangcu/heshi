package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"util"

	"github.com/satori/go.uuid"
)

func (d *diamond) composeInsertQuery() string {
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
	return q
}

func (d *diamond) composeUpdateQuery() string {
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
	return q
}

// 	params := make(map[string]interface{})
func (d *diamond) parmsKV() map[string]interface{} {
	params := make(map[string]interface{})
	if d.StockRef != "" {
		params["stock_ref"] = d.StockRef
	}
	if d.Shape != "" {
		params["shape"] = d.Shape
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
	if d.GradingLab != "" {
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

func importDiamondsCustomizeHeaders(headers map[string]string, records [][]string) ([][]string, error) {
	originalHeaders := []string{}
	ignoredRows := [][]string{}
	//get headers
	for index := 0; index < len(records); index++ {
		if index == 0 {
			originalHeaders = records[0]
		}
	}
	for index := 0; index < len(records); index++ {
		//process records
		if index != 0 {
			ignored := false
			d := diamond{}
			record := records[index]
			util.Println("processsing " + strconv.Itoa(index))
			for header, oriheader := range headers {
				for i := 0; i < len(originalHeaders); i++ {
					if originalHeaders[i] == oriheader {
						switch header {
						case "diamond_id":
							d.DiamondID = record[i]
						case "stock_ref":
							d.StockRef = record[i]
						case "shape":
							d.Shape = record[i]
						case "carat":
							cValue, err := strconv.ParseFloat(record[i], 64)
							if err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							}
							if cValue == 0 {
								ignored = true
							}
							d.Carat = cValue
						case "color":
							d.Color = record[i]
						case "clarity":
							d.Clarity = record[i]
						case "grading_lab":
							d.GradingLab = record[i]
						case "certificate_number":
							d.CertificateNumber = record[i]
						case "cut_grade":
							d.CutGrade = record[i]
						case "polish":
							d.Polish = record[i]
						case "symmetry":
							d.Symmetry = record[i]
						case "fluorescence_intensity":
							d.FluorescenceIntensity = record[i]
						case "country":
							d.Country = record[i]
						case "supplier":
							d.Supplier = record[i]
						case "price_no_added_value":
							cValue, err := strconv.ParseFloat(strings.Replace(record[i], ",", "", -1), 64)
							if err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							}
							if cValue == 0 {
								ignored = true
							}
							d.PriceNoAddedValue = cValue
						case "price_retail":
							cValue, err := strconv.ParseFloat(strings.Replace(record[i], ",", "", -1), 64)
							if err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							}
							if cValue == 0 {
								ignored = true
							}
							d.PriceRetail = cValue
						case "clarity_number":
							d.ClarityNumber = record[i]
						case "cut_number":
							d.CutNumber = record[i]
						}
						break
					}
				}
			}
			//insert into db
			if !ignored {
				var id string
				if err := dbQueryRow(fmt.Sprintf("SELECT id FROM diamonds WHERE stock_ref='%s'", d.StockRef)).Scan(&id); err != nil {
					if err == sql.ErrNoRows {
						d.ID = uuid.NewV4().String()
						q := d.composeInsertQuery()
						if _, err := dbExec(q); err != nil {
							return nil, err
							// ignoredRows = append(ignoredRows, record)
						}
					} else {
						// ignoredRows = append(ignoredRows, record)
						return nil, err
					}
				}

				d.ID = id
				q := d.composeUpdateQuery()
				if _, err := dbExec(q); err != nil {
					// ignoredRows = append(ignoredRows, record)
					return nil, err
				}
			}
			util.Println("finish process")
		}
	}
	return ignoredRows, nil
}
