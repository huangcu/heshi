package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"strconv"
	"strings"
	"util"

	uuid "github.com/satori/go.uuid"
)

func validateJewelryHeaders(headers []string) []string {
	var missingHeaders []string
	for k, header := range jewelryHeaders {
		if !util.IsInArrayString(header, headers) && k < 8 {
			missingHeaders = append(missingHeaders, header)
		}
	}
	return missingHeaders
}

func importJewelryProducts(file string) ([][]string, error) {
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
		j := jewelry{}
		record := records[index]
		util.Printf("processsing row: %d, %s", index, record)
		for i, header := range originalHeaders {
			switch header {
			case "stock_id":
				j.StockID = strings.ToUpper(record[i])
			case "name":
				j.Name = record[i]
			case "need_diamond":
				j.NeedDiamond = strings.ToUpper(record[i])
			case "category":
				category, err := jewelryCategory(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				} else {
					j.Category = category
				}
			case "material":
				j.Material = jewelryMaterial(record[i])
			case "dia_shape":
				if s, err := jewelryShape(record[i]); err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				} else {
					j.DiaShape = s
				}
			case "metal_weight":
				cValue, err := util.StringToFloat(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				if cValue == 0 {
					ignored = true
				}
				j.MetalWeight = cValue
			case "mounting_type":
				mt, err := jewelryMountingType(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				} else {
					j.MountingType = mt
				}
			case "price":
				sValue, err := util.StringToFloat(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.Price = sValue

			case "unit_number":
				j.UnitNumber = record[i]
			case "dia_size_min":
				sValue, err := util.StringToFloat(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.DiaSizeMin = sValue
			case "dia_size_max":
				sValue, err := util.StringToFloat(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.DiaSizeMax = sValue
			case "main_dia_num":
				sValue, err := strconv.Atoi(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.MainDiaNum = int64(util.AbsInt(sValue))
			case "main_dia_size":
				sValue, err := util.StringToFloat(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.MainDiaSize = sValue
			case "small_dias":
				j.SmallDias = strings.ToUpper(record[i])
			case "small_dia_num":
				sValue, err := strconv.Atoi(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.SmallDiaNum = int64(util.AbsInt(sValue))
			case "small_dia_carat":
				sValue, err := util.StringToFloat(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.SmallDiaCarat = sValue
			case "video_link":
				j.VideoLink = record[i]
			case "text":
				j.Text = record[i]
			case "online":
				j.Online = strings.ToUpper(record[i])
			case "verified":
				j.Verified = strings.ToUpper(record[i])
			case "in_stock":
				j.InStock = strings.ToUpper(record[i])
			case "featured":
				j.Featured = strings.ToUpper(record[i])
			case "stock_quantity":
				sValue, err := strconv.Atoi(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.StockQuantity = util.AbsInt(sValue)
			case "profitable":
				j.Profitable = strings.ToUpper(record[i])
			case "free_acc":
				j.FreeAcc = strings.ToUpper(record[i])
			}
		}
		//insert into db
		if !ignored {
			var id string
			if err := dbQueryRow(fmt.Sprintf("SELECT id FROM jewelrys WHERE stock_id='%s'", j.StockID)).Scan(&id); err != nil {
				if err == sql.ErrNoRows {
					j.ID = uuid.NewV4().String()
					q := j.composeInsertQuery()
					if _, err := dbExec(q); err != nil {
						return nil, err
						// ignoredRows = append(ignoredRows, record)
					}
				} else {
					// ignoredRows = append(ignoredRows, record)
					return nil, err
				}
			} else {
				j.ID = id
				q := j.composeUpdateQuery()
				if _, err := dbExec(q); err != nil {
					// ignoredRows = append(ignoredRows, record)
					return nil, err
				}
			}
		}
	}
	util.Println("finish process jewelry")
	return ignoredRows, nil
}

func jewelryCategory(category string) (string, error) {
	cate := strings.ToUpper(category)
	if util.IsInArrayString(cate, VALID_CATEGORY) {
		return cate, nil
	}
	return "", errors.Newf("%s not a valid category", category)
}

func jewelryMaterial(material string) string {
	m := strings.ToUpper(strings.Replace(material, " ", "_", -1))
	if util.IsInArrayString(m, VALID_MATERIAL) {
		return m
	}
	return "UNKNOW-" + m
}

func jewelryShape(shapeStr string) (string, error) {
	var jShapes []string
	shapes := strings.Split(shapeStr, ",")
	for _, shape := range shapes {
		s, err := diamondShape(shape)
		if err != nil {
			return "", err
		}
		jShapes = append(jShapes, s)
	}
	return strings.Join(jShapes, ","), nil
}

func jewelryMountingType(mountingType string) (string, error) {
	mt := strings.ToUpper(mountingType)
	if util.IsInArrayString(mountingType, VALID_MOUNTING_TYPE) {
		return mt, nil
	}
	return "", errors.Newf("%s is not a valid mounting type", mountingType)
}
