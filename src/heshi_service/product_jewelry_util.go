package main

import (
	"fmt"
	"strings"
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
		}
	}
	q = fmt.Sprintf("%s WHERE id='%s'", strings.TrimSuffix(q, ","), j.ID)
	return q
}

// 	params := make(map[string]interface{})
//TODO validate input
func (j *jewelry) parmsKV() map[string]interface{} {
	params := make(map[string]interface{})
	if j.StockID != "" {
		params["stock_id"] = j.StockID
	}
	if j.Category != 0 {
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
	if j.NameSuffix != 0 {
		params["name_suffix"] = j.NameSuffix
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
	if j.Text != "" {
		params["text"] = j.Text
	}
	if j.Online != "" {
		params["online"] = j.Online
	}
	if j.Verified != "" {
		params["verified"] = j.Verified
	}
	if j.InStock != "" {
		params["in_stock"] = j.InStock
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
