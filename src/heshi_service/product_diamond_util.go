package main

import (
	"fmt"
	"strconv"
	"strings"
	"util"
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
	oldStockRefList, err := getAllStockRef()
	if err != nil {
		return nil, err
	}
	var suppliers []string
	suppliers, err = getAllValidSupplierName()
	if err != nil {
		util.Traceln("Fail to get all suppliers name from db, use default predefined: %s",
			strings.Join(VALID_SUPPLIER_NAME, ","))
		suppliers = VALID_SUPPLIER_NAME
	}
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
							d.Shape = diamondShape(record[i])
						case "carat":
							cValue, err := util.StringToFloat(record[i])
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
							d.Clarity = diamondClarity(record[i])
						case "grading_lab":
							d.GradingLab = diamondGradingLab(record[i])
						case "certificate_number":
							d.CertificateNumber = record[i]
						case "cut_grade":
							d.CutGrade = diamondCutGradeSymmetryPolish(record[i])
						case "polish":
							d.Polish = diamondCutGradeSymmetryPolish(record[i])
						case "symmetry":
							d.Symmetry = diamondCutGradeSymmetryPolish(record[i])
						case "fluorescence_intensity":
							d.FluorescenceIntensity = diamondFluo(record[i])
						case "country":
							d.Country = record[i]
						case "supplier":
							d.Supplier = diamondSupplier(record[i], suppliers)
						case "price_no_added_value":
							cValue, err := util.StringToFloat(record[i])
							if err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							}
							if cValue == 0 {
								ignored = true
							}
							d.PriceNoAddedValue = cValue
						}
						break
					}
				}
			}
			//insert into db
			if !ignored {
				if err := d.composeStockRefWithSupplierPrefix(); err != nil {
					//TODO
					return nil, err
				}
				if err := d.processDiamondRecord(); err != nil {
					//TODO return err for now!
					return nil, err
				}
				//remove it from old stock ref map
				delete(oldStockRefList, d.StockRef)
			}
		}
	}
	util.Println("finish process diamond")
	if err := offlineDiamondsNoLongerExist(oldStockRefList); err != nil {
		return ignoredRows, err
	}
	return ignoredRows, nil
}

//TODO double check url GIA
func composeCertifcateLink(gradingLab, certificate string) string {
	switch gradingLab {
	case "HRD":
		return fmt.Sprintf("https://my.hrdantwerp.com/?L=&record_number=%s&certificatetype=MC", certificate)
	case "GIA":
		return "http://www.gia.edu/cs/Satellite?pagename=GST%2FDispatcher&childpagename=GIA%2FPage%2FReportCheck&c=Page&cid=1355954554547&reportno=" + certificate
	case "IGI":
		return fmt.Sprintf("http://www.igiworldwide.com/verify.php?r=%s", certificate)
	default:
		return ""
	}
}
