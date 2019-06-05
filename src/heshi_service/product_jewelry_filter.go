package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"util"

	"github.com/gin-gonic/gin"
)

// ref - search by name and ID
func searchJewelrys(c *gin.Context) ([]jewelry, error) {
	q := fmt.Sprintf(`SELECT jewelrys.id, stock_id, category, unit_number, dia_shape, material, metal_weight, 
		need_diamond, name, dia_size_min, dia_size_max, small_dias, small_dia_num, small_dia_carat, 
	 mounting_type, main_dia_num, main_dia_size, video_link, images, text, jewelrys.status, verified, 
	 featured, price, stock_quantity, profitable, totally_scanned, free_acc, last_scan_at,offline_at,
	 promotions.id, prom_type, prom_discount, prom_price, begin_at, end_at, promotions.status 
	 FROM jewelrys 
	 LEFT JOIN promotions ON jewelrys.promotion_id=promotions.id 
	 WHERE stock_id = '%s' 
	 AND jewelrys.status in ('AVAILABLE','OFFLINE') 
	 AND stock_quantity > 0`, strings.ToUpper(c.PostForm("ref")))

	class := strings.ToUpper(c.Query("class"))
	needDiamond := ""
	if class == "NOMOUNTING" {
		needDiamond = "NO"
	}
	if class == "MOUNTING" {
		needDiamond = "YES"
	}
	if needDiamond != "" {
		q = fmt.Sprintf("%s AND need_diamond= '%s'", q, needDiamond)
	}
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	js, err := composeJewelry(rows)
	if err != nil {
		return nil, err
	}
	return js, nil
}

func filterJewelrys(c *gin.Context) ([]jewelry, error) {
	q, err := composeFilterJewelryQuery(c)
	if err != nil {
		return nil, err
	}
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	js, err := composeJewelry(rows)
	if err != nil {
		return nil, err
	}
	return js, nil
}

// class:mounting or complete
// category:2
// size:
// material:
// price:
// mounting_type:
// sds:
// diashape:
// crrpage:
func composeFilterJewelryQuery(c *gin.Context) (string, error) {
	var querys []string
	class := strings.ToUpper(c.Query("class"))
	needDiamond := ""
	if class == "NOMOUNTING" {
		needDiamond = "NO"
	}
	if class == "MOUNTING" {
		needDiamond = "YES"
	}
	if needDiamond != "" {
		querys = append(querys, fmt.Sprintf("need_diamond='%s'", needDiamond))
	}

	if c.PostForm("category") != "" {
		querys = append(querys, fmt.Sprintf("category='%s'", strings.ToUpper(c.PostForm("category"))))
	}
	if c.PostForm("material") != "" {
		querys = append(querys, fmt.Sprintf("material='%s'", strings.ToUpper(c.PostForm("material"))))
	}
	if c.PostForm("size") != "" {
		cValue, err := strconv.ParseFloat(c.PostForm("size"), 64)
		if err != nil {
			return "", err
		}
		querys = append(querys, fmt.Sprintf("dia_size_min<=%f", math.Abs(cValue)))
		querys = append(querys, fmt.Sprintf("dia_size_max>='%f'", math.Abs(cValue)))
	}
	if c.PostForm("price") != "" {
		price, err := strconv.Atoi(c.PostForm("price"))
		if err != nil {
			return "", err
		}
		if price == 300 {
			querys = append(querys, fmt.Sprintf("price<=%d", price))
		} else if price == 1500 {
			querys = append(querys, fmt.Sprintf("price>=%d", price))
		} else {
			maxPrice := price + 300
			querys = append(querys, fmt.Sprintf("price<=%d AND price>=%d", maxPrice, price))
		}
	}
	sds := c.PostForm("small_dias")
	if sds != "" {
		querys = append(querys, fmt.Sprintf("small_dias='%s'", sds))
	}

	if c.PostForm("mounting_type") != "" {
		querys = append(querys, fmt.Sprintf("mounting_type='%s'", strings.ToUpper(c.PostForm("mounting_type"))))
	}

	if c.PostForm("dia_shape") != "" {
		querys = append(querys, fmt.Sprintf("dia_shape LIKE '%s'", strings.ToUpper(c.PostForm("dia_shape"))))
	}

	if c.PostForm("stock_quantity") != "" {
		querys = append(querys, fmt.Sprintf("stock_quantity > '%s'", strings.ToUpper(c.PostForm("stock_quantity"))))
	} else {
		querys = append(querys, "stock_quantity > 0")
	}

	if c.PostForm("status") != "" {
		querys = append(querys, fmt.Sprintf("jewelrys.status='%s'", strings.ToUpper(c.PostForm("status"))))
	} else {
		querys = append(querys, "jewelrys.status in ('AVAILABLE','OFFLINE')")
	}

	var limit string
	currentPage := 1
	if c.PostForm("crr_page") != "" {
		var err error
		currentPage, err = strconv.Atoi(c.PostForm("crr_page"))
		if err != nil {
			return "", err
		}
		//32 records per page
		limit = fmt.Sprintf("LIMIT 32 OFFSET %d", util.AbsInt(currentPage-1)*32)
	}

	// TODO name, shouldn't group by name, name can be same???
	// WHERE (%s) GROUP BY name
	//  ORDER BY jewelrys.status DESC, stock_quantity DESC, jewelrys.created_at DESC %s`,
	q := fmt.Sprintf(`SELECT jewelrys.id, stock_id, category, unit_number, dia_shape, material, metal_weight, 
		need_diamond, name, dia_size_min, dia_size_max, small_dias, small_dia_num, small_dia_carat, 
	 mounting_type, main_dia_num, main_dia_size, video_link, images, text, jewelrys.status, verified, 
	 featured, price, stock_quantity, profitable, totally_scanned, free_acc, last_scan_at,offline_at,
	 promotions.id, prom_type, prom_discount, prom_price, begin_at, end_at, promotions.status 
	 FROM jewelrys 
	 LEFT JOIN promotions ON jewelrys.promotion_id=promotions.id
	 WHERE (%s) 
	 ORDER BY jewelrys.status DESC, stock_quantity DESC, jewelrys.created_at DESC %s`,
		strings.Join(querys, ") AND ("), limit)
	return q, nil
}
