package main

import (
	"fmt"
	"heshi/errors"
	"strconv"
	"strings"
	"time"
	"util"
)

func (j *jewelry) composeInsertQuery() string {
	params := j.parmsKV()
	q := `INSERT INTO jewelrys (id`
	va := fmt.Sprintf(`VALUES ('%s'`, j.ID)
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
		case time.Time:
			va = fmt.Sprintf("%s, '%s'", va, v.(time.Time).Format(timeFormat))
		}
	}
	q = fmt.Sprintf("%s) %s)", q, va)
	return q
}

func (j *jewelry) composeUpdateQuery() string {
	params := j.parmsKV()
	q := `UPDATE jewelrys SET`
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
		case time.Time:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(time.Time).Format(timeFormat))
		}
	}
	q = fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, j.ID)
	return q
}

// 	params := make(map[string]interface{})
//TODO validate input
func (j *jewelry) parmsKV() map[string]interface{} {
	params := make(map[string]interface{})
	if j.StockID != "" {
		params["stock_id"] = j.StockID
	}
	if j.Category != "" {
		params["category"] = j.Category
	}
	if j.MetalWeight != 0 {
		params["metal_weight"] = j.MetalWeight
	}
	if j.Material != "" {
		params["material"] = j.Material
	}
	if j.NeedDiamond != "" {
		params["need_diamond"] = j.NeedDiamond
	}
	if j.Name != "" {
		params["name"] = j.Name
	}
	if j.DiaSizeMin != 0 {
		params["dia_size_min"] = j.DiaSizeMin
	}
	if j.DiaSizeMax != 0 {
		params["dia_size_max"] = j.DiaSizeMax
	}
	if j.MountingType != "" {
		params["mounting_type"] = j.MountingType
	}
	if j.Price != 0 {
		params["price"] = j.Price
	}
	if j.UnitNumber != "" {
		params["unit_number"] = j.UnitNumber
	}
	if j.DiaShape != "" {
		params["dia_shape"] = j.DiaShape
	}
	if j.SmallDias != "" {
		params["small_dias"] = j.SmallDias
	}
	if j.SmallDiaNum != 0 {
		params["small_dia_num"] = j.SmallDiaNum
	}
	if j.SmallDiaCarat != 0 {
		params["small_dia_carat"] = j.SmallDiaCarat
	}
	if j.MainDiaNum != 0 {
		params["main_dia_num"] = j.MainDiaNum
	}
	if j.MainDiaSize != 0 {
		params["main_dia_size"] = j.MainDiaSize
	}
	if j.Featured != "" {
		params["featured"] = j.Featured
	}
	if j.VideoLink != "" {
		params["video_link"] = j.VideoLink
	}
	if len(j.Images) != 0 {
		params["images"] = strings.Join(j.Images, ";")
	}
	if j.Text != "" {
		params["text"] = j.Text
	}
	if j.Status != "" {
		params["status"] = j.Status
	}
	if j.Verified != "" {
		params["verified"] = j.Verified
	}
	if j.Featured != "" {
		params["featured"] = j.Featured
	}
	if j.StockQuantity != 0 {
		params["stock_quantity"] = j.StockQuantity
	}
	if j.Profitable != "" {
		params["profitable"] = j.Profitable
	}
	if j.TotallyScanned != 0 {
		params["totally_scanned"] = j.TotallyScanned
	}
	if j.FreeAcc != "" {
		params["free_acc"] = j.FreeAcc
	}
	return params
}

//if from importCSV, it is otherwize, no need to check duplication of stock_id
func (j *jewelry) validateJewelryReq(update bool) ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage
	if !update && j.MetalWeightStr == "" {
		vemsg = append(vemsg, vemsgMetalWeightEmpty)
	} else if j.MetalWeightStr != "" {
		cValue, err := util.StringToFloat(j.MetalWeightStr)
		if err != nil {
			vemsg = append(vemsg, vemsgMetalWeightNotValid)
		} else if cValue == 0 {
			vemsg = append(vemsg, vemsgMetalWeightNotValid)
		} else {
			j.MetalWeight = cValue
		}
	}

	if j.DiaSizeMinStr != "" {
		pValue, err := util.StringToFloat(j.DiaSizeMinStr)
		if err != nil {
			vemsg = append(vemsg, vemsgDiaSizeMinNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgDiaSizeMinNotValid)
		} else {
			j.DiaSizeMin = pValue
		}
	}

	if j.DiaSizeMaxStr != "" {
		pValue, err := util.StringToFloat(j.DiaSizeMaxStr)
		if err != nil {
			vemsg = append(vemsg, vemsgDiaSizeMaxNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgDiaSizeMaxNotValid)
		} else {
			j.DiaSizeMax = pValue
		}
	}

	if j.MainDiaNumStr != "" {
		pValue, err := strconv.Atoi(j.MainDiaNumStr)
		if err != nil {
			vemsg = append(vemsg, vemsgMainDiaNumNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgMainDiaNumNotValid)
		} else {
			j.MainDiaNum = int64(util.AbsInt(pValue))
		}
	}
	if j.MainDiaSizeStr != "" {
		pValue, err := util.StringToFloat(j.MainDiaSizeStr)
		if err != nil {
			vemsg = append(vemsg, vemsgMainDiaSizeNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgMainDiaSizeNotValid)
		} else {
			j.MainDiaSize = pValue
		}
	}

	if j.SmallDiaNumStr != "" {
		pValue, err := strconv.Atoi(j.SmallDiaNumStr)
		if err != nil {
			vemsg = append(vemsg, vemsgSmallDiaNumNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgSmallDiaNumNotValid)
		} else {
			j.SmallDiaNum = int64(util.AbsInt(pValue))
		}
	}
	if j.SmallDiaCaratStr != "" {
		pValue, err := util.StringToFloat(j.SmallDiaCaratStr)
		if err != nil {
			vemsg = append(vemsg, vemsgSmallDiaCaratNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgSmallDiaCaratNotValid)
		} else {
			j.SmallDiaCarat = pValue
		}
	}
	if j.StockQuantityStr != "" {
		pValue, err := strconv.Atoi(j.StockQuantityStr)
		if err != nil {
			vemsg = append(vemsg, vemsgStockQuantityNotValidJ)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgStockQuantityNotValidJ)
		} else {
			j.StockQuantity = util.AbsInt(pValue)
		}
	}
	if !update && j.PriceStr == "" {
		vemsg = append(vemsg, vemsgPriceEmpty)
	} else if j.PriceStr != "" {
		pValue, err := util.StringToFloat(j.PriceStr)
		if err != nil {
			vemsg = append(vemsg, vemsgPriceNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgPriceNotValid)
		} else {
			j.Price = pValue
		}
	}

	if !update && j.StockID == "" {
		vemsgNotValid.Message = "jewelry stock id can not be empty"
		vemsg = append(vemsg, vemsgNotValid)
	}
	if j.StockID != "" {
		if update {
			if exist, err := isItemExistInDbByPropertyWithDifferentID("jewelrys", "stock_id", j.StockID, j.ID); err != nil {
				return nil, err
			} else if exist {
				vemsgAlreadyExist.Message = "jewelry stock_ref " + j.StockID + " already exists"
				vemsg = append(vemsg, vemsgAlreadyExist)
			}
		} else {
			if exist, err := isItemExistInDbByProperty("jewelrys", "stock_id", j.StockID); err != nil {
				return nil, err
			} else if exist {
				vemsgAlreadyExist.Message = "jewelry stock_ref " + j.StockID + " already exists"
				vemsg = append(vemsg, vemsgAlreadyExist)
			}
		}
	}

	if !update && j.Name == "" {
		vemsgNotValid.Message = "jewelry name can not be empty"
		vemsg = append(vemsg, vemsgNotValid)
	}
	if !update && j.NeedDiamond == "" {
		vemsgNotValid.Message = "please set if jewelry  need diamond or not"
		vemsg = append(vemsg, vemsgNotValid)
	}
	if !update && j.Category == "" {
		vemsgNotValid.Message = "jewelry category can not be empty"
		vemsg = append(vemsg, vemsgNotValid)
	} else if j.Category != "" {
		cate, err := jewelryCategory(j.Category)
		if err != nil {
			return nil, err
		}
		j.Category = cate
	}
	if !update && j.MountingType == "" {
		vemsgNotValid.Message = "jewelry mounting type can not be empty"
		vemsg = append(vemsg, vemsgNotValid)
	} else if j.MountingType != "" {
		mt, err := jewelryMountingType(j.MountingType)
		if err != nil {
			return nil, err
		}
		j.MountingType = mt
	}

	if !update && j.Material == "" {
		vemsgNotValid.Message = "jewelry material can not be empty"
		vemsg = append(vemsg, vemsgNotValid)
	} else if j.Material != "" {
		j.Material = jewelryMaterial(j.Material)
	}

	if !update && j.DiaShape == "" {
		vemsgNotValid.Message = "jewelry diamond shape can not be empty"
		vemsg = append(vemsg, vemsgNotValid)
	} else if j.DiaShape != "" {
		s, err := jewelryShape(j.DiaShape)
		if err != nil {
			return nil, err
		}
		j.DiaShape = s
	}
	//TODO Featured/Online value validate??? - value can only be YES OR NO
	return vemsg, nil
}

//TODO TOBE REMOVED
func (j *jewelry) validateJewelryUpdateReq() ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage
	if j.MetalWeightStr != "" {
		cValue, err := util.StringToFloat(j.MetalWeightStr)
		if err != nil {
			vemsg = append(vemsg, vemsgMetalWeightNotValid)
		} else if cValue == 0 {
			vemsg = append(vemsg, vemsgMetalWeightNotValid)
		} else {
			j.MetalWeight = cValue
		}
	}

	if j.DiaSizeMinStr != "" {
		pValue, err := util.StringToFloat(j.DiaSizeMinStr)
		if err != nil {
			vemsg = append(vemsg, vemsgDiaSizeMinNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgDiaSizeMinNotValid)
		} else {
			j.DiaSizeMin = pValue
		}
	}

	if j.DiaSizeMaxStr != "" {
		pValue, err := util.StringToFloat(j.DiaSizeMaxStr)
		if err != nil {
			vemsg = append(vemsg, vemsgDiaSizeMaxNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgDiaSizeMaxNotValid)
		} else {
			j.DiaSizeMax = pValue
		}
	}

	if j.MainDiaNumStr != "" {
		pValue, err := strconv.Atoi(j.MainDiaNumStr)
		if err != nil {
			vemsg = append(vemsg, vemsgMainDiaNumNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgMainDiaNumNotValid)
		} else {
			j.MainDiaNum = int64(util.AbsInt(pValue))
		}
	}
	if j.MainDiaSizeStr != "" {
		pValue, err := util.StringToFloat(j.MainDiaSizeStr)
		if err != nil {
			vemsg = append(vemsg, vemsgMainDiaSizeNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgMainDiaSizeNotValid)
		} else {
			j.MainDiaSize = pValue
		}
	}

	if j.SmallDiaNumStr != "" {
		pValue, err := strconv.Atoi(j.SmallDiaNumStr)
		if err != nil {
			vemsg = append(vemsg, vemsgSmallDiaNumNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgSmallDiaNumNotValid)
		} else {
			j.SmallDiaNum = int64(util.AbsInt(pValue))
		}
	}
	if j.SmallDiaCaratStr != "" {
		pValue, err := util.StringToFloat(j.SmallDiaCaratStr)
		if err != nil {
			vemsg = append(vemsg, vemsgSmallDiaCaratNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgSmallDiaCaratNotValid)
		} else {
			j.SmallDiaCarat = pValue
		}
	}
	if j.StockQuantityStr != "" {
		pValue, err := strconv.Atoi(j.StockQuantityStr)
		if err != nil {
			vemsg = append(vemsg, vemsgStockQuantityNotValidJ)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgStockQuantityNotValidJ)
		} else {
			j.StockQuantity = util.AbsInt(pValue)
		}
	}
	if j.PriceStr != "" {
		pValue, err := util.StringToFloat(j.PriceStr)
		if err != nil {
			vemsg = append(vemsg, vemsgPriceNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgPriceNotValid)
		} else {
			j.Price = pValue
		}
	}

	if j.StockID != "" {
		if exist, err := isItemExistInDbByPropertyWithDifferentID("jewelrys", "stock_id", j.StockID, j.ID); err != nil {
			return nil, err
		} else if exist {
			vemsgAlreadyExist.Message = "jewelry stock_ref " + j.StockID + " already exists"
			vemsg = append(vemsg, vemsgAlreadyExist)
		}
	}

	if j.Category != "" {
		cate, err := jewelryCategory(j.Category)
		if err != nil {
			return nil, err
		}
		j.Category = cate
	}

	if j.MountingType != "" {
		mt, err := jewelryMountingType(j.MountingType)
		if err != nil {
			return nil, err
		}
		j.MountingType = mt
	}

	if j.Material != "" {
		j.Material = jewelryMaterial(j.Material)
	}

	if j.DiaShape != "" {
		s, err := jewelryShape(j.DiaShape)
		if err != nil {
			return nil, err
		}
		j.DiaShape = s
	}
	//TODO Featured/Online value validate??? - value can only be YES OR NO
	return vemsg, nil
}

func (j *jewelry) isJewelryExistByStockID() error {
	var id string
	if err := dbQueryRow("SELECT id FROM jewelrys WHERE stock_id='?' AND status IN ('AVAILABLE', 'OFFLINE')", j.StockID).Scan(&id); err != nil {
		return err
	}
	j.ID = id
	return nil
}

func isJewelryExistByID(id string) (bool, error) {
	var count int
	if err := dbQueryRow("SELECT COUNT(*) FROM jewelrys WHERE id='?' AND status IN ('AVAILABLE', 'OFFLINE')", id).Scan(&count); err != nil {
		return false, err
	}
	return count == 1, nil
}
