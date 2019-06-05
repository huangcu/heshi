package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"strings"
	"util"
)

// <td width="88">唯一商品号(StockID)</td>
// <td width="88">货号(Name)</td>
// <td width="88">材料</td>
// <td width="88">金重</td>
// <td width="88">是否空托</td>
// <td width="88">最小钻石尺寸</td>
// <td width="88">最大钻石尺寸</td>
// <td width="88">镶碎钻</td>
// <td width="88">小钻数量</td>
// <td width="88">小钻总重</td>
// <td width="88">镶嵌方式</td>
// <td width="88">价格</td>
//TODO must have headers when import
func validateJewelryHeaders(headers []string) []string {
	var missingHeaders []string
	for k, header := range jewelryHeaders {
		if !util.IsInArrayString(header, headers) && k < 8 {
			missingHeaders = append(missingHeaders, header)
		}
	}
	return missingHeaders
}

func importJewelryProducts(uid, file, category string) ([]util.Row, error) {
	oldStockIDList, err := getAllStockIDBySubCategory(category)
	if err != nil {
		return nil, err
	}
	rows, err := util.ParseCSVToStruct(file)
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, errors.New("uploaded file has no rows")
	}

	unimportRows := []util.Row{}
	//get headers
	originalHeaders := rows[0]

	//process rows
	util.Println("start process jewelry")
	for index := 1; index < len(rows); index++ {
		j := jewelry{}
		row := rows[index]
		record := row.Value
		util.Printf("processsing row: %d, %s", index, record)
		for i, header := range originalHeaders.Value {
			switch header {
			case "stock_id":
				if record[i] == "" {
					row.Ignored = true
					row.Message = append(row.Message, "jewelry stock id cannot be empty")
					break
				}
				j.StockID = strings.ToUpper(record[i])
			case "name":
				j.Name = strings.ToUpper(record[i])
			case "need_diamond":
				j.NeedDiamond = strings.ToUpper(record[i])
			case "category":
				j.Category = strings.ToUpper(record[i])
			case "material":
				j.Material = strings.ToUpper(record[i])
			case "dia_shape":
				j.DiaShape = strings.ToUpper(record[i])
			case "metal_weight":
				j.MetalWeightStr = record[i]
			case "mounting_type":
				j.MountingType = strings.ToUpper(record[i])
			case "price":
				j.PriceStr = record[i]
			case "unit_number":
				j.UnitNumber = record[i]
			case "dia_size_min":
				j.DiaSizeMinStr = record[i]
			case "dia_size_max":
				j.DiaSizeMaxStr = record[i]
			case "main_dia_num":
				j.MainDiaNumStr = record[i]
			case "main_dia_size":
				j.MainDiaSizeStr = record[i]
			case "small_dias":
				j.SmallDias = strings.ToUpper(record[i])
			case "small_dia_num":
				j.SmallDiaNumStr = record[i]
			case "small_dia_carat":
				j.SmallDiaCaratStr = record[i]
			case "video_link":
				j.VideoLink = record[i]
			case "text":
				j.Text = record[i]
			case "status":
				j.Status = strings.ToUpper(record[i])
			case "verified":
				j.Verified = strings.ToUpper(record[i])
			case "featured":
				j.Featured = strings.ToUpper(record[i])
			case "stock_quantity":
				j.StockQuantityStr = record[i]
			case "profitable":
				j.Profitable = strings.ToUpper(record[i])
			case "free_acc":
				j.FreeAcc = strings.ToUpper(record[i])
			case "image", "image1", "image2", "image3", "image4", "image5":
				j.Images = append(j.Images, record[i])
			}
		}

		if row.Ignored {
			unimportRows = append(unimportRows, row)
			//move on to next row
			continue
		}

		//validate row
		var id string
		q := fmt.Sprintf("SELECT id FROM jewelrys WHERE stock_id='%s'", j.StockID)
		if err := dbQueryRow(q).Scan(&id); err != nil {
			//new record
			if err == sql.ErrNoRows {
				//validate as new request
				if vemsg, err := j.validateJewelryReq(false); err != nil {
					return nil, err
				} else if len(vemsg) != 0 {
					row.Ignored = true
					for _, v := range vemsg {
						row.Message = append(row.Message, v.Message)
					}
					unimportRows = append(unimportRows, row)
					//move on to next row
					continue
				}
				//pass validation, insert into db
				j.ID = newV4()
				j.jewelryImages()
				q := j.composeInsertQuery()
				if _, err := dbExec(q); err != nil {
					util.Printf("fail to add jewelry item. stock id: %s; err: %s", j.StockID, errors.GetMessage(err))
					return nil, err
				}
				util.Printf("new jewelry item added. stock id: %s", j.StockID)
			} else {
				return nil, err
			}
		}

		//already exist, validate as update request
		if vemsg, err := j.validateJewelryReq(true); err != nil {
			return nil, err
		} else if len(vemsg) != 0 {
			row.Ignored = true
			for _, v := range vemsg {
				row.Message = append(row.Message, v.Message)
			}
			unimportRows = append(unimportRows, row)
			//move on to next row
			continue
		}
		//pass validation, update db
		j.ID = id
		j.jewelryImages()
		q = j.composeUpdateQueryTrack(uid)
		if _, err := dbExec(q); err != nil {
			util.Printf("fail to update jewelry item. stock id: %s; err; %s", j.StockID, errors.GetMessage(err))
			return nil, err
		}
		util.Printf("jewelry item updated. stock id: %s", j.StockID)
		// go newHistoryRecords(uid, "jewelrys", j.ID, j.parmsKV())
		//remove updated stockID from old one as this has been scanned and processed
		delete(oldStockIDList, j.StockID)
	}
	util.Println("finish process jewelry")
	if err := offlineJewelrysNoLongerExist(oldStockIDList); err != nil {
		return unimportRows, err
	}
	return unimportRows, nil
}

func jewelryCategory(category string) (string, error) {
	cate := strings.ToUpper(category)
	if util.IsInArrayString(cate, validCategory) {
		return cate, nil
	}
	return "", errors.Newf("%s not a valid category", category)
}

func jewelryMaterial(material string) string {
	m := strings.ToUpper(strings.Replace(material, " ", "_", -1))
	if util.IsInArrayString(m, validMaterial) {
		return m
	}
	return "UNKNOWN-" + m
}

func jewelryShape(shapeStr string) (string, error) {
	var jShapes []string
	shapes := strings.Split(formatInputString(shapeStr), ",")
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
	if util.IsInArrayString(mountingType, validMountingType) {
		return mt, nil
	}
	return "", errors.Newf("%s is not a valid mounting type", mountingType)
}

func (j *jewelry) jewelryImages() {
	var imageNames []string
	for _, imageName := range j.Images {
		name := fmt.Sprintf("beyoudiamond-image-%s-%s", j.StockID, imageName)
		imageNames = append(imageNames, name)
	}
	j.Images = imageNames
}

// <option value="JP">素金吊坠／项链</option> 1
// <option value="JR">素金戒指</option> 2
// <option value="JE">素金耳环／耳钉</option> 3
// <option value="ZP">镶碎钻吊坠／项链</option> 1 | 5
// <option value="ZR">镶碎钻戒指</option> 2
// <option value="ZE">镶碎钻耳环／耳钉</option> 3
// <option value="CP">成品吊坠／项链</option> 1 | 5 /NO
// <option value="CR">成品戒指</option> 2 /NO
// <option value="CE">成品耳环／耳钉</option> 3/NO
func getAllStockIDBySubCategory(subCategory string) (map[string]struct{}, error) {
	q := ""
	switch strings.ToUpper(subCategory) {
	case "JR":
		q = "small_dias='NO' AND need_diamond='YES' AND category='RING' ORDER BY id ASC"
	case "JE":
		q = "small_dias='NO' AND need_diamond='YES' AND category='EARRING' ORDER BY id ASC"
	case "JP":
		q = "small_dias='NO' AND need_diamond='YES' AND (category='PENDANT' OR category='NECKLACE') ORDER BY id ASC"
	case "ZR":
		q = "small_dias='YES' AND need_diamond='YES' AND category='RING' ORDER BY id ASC"
	case "ZE":
		q = "small_dias='YES' AND need_diamond='YES' AND category='EARRING' ORDER BY id ASC"
	case "ZP":
		q = "small_dias='YES' AND need_diamond='YES' AND (category='PENDANT' OR category='NECKLACE') ORDER BY id ASC"
	case "CR":
		q = "need_diamond='NO' AND category='RING' ORDER BY id ASC"
	case "CE":
		q = "need_diamond='NO' AND category='EARRING' ORDER BY id ASC"
	case "CP":
		q = "need_diamond='NO' AND (category='PENDANT' OR category='NECKLACE') ORDER BY id ASC"
	default:
		return nil, errors.New("missing upload sub category")
	}
	stockIds := make(map[string]struct{})
	q = fmt.Sprintf("SELECT stock_id FROM jewelrys WHERE status IN ('AVAILABLE','OFFLINE') AND %s", q)
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stockID string
		if err := rows.Scan(&stockID); err != nil {
			return nil, err
		}
		var s struct{}
		stockIds[stockID] = s
	}
	return stockIds, nil
}

//下线不存在的钻石 //TODO return or just trace err ???
func offlineJewelrysNoLongerExist(stockIDList map[string]struct{}) error {
	util.Traceln("Start to offline all jewelrys no longer exists.")
	for k := range stockIDList {
		q := fmt.Sprintf("UPDATE jewelrys SET status='OFFLINE',updated_at=(CURRENT_TIMESTAMP) WHERE stock_id ='%s'", k)
		util.Tracef("Offline jewelry stock_id: %s.\n", k)
		if _, err := dbExec(q); err != nil {
			util.Tracef("error when offline jewelry. stock_id: %s. err: \n", k, errors.GetMessage(err))
			return err
		}
	}
	util.Traceln("Finished offline all jewelrys no longer exists")
	return nil
}
