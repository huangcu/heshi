package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"math"
	"strconv"
	"strings"
	"util"

	uuid "github.com/satori/go.uuid"
)

func validateDiamondHeaders(headers []string) []string {
	var missingHeaders []string
	for _, header := range diamondHeaders {
		if !util.IsInArrayString(header, headers) {
			missingHeaders = append(missingHeaders, header)
		}
	}
	return missingHeaders
}

func importDiamondProducts(file string) ([][]string, error) {
	originalHeaders := []string{}
	records, err := util.ParseCSVToArrays(file)
	if err != nil {
		return nil, err
	}
	if len(records) < 1 {
		return nil, errors.New("uploaded file has no records")
	}

	ignoredRows := [][]string{}
	//get headers
	originalHeaders = records[0]

	//process records
	for index := 1; index < len(records); index++ {
		ignored := false
		d := diamond{}
		record := records[index]
		util.Printf("processsing row: %d, %s", index, record)
		for i, header := range originalHeaders {
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
				d.Carat = math.Abs(cValue)
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
				d.PriceNoAddedValue = math.Abs(cValue)
				//THIS PRICE IS CACULATED BASE ON PRICE SETTING
			// case "price_retail":
			// 	cValue, err := strconv.ParseFloat(strings.Replace(record[i], ",", "", -1), 64)
			// 	if err != nil {
			// 		ignoredRows = append(ignoredRows, record)
			// 		ignored = true
			// 	}
			// 	if cValue == 0 {
			// 		ignored = true
			// 	}
			// 	d.PriceRetail = math.Abs(cValue)
			case "clarity_number":
				d.ClarityNumber = record[i]
			case "cut_number":
				d.CutNumber = record[i]
			}
		}
		//handle db
		if !ignored {
			if err := d.processDiamondRecord(); err != nil {
				//TODO return err for now!
				return nil, err
			}
		}
	}
	util.Println("finish process diamond")
	return ignoredRows, nil
}

func (d *diamond) processDiamondRecord() error {
	q := "SELECT price_no_added_value, price_retail FROM diamonds WHERE grading_lab = '' AND certificate_number = ''"
	var id, status string
	var priceNoAddedValue, priceRetail float64
	q = fmt.Sprintf("SELECT id, price_no_added_value, price_retail, status FROM diamonds WHERE stock_ref='%s'", d.StockRef)
	if err := db.QueryRow(q).Scan(&id, &priceNoAddedValue, &priceRetail, &status); err != nil {
		//item not exist in db
		if err == sql.ErrNoRows {
			d.ID = uuid.NewV4().String()
			q := d.composeInsertQuery()
			if _, err := dbExec(q); err != nil {
				util.Tracef(`fail to add diamond item. diamond: %s; certificate_number: %s; grading_lab: %s; retail price %f`,
					d.StockRef, d.CertificateNumber, d.GradingLab, d.PriceRetail)
				return err
			}
			util.Tracef(`diamond item added! diamond: %s; certificate_number: %s; grading_lab: %s; retail price %f`,
				d.StockRef, d.CertificateNumber, d.GradingLab, d.PriceRetail)
			return nil
		}
		return err
	}
	//item alread exist in db
	if err := d.processPrice(); err != nil {
		return err
	}
	if status != "SOLD" && status != "RESERVED" && (d.PriceRetail-priceRetail) > 5 {
		q := d.composeInsertQuery()
		if _, err := dbExec(q); err != nil {
			util.Tracef(`retail price changed, but failed to update. diamond: %s; certificate_number: %s; grading_lab: %s; original price: %f; new price should be %f`,
				d.StockRef, d.CertificateNumber, d.GradingLab, priceRetail, d.PriceRetail)
			return err
		}
		util.Tracef(`retail price changed for diamond: %s; certificate_number: %s; grading_lab: %s; original price: %f; new price %f`,
			d.StockRef, d.CertificateNumber, d.GradingLab, priceRetail, d.PriceRetail)
	}

	return nil
}

func (d *diamond) processPrice() error {
	q := fmt.Sprintf(`SELECT carat_from, carat_to, color, clarity, cut_grade, polish, 
		symmetry, grading_lab, fluo, shape, the_para_value FROM price_settings_universal 
		WHERE supplier_id = '%s' AND status='active' ORDER BY priority ASC`, d.Supplier)
	rows, err := dbQuery(q)
	if err != nil {
		return err
	}
	for rows.Next() {
		var color, clarity, cutGrade, polish, symmetry, gradingLab, fluo, shape string
		var caratFrom, caratTo, theParaValue float64
		if err := rows.Scan(&caratFrom, &caratTo, &color, &clarity, &cutGrade, &polish, &symmetry,
			&gradingLab, &fluo, &shape, &theParaValue); err != nil {
			return err
		}

		if d.Shape == "BR" {
			if d.Carat > caratFrom && d.Carat <= caratTo && util.IsInArrayString(d.Color, strings.Split(color, ",")) &&
				util.IsInArrayString(d.Clarity, strings.Split(clarity, ",")) &&
				util.IsInArrayString(d.CutGrade, strings.Split(cutGrade, ",")) &&
				util.IsInArrayString(d.Polish, strings.Split(polish, ",")) &&
				util.IsInArrayString(d.Symmetry, strings.Split(symmetry, ",")) &&
				util.IsInArrayString(d.FluorescenceIntensity, strings.Split(fluo, ",")) &&
				// util.IsInArrayString(d.Shape, strings.Split(shape, ",")) &&
				util.IsInArrayString(d.GradingLab, strings.Split(gradingLab, ",")) {
				d.PriceRetail = d.PriceNoAddedValue * theParaValue
			}
		} else {
			if d.Carat > caratFrom && d.Carat <= caratTo && util.IsInArrayString(d.Color, strings.Split(color, ",")) &&
				util.IsInArrayString(d.Clarity, strings.Split(clarity, ",")) &&
				util.IsInArrayString(d.Polish, strings.Split(polish, ",")) &&
				util.IsInArrayString(d.Symmetry, strings.Split(symmetry, ",")) &&
				util.IsInArrayString(d.FluorescenceIntensity, strings.Split(fluo, ",")) &&
				util.IsInArrayString(d.Shape, strings.Split(shape, ",")) &&
				util.IsInArrayString(d.GradingLab, strings.Split(gradingLab, ",")) {
				d.PriceRetail = d.PriceNoAddedValue * theParaValue
			}
		}
		if d.PriceRetail != 0 {
			return nil
		}
	}
	return nil
}
