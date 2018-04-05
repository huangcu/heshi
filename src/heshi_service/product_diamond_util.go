package main

import (
	"fmt"
	"heshi/errors"
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
	q = fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, d.ID)
	return q
}

// 	params := make(map[string]interface{})
func (d *diamond) parmsKV() map[string]interface{} {
	params := make(map[string]interface{})
	if d.DiamondID != "" {
		params["diamond_id"] = d.DiamondID
	}
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
	if d.RecommendWords != "" {
		params["recommend_words"] = d.RecommendWords
	}
	if d.ExtraWords != "" {
		params["extra_words"] = d.ExtraWords
	}
	if len(d.Images) != 0 {
		params["images"] = strings.Join(d.Images, ";")
	}
	if d.Status != "" {
		params["status"] = d.Status
	}
	if d.OrderedBy != "" {
		params["ordered_by"] = d.OrderedBy
	}
	if d.PickedUp != "" {
		params["picked_up"] = d.PickedUp
	}
	if d.SoldPrice != 0 {
		params["sold_price"] = d.SoldPrice
	}
	if d.Profitable != "" {
		params["profitable"] = d.Profitable
	}
	return params
}

// TODO to be removed
func importDiamondsCustomizeHeaders(headers map[string]string, records [][]string) ([][]string, error) {
	oldStockRefList, err := getAllStockRef()
	if err != nil {
		return nil, err
	}
	var suppliers []string
	suppliers, err = getAllValidSupplierName()
	if err != nil {
		util.Printf("Fail to get all suppliers name from db, use default predefined: %s",
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
	util.Println("start process diamond")
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
							d.DiamondID = strings.ToUpper(record[i])
						case "stock_ref":
							d.StockRef = strings.ToUpper(record[i])
						case "shape":
							if s, err := diamondShape(record[i]); err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							} else {
								d.Shape = s
							}
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
							if c, err := diamondColor(record[i]); err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							} else {
								d.Color = c
							}
						case "clarity":
							if c, err := diamondClarity(record[i]); err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							} else {
								d.Clarity = c
							}
						case "grading_lab":
							if s, err := diamondGradingLab(record[i]); err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							} else {
								d.GradingLab = s
							}
						case "certificate_number":
							d.CertificateNumber = strings.ToUpper(record[i])
						case "cut_grade":
							if c, err := diamondCutGradeSymmetryPolish(record[i]); err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							} else {
								d.CutGrade = c
							}
						case "polish":
							if c, err := diamondCutGradeSymmetryPolish(record[i]); err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							} else {
								d.Polish = c
							}
						case "symmetry":
							if c, err := diamondCutGradeSymmetryPolish(record[i]); err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							} else {
								d.Symmetry = c
							}
						case "fluorescence_intensity":
							if c, err := diamondFluo(record[i]); err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							} else {
								d.FluorescenceIntensity = c
							}
						case "country":
							d.Country = strings.ToUpper(record[i])
						case "supplier":
							if s, err := diamondSupplier(record[i], suppliers); err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							} else {
								d.Supplier = s
							}
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
						case "price_retail":
							cValue, err := util.StringToFloat(record[i])
							if err != nil {
								ignoredRows = append(ignoredRows, record)
								ignored = true
							}
							if cValue == 0 {
								ignored = true
							}
							d.PriceRetail = cValue
						case "featured":
							d.Featured = strings.ToUpper(record[i])
						case "recommend_words":
							d.Featured = strings.ToUpper(record[i])
						case "extra_words":
							d.Featured = strings.ToUpper(record[i])
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

func (d *diamond) validateDiamondReq(update bool) ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage
	if !update && d.CaratStr == "" {
		vemsg = append(vemsg, vemsgDiamondCaratEmpty)
	} else if d.CaratStr != "" {
		cValue, err := util.StringToFloat(d.CaratStr)
		if err != nil {
			vemsg = append(vemsg, vemsgDiamondCaratNotValid)
		} else if cValue == 0 {
			vemsg = append(vemsg, vemsgDiamondCaratNotValid)
		} else {
			d.Carat = cValue
		}
	}

	if !update && d.PriceNoAddedValueStr == "" {
		vemsg = append(vemsg, vemsgDiamondRawPriceEmpty)
	} else if d.PriceNoAddedValueStr != "" {
		pValue, err := util.StringToFloat(d.PriceNoAddedValueStr)
		if err != nil {
			vemsg = append(vemsg, vemsgDiamondRawPriceNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgDiamondRawPriceNotValid)
		} else {
			d.PriceNoAddedValue = pValue
		}
	}

	if !update && d.PriceRetailStr == "" {
		vemsg = append(vemsg, vemsgDiamondRetailPriceEmpty)
	} else if d.PriceRetailStr != "" {
		pValue, err := util.StringToFloat(d.PriceRetailStr)
		if err != nil {
			vemsg = append(vemsg, vemsgDiamondRetailPriceNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgDiamondRetailPriceNotValid)
		} else {
			d.PriceRetail = pValue
		}
	}

	if !update && d.StockRef == "" {
		vemsgEmpty.Message = "diamond stock ref can not be empty"
		vemsg = append(vemsg, vemsgEmpty)
	}
	if d.StockRef != "" {
		if err := d.composeStockRefWithSupplierPrefix(); err != nil {
			return nil, err
		}
		if update {
			if exist, err := isItemExistInDbByPropertyWithDifferentID("diamonds", "stock_ref", d.StockRef, d.ID); err != nil {
				return nil, err
			} else if exist {
				vemsgAlreadyExist.Message = "diamond stock_ref " + d.StockRef + " already exists"
				vemsg = append(vemsg, vemsgAlreadyExist)
			}
		} else {
			if exist, err := isItemExistInDbByProperty("diamonds", "stock_ref", d.StockRef); err != nil {
				return nil, err
			} else if exist {
				vemsgAlreadyExist.Message = "diamond stock_ref " + d.StockRef + " already exists"
				vemsg = append(vemsg, vemsgAlreadyExist)
			}
		}
	}
	if !update && d.DiamondID == "" {
		vemsgEmpty.Message = "diamond id can not be empty"
		vemsg = append(vemsg, vemsgNotValid)
	}
	if d.DiamondID != "" {
		if update {
			if exist, err := isItemExistInDbByPropertyWithDifferentID("diamonds", "diamond_id", d.DiamondID, d.ID); err != nil {
				return nil, err
			} else if exist {
				vemsgAlreadyExist.Message = "diamond id " + d.DiamondID + " already exists"
				vemsg = append(vemsg, vemsgAlreadyExist)
			}
		} else {
			if exist, err := isItemExistInDbByProperty("diamonds", "diamond_id", d.DiamondID); err != nil {
				return nil, err
			} else if exist {
				vemsgAlreadyExist.Message = "diamond id " + d.DiamondID + " already exists"
				vemsg = append(vemsg, vemsgAlreadyExist)
			}
		}
	}

	if !update && d.Shape == "" {
		vemsgEmpty.Message = "diamond shape can not be empty"
		vemsg = append(vemsg, vemsgEmpty)
	} else if d.Shape != "" {
		if s, err := diamondShape(d.Shape); err != nil {
			vemsgNotValid.Message = err.Error()
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.Shape = s
		}
	}

	if !update && d.GradingLab == "" {
		vemsgEmpty.Message = "diamond grading lab can not be empty"
		vemsg = append(vemsg, vemsgEmpty)
	} else if d.GradingLab != "" {
		if s, err := diamondGradingLab(d.GradingLab); err != nil {
			vemsgNotValid.Message = err.Error()
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.GradingLab = s
		}
	}

	if !update && d.Color == "" {
		vemsgEmpty.Message = "diamond color can not be empty"
		vemsg = append(vemsg, vemsgEmpty)
	} else if d.Color != "" {
		if d.Color != "" {
			if c, err := diamondColor(d.Color); err != nil {
				vemsgNotValid.Message = err.Error()
				vemsg = append(vemsg, vemsgNotValid)
			} else {
				d.Color = c
			}
		}
	}

	if !update && d.Clarity == "" {
		vemsgEmpty.Message = "diamond clarity can not be empty"
		vemsg = append(vemsg, vemsgEmpty)
	} else if d.Clarity != "" {
		if c, err := diamondClarity(d.Clarity); err != nil {
			vemsgNotValid.Message = err.Error()
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.Clarity = c
		}
	}

	if !update && d.CutGrade == "" {
		vemsgEmpty.Message = "diamond cut grade can not be empty"
		vemsg = append(vemsg, vemsgEmpty)
	} else if d.CutGrade != "" {
		if c, err := diamondCutGradeSymmetryPolish(d.CutGrade); err != nil {
			vemsgNotValid.Message = err.Error()
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.CutGrade = c
		}
	}

	if !update && d.Polish == "" {
		vemsgEmpty.Message = "diamond polish can not be empty"
		vemsg = append(vemsg, vemsgEmpty)
	} else if d.Polish != "" {
		if c, err := diamondCutGradeSymmetryPolish(d.Polish); err != nil {
			vemsgNotValid.Message = err.Error()
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.Polish = c
		}
	}

	if !update && d.Symmetry == "" {
		vemsgEmpty.Message = "diamond symmetry can not be empty"
		vemsg = append(vemsg, vemsgEmpty)
	} else if d.Symmetry != "" {
		if c, err := diamondCutGradeSymmetryPolish(d.Symmetry); err != nil {
			vemsgNotValid.Message = err.Error()
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.Symmetry = c
		}
	}

	if !update && d.FluorescenceIntensity == "" {
		vemsgEmpty.Message = "diamond fluorescence intensity can not be empty"
		vemsg = append(vemsg, vemsgEmpty)
	} else if d.FluorescenceIntensity != "" {
		if c, err := diamondFluo(d.FluorescenceIntensity); err != nil {
			vemsgNotValid.Message = err.Error()
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.FluorescenceIntensity = c
		}
	}

	if !update && d.Country == "" {
		vemsgEmpty.Message = "diamond country can not be empty"
		vemsg = append(vemsg, vemsgEmpty)
	} else if d.Country != "" {
		if c, err := diamondCountry(d.Country); err != nil {
			vemsgNotValid.Message = err.Error()
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.Country = c
		}
	}

	if !update && d.Supplier == "" {
		vemsgEmpty.Message = "diamond supplier can not be empty"
		vemsg = append(vemsg, vemsgEmpty)
	} else if d.Supplier != "" {
		if s, err := diamondSupplierPageReq(d.Supplier); err != nil {
			vemsgEmpty.Message = errors.GetMessage(err)
			vemsg = append(vemsg, vemsgEmpty)
		} else {
			d.Supplier = s
		}
	}

	//TODO Status - when new, the status is always available??
	//TODO Featured/ Profitable - value can only be YES OR NO
	return vemsg, nil
}

//TODO TOBE REMOVED
func (d *diamond) validateDiamondUpdateReq() ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage
	if d.CaratStr != "" {
		cValue, err := util.StringToFloat(d.CaratStr)
		if err != nil {
			vemsg = append(vemsg, vemsgDiamondCaratNotValid)
		} else if cValue == 0 {
			vemsg = append(vemsg, vemsgDiamondCaratNotValid)
		} else {
			d.Carat = cValue
		}
	}

	if d.PriceNoAddedValueStr != "" {
		pValue, err := util.StringToFloat(d.PriceNoAddedValueStr)
		if err != nil {
			vemsg = append(vemsg, vemsgDiamondRawPriceNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgDiamondRawPriceNotValid)
		} else {
			d.PriceNoAddedValue = pValue
		}
	}

	if d.PriceRetailStr == "" {
		pValue, err := util.StringToFloat(d.PriceRetailStr)
		if err != nil {
			vemsg = append(vemsg, vemsgDiamondRetailPriceNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgDiamondRetailPriceNotValid)
		} else {
			d.PriceRetail = pValue
		}
	}
	if d.StockRef != "" {
		if exist, err := isItemExistInDbByPropertyWithDifferentID("diamonds", "stock_ref", d.StockRef, d.ID); err != nil {
			return nil, err
		} else if exist {
			vemsgAlreadyExist.Message = "diamond stock_ref " + d.StockRef + " already exists"
			vemsg = append(vemsg, vemsgAlreadyExist)
		}
	}
	if d.DiamondID != "" {
		if exist, err := isItemExistInDbByPropertyWithDifferentID("diamonds", "diamond_id", d.DiamondID, d.ID); err != nil {
			return nil, err
		} else if exist {
			vemsgAlreadyExist.Message = "diamond id " + d.DiamondID + " already exists"
			vemsg = append(vemsg, vemsgAlreadyExist)
		}
	}
	if d.Shape != "" {
		s, err := diamondShape(d.Shape)
		if err != nil {
			return nil, err
		}
		d.Shape = s
	}

	if d.Color != "" {
		if c, err := diamondColor(d.Color); err != nil {
			vemsgNotValid.Message = "diamond color is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.Color = c
		}
	}

	if d.Clarity != "" {
		if c, err := diamondClarity(d.Clarity); err != nil {
			vemsgNotValid.Message = "diamond clarity is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.Clarity = c
		}
	}

	if d.CutGrade != "" {
		if c, err := diamondCutGradeSymmetryPolish(d.CutGrade); err != nil {
			vemsgNotValid.Message = "diamond cut grade is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.CutGrade = c
		}
	}

	if d.Polish != "" {
		if c, err := diamondCutGradeSymmetryPolish(d.Polish); err != nil {
			vemsgNotValid.Message = "diamond polish is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.Polish = c
		}
	}

	if d.Symmetry != "" {
		if c, err := diamondCutGradeSymmetryPolish(d.Symmetry); err != nil {
			vemsgNotValid.Message = "diamond symmetry is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.Symmetry = c
		}
	}

	if d.FluorescenceIntensity != "" {
		if c, err := diamondFluo(d.FluorescenceIntensity); err != nil {
			vemsgNotValid.Message = "diamond fluorescence is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.FluorescenceIntensity = c
		}
	}

	if d.Supplier != "" {
		if s, err := diamondSupplierPageReq(d.Supplier); err != nil {
			vemsgNotValid.Message = errors.GetMessage(err)
			vemsg = append(vemsg, vemsgNotValid)
		} else {
			d.Supplier = s
		}
	}

	//TODO Status - when new, the status is always available??
	//TODO Featured/ Profitable - value can only be YES OR NO
	return vemsg, nil
}

func diamondSupplierPageReq(supplier string) (string, error) {
	suppliers, err := getAllValidSupplierName()
	if err != nil {
		util.Traceln("Fail to get all suppliers name from db, use default predefined: %s",
			strings.Join(VALID_SUPPLIER_NAME, ","))
		suppliers = VALID_SUPPLIER_NAME
	}
	if len(supplier) != 0 {
		if util.IsInArrayString(strings.ToUpper(supplier), suppliers) {
			return strings.ToUpper(supplier), nil
		}
	}
	return "", errors.Newf("supplier %s not exist, please add first!", supplier)
}

//TODO compose stock ref with supplier prefix?????
func diamondStockRef(stockRef, supplierPrefix string) string {
	return supplierPrefix + stockRef
}

func (d *diamond) isDiamondExistByDiamondID() error {
	var id string
	if err := dbQueryRow("SELECT id FROM diamonds WHERE diamond_id=?", d.DiamondID).Scan(&id); err != nil {
		return err
	}
	d.ID = id
	return nil
}

func isDiamondExistByID(id string) (bool, error) {
	var count int
	if err := dbQueryRow("SELECT COUNT(*) FROM diamonds WHERE id=?", id).Scan(&count); err != nil {
		return false, err
	}
	return count == 1, nil
}
