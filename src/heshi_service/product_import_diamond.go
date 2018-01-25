package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
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
				d.DiamondID = strings.ToUpper(record[i])
			case "stock_ref":
				d.StockRef = strings.ToUpper(record[i])
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
				d.Color = diamondColor(record[i])
			case "clarity":
				d.Clarity = diamondClarity(record[i])
			case "grading_lab":
				d.GradingLab = diamondGradingLab(record[i])
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
			}
		}
		//handle db
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
	util.Println("finish process diamond")
	if err := offlineDiamondsNoLongerExist(oldStockRefList); err != nil {
		return ignoredRows, err
	}
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
	util.Tracef("Start to offline all diamonds no longer exists")
	for k := range stockRefList {
		q := fmt.Sprintf("UPDATE diamonds SET status='OFFLINE' WHERE stock_ref ='%s'", k)
		if _, err := dbExec(q); err != nil {
			util.Tracef("error when offline diamond. stock_ref: %s. err: ", k, errors.GetMessage(err))
			return err
		}
	}
	util.Tracef("Finished offline all diamonds no longer exists")
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
	return strings.ToUpper(color)
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
	// return "UNKOWN-" + strings.ToUpper(color)
}

func diamondShape(shape string) string {
	if len(shape) != 0 {
		switch strings.ToUpper(shape) {
		case "BR", "ROUND":
			return "BR"
		case "PS", "PEAR":
			return "PS"
		case "PR", "PRICESS":
			return "PR"
		case "HS", "HEART":
			return "HS"
		case "MQ", "MARQUISE":
			return "MQ"
		case "OV", "OVAL":
			return "OV"
		case "EM", "EMERALD":
			return "EM"
		case "CU", "CUSHION":
			return "CU"
		case "AS", "ASSCHER":
			return "AS"
		case "RAD", "RADIANT", "RA":
			return "RAD"
		case "RBC", "RCRB", "RC", "PE", "HT", "CMB":
			return strings.ToUpper(shape)
		}
	}
	return "-"
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
func diamondGradingLab(gradingLab string) string {
	if util.IsInArrayString(strings.ToUpper(gradingLab), VALID_GRADING_LAB) {
		return strings.ToUpper(gradingLab)
	}
	return "UNKOWN-" + strings.ToUpper(gradingLab)
}
