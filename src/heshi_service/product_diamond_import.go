package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"strings"
	"util"
)

func validateDiamondHeaders(headers []string) []string {
	var missingHeaders []string
	for k, header := range diamondHeaders {
		if !util.IsInArrayString(header, headers) && k < 15 {
			missingHeaders = append(missingHeaders, header)
		}
	}
	return missingHeaders
}

//TODO better validate import(new VS update data validation - > compare with jewelrys)
func importDiamondProducts(file string) ([]util.Row, error) {
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

	rows, err := util.ParseCSVToStruct(file)
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, errors.New("uploaded file has no records")
	}

	unimportRows := []util.Row{}
	//get headers
	originalHeaders := rows[0]

	//process records
	for index := 1; index < len(rows); index++ {
		d := diamond{}
		row := rows[index]
		record := row.Value
		util.Printf("processsing row: %d, %s", index, record)
		for i, header := range originalHeaders.Value {
			switch header {
			case "diamond_id":
				if record[i] == "" {
					row.Ignored = true
					break
				}
				d.DiamondID = strings.ToUpper(record[i])
			case "stock_ref":
				if record[i] == "" {
					row.Ignored = true
					break
				}
				d.StockRef = strings.ToUpper(record[i])
			case "shape":
				if s, err := diamondShape(record[i]); err != nil {
					row.Message = append(row.Message, errors.GetMessage(err))
					row.Ignored = true
				} else {
					d.Shape = s
				}
			case "carat":
				cValue, err := util.StringToFloat(record[i])
				if err != nil {
					row.Message = append(row.Message, errors.GetMessage(err))
					row.Ignored = true
				}
				if cValue == 0 {
					row.Ignored = true
				}
				d.Carat = cValue
			case "color":
				d.Color = diamondColor(record[i])
			case "clarity":
				d.Clarity = diamondClarity(record[i])
			case "grading_lab":
				if s, err := diamondGradingLab(record[i]); err != nil {
					row.Message = append(row.Message, errors.GetMessage(err))
					row.Ignored = true
				} else {
					d.GradingLab = s
				}
				//TODO certificate number duplicate??
			case "certificate_number":
				d.CertificateNumber = strings.ToUpper(record[i])
			case "cut_grade":
				d.CutGrade = diamondCutGradeSymmetryPolish(record[i])
			case "polish":
				d.Polish = diamondCutGradeSymmetryPolish(record[i])
			case "symmetry":
				d.Symmetry = diamondCutGradeSymmetryPolish(record[i])
			case "fluorescence_intensity":
				d.FluorescenceIntensity = diamondFluo(record[i])
			case "country":
				//TODO format country
				d.Country = strings.ToUpper(record[i])
			case "supplier":
				if s, err := diamondSupplier(record[i], suppliers); err != nil {
					row.Message = append(row.Message, errors.GetMessage(err))
					row.Ignored = true
				} else {
					d.Supplier = s
				}
				// 				"price_retail",
			case "price_no_added_value":
				cValue, err := util.StringToFloat(record[i])
				if err != nil {
					row.Message = append(row.Message, errors.GetMessage(err))
					row.Ignored = true
				}
				if cValue == 0 {
					row.Ignored = true
				}
				d.PriceNoAddedValue = cValue
			case "price_retail":
				cValue, err := util.StringToFloat(record[i])
				if err != nil {
					row.Message = append(row.Message, errors.GetMessage(err))
					row.Ignored = true
				}
				if cValue == 0 {
					row.Ignored = true
				}
				d.PriceRetail = cValue
			case "featured":
				d.Featured = strings.ToUpper(record[i])
			case "recommand_words":
				d.Featured = strings.ToUpper(record[i])
			case "extra_words":
				d.Featured = strings.ToUpper(record[i])
			case "image", "image1", "image2", "image3", "image4", "image5":
				d.Images = append(d.Images, record[i])
			}
		}

		if row.Ignored {
			unimportRows = append(unimportRows, row)
			continue
		}
		//handle db
		if !row.Ignored {
			if err := d.composeStockRefWithSupplierPrefix(); err != nil {
				//TODO
				return nil, err
			}
			d.diamondImages()

			if err := d.processDiamondRecord(); err != nil {
				//TODO return err for now!
				return nil, err
			}
			//remove it from old stock ref map
			delete(oldStockRefList, d.StockRef)
		}
	}
	util.Println("finish process diamond")
	if err := offlineDiamondsNoLongerExist(oldStockRefList); err != nil {
		return unimportRows, err
	}
	return unimportRows, nil
}

func (d *diamond) processDiamondRecord() error {
	if err := d.processPrice(); err != nil {
		return err
	}
	var id, status string
	var priceNoAddedValue, priceRetail float64
	q := fmt.Sprintf("SELECT id, price_no_added_value, price_retail, status FROM diamonds WHERE stock_ref='%s'", d.StockRef)
	if err := db.QueryRow(q).Scan(&id, &priceNoAddedValue, &priceRetail, &status); err != nil {
		//item not exist in db
		if err == sql.ErrNoRows {
			d.ID = newV4()
			q := d.composeInsertQuery()
			if _, err := dbExec(q); err != nil {
				util.Tracef(`fail to add diamond item. diamond: %s; certificate_number: %s; grading_lab: %s; retail price %f.\n`,
					d.StockRef, d.CertificateNumber, d.GradingLab, d.PriceRetail)
				return err
			}
			util.Tracef(`diamond item added! diamond: %s; certificate_number: %s; grading_lab: %s; retail price %f.\n`,
				d.StockRef, d.CertificateNumber, d.GradingLab, d.PriceRetail)
			return nil
		}
		return err
	}
	//item alread exist in db
	if status != "SOLD" && status != "RESERVED" && (d.PriceRetail-priceRetail) > 5 {
		q := d.composeUpdateQuery()
		if _, err := dbExec(q); err != nil {
			util.Tracef(`retail price changed, but failed to update. diamond: %s; certificate_number: %s; grading_lab: %s; original price: %f; new price should be %f.\n`,
				d.StockRef, d.CertificateNumber, d.GradingLab, priceRetail, d.PriceRetail)
			return err
		}
		util.Tracef(`retail price changed for diamond: %s; certificate_number: %s; grading_lab: %s; original price: %f; new price %f.\n`,
			d.StockRef, d.CertificateNumber, d.GradingLab, priceRetail, d.PriceRetail)
	}

	return nil
}

func (d *diamond) processPrice() error {
	//set the price already, no need to caculate
	if d.PriceRetail != 0 {
		return nil
	}
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

//TODO is stock_ref need compose with prefix???? - change new req handler accordingly
func (d *diamond) composeStockRefWithSupplierPrefix() error {
	q := fmt.Sprintf(`SELECT prefix FROM suppliers WHERE name='%s'`, d.Supplier)
	var prefix string
	if err := dbQueryRow(q).Scan(&prefix); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		prefix = "SYS"
	}
	stockRef := d.StockRef
	d.StockRef = fmt.Sprintf("%s-%s", prefix, stockRef)
	return nil
}

func getAllStockRef() (map[string]struct{}, error) {
	rows, err := dbQuery("SELECT stock_ref FROM diamonds WHERE status!='OFFLINE'")
	if err != nil {
		return nil, err
	}
	stockRefs := make(map[string]struct{})
	for rows.Next() {
		var stockRef string
		if err := rows.Scan(&stockRef); err != nil {
			return nil, err
		}
		//empty struct comsumes 0 bytes
		var s struct{}
		stockRefs[stockRef] = s
	}
	return stockRefs, nil
}

func getAllValidSupplierName() ([]string, error) {
	var suppiers []string
	rows, err := dbQuery("SELECT name FROM suppliers")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		suppiers = append(suppiers, name)
	}
	return suppiers, nil
}

//下线不存在的钻石 //TODO return or just trace err ???
func offlineDiamondsNoLongerExist(stockRefList map[string]struct{}) error {
	util.Tracef("Start to offline all diamonds no longer exists.\n")
	for k := range stockRefList {
		q := fmt.Sprintf("UPDATE diamonds SET status='OFFLINE' WHERE stock_ref ='%s'", k)
		util.Tracef("Offline diamond stock_ref: %s.\n", k)
		if _, err := dbExec(q); err != nil {
			util.Tracef("error when offline diamond. stock_ref: %s. err: .\n", k, errors.GetMessage(err))
			return err
		}
	}
	util.Tracef("Finished offline all diamonds no longer exists.\n")
	return nil
}

func diamondClarity(clarity string) string {
	if len(clarity) != 0 {
		if util.IsInArrayString(strings.ToUpper(clarity), VALID_CLARITY) {
			return strings.ToUpper(clarity)
		}
	}
	return "-"
}

func diamondFluo(fluo string) string {
	if len(fluo) != 0 {
		p := strings.ToUpper(fluo)
		if p == "VERY STRONG" || p == "VST" {
			return "VST"
		}
		if p == "STRONG" || p == "STG" {
			return "STG"
		}
		if p == "SLIGHT" || p == "SLT" || p == "SL" {
			return "SLT"
		}
		if p == "VERY SLIGHT" || p == "VSL" {
			return "VSL"
		}
		if p == "MEDIUM" || p == "MED" || string(p[0]) == "M" {
			return "MED"
		}
		if p == "FAINT" || p == "FNT" || string(p[0]) == "F" {
			return "FNT"
		}
		if p == "NONE" || p == "NON" || string(p[0]) == "N" {
			return "NONE"
		}
	}
	return "UNKOWN-" + strings.ToUpper(fluo)

}

func diamondCutGradeSymmetryPolish(cutGrade string) string {
	if len(cutGrade) != 0 {
		p := strings.ToUpper(cutGrade)
		if p == "EXC" || p == "EXCELLENT" || string(p[0]) == "E" {
			return "EX"
		}
		if p == "VERY GOOD" || string(p[0]) == "V" {
			return "VG"
		}
		if p == "GOOD" || p == "GD" || string(p[0]) == "G" {
			return "G"
		}
		if p == "FAIR" || string(p[0]) == "F" {
			return "F"
		}
	}
	return "UNKOWN-" + strings.ToUpper(cutGrade)
}

//TODO
func diamondColor(color string) string {
	if len(color) != 0 {
		switch strings.ToUpper(color) {
		case "FY", "FANCY YELLOW":
			return "FY"
		case "FLY":
			return "FLY"
		case "FANCY BROWNISH YELLOW", "FBY":
			return "FBY"
		case "FANCY LIGHT BROWNISH YELLOW", "FLBY":
			return "FLBY"
		case "FANCY INTENSE YELLOW", "FIY":
			return "FIY"
		case "FVY", "FANCY VIVID YELLOW":
			return "FVY"
		case "FLBGY":
			return "FLBGY"
		}
	}
	//  D
	//  E
	//  F
	//  G
	//  H
	//  I
	//  J
	//  K
	//  L
	//  M
	//  N
	//  O
	//  P
	//  Q
	//  R
	//  S
	//  T
	//  U
	//  V
	//  W
	//  X
	//  Y
	//  Z
	return "UNKOWN-" + strings.ToUpper(color)
}

func diamondShape(shape string) (string, error) {
	if len(shape) != 0 {
		switch strings.ToUpper(shape) {
		case "BR", "ROUND":
			return "BR", nil
		case "PS", "PEAR":
			return "PS", nil
		case "PR", "PRICESS":
			return "PR", nil
		case "HS", "HEART":
			return "HS", nil
		case "MQ", "MARQUISE":
			return "MQ", nil
		case "OV", "OVAL":
			return "OV", nil
		case "EM", "EMERALD":
			return "EM", nil
		case "CU", "CUSHION":
			return "CU", nil
		case "AS", "ASSCHER":
			return "AS", nil
		case "RAD", "RADIANT", "RA":
			return "RAD", nil
		case "RBC", "RCRB", "RC", "PE", "HT", "CMB":
			return strings.ToUpper(shape), nil
		default:
			return "", errors.Newf("%s is not a valid shape", shape)
		}
	}
	return "", errors.Newf("shape cannot be empty")
}

//TODO should return error - > to add new suppliers
func diamondSupplier(supplier string, suppliers []string) (string, error) {
	if len(supplier) != 0 {
		if util.IsInArrayString(strings.ToUpper(supplier), suppliers) {
			return strings.ToUpper(supplier), nil
		}
	}
	return "", errors.Newf("supplier %s not exist, please add first!", supplier)
}

//TODO should return error ????
func diamondGradingLab(gradingLab string) (string, error) {
	if util.IsInArrayString(strings.ToUpper(gradingLab), VALID_GRADING_LAB) {
		return strings.ToUpper(gradingLab), nil
	}
	return "", errors.Newf("%s is not a valid grading lab", gradingLab)
}

func (d *diamond) diamondImages() {
	var imageNames []string
	for _, imageName := range d.Images {
		name := fmt.Sprintf("beyoudiamond-image-%s-%s", d.StockRef, imageName)
		imageNames = append(imageNames, name)
	}
	d.Images = imageNames
}
