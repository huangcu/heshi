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
	// 	if(isset($_GET['ref'])){
	// 	$searchkey_raw=strtoupper($_GET['ref']);
	// 	$toreplace=array('JR','JE','JP','ZR','ZE','ZP','CR','CE','CP');
	// 	$searchkey_p1=str_replace($toreplace, '', $searchkey_raw);
	// 	$searchkey=trim(str_replace('-', '', $searchkey_p1));
	//   $sql_num='SELECT COUNT(*) AS num FROM jewelry WHERE need_diamond = "'.$need_diamond.'" AND id = ? AND online = "YES"';
	// 	$stmt_num=$conn->prepare($sql_num);
	// 	$stmt_num->execute(array($searchkey));
	// }else{
	// 	$sql_num='SELECT COUNT(DISTINCT name) AS num FROM jewelry WHERE need_diamond = "'.$need_diamond.'" '.$query_category.$query_size.$query_material.$query_price.$query_mountingtype.$query_sds.$query_diashape.'';
	// 	$stmt_num=$conn->query($sql_num);
	// }
	// if(isset($_GET['ref'])){
	//   $sql='SELECT * FROM jewelry WHERE need_diamond = "'.$need_diamond.'" AND id = ? AND online = "YES"';
	// 	$stmt=$conn->prepare($sql);
	// 	$stmt->execute(array($searchkey));
	// }else{
	// 	$sql='SELECT * FROM jewelry WHERE need_diamond = "'.$need_diamond.'" '.$query_category.$query_size.$query_material.$query_price.$query_mountingtype.$query_sds.$query_diashape.' GROUP BY name ORDER BY online DESC, stock_quantity DESC, created_at DESC LIMIT '.$startFrom.',32';
	// 	$stmt=$conn->query($sql);
	// }
	q := fmt.Sprintf(`SELECT id, stock_id, category, unit_number, dia_shape, material, metal_weight, need_diamond, name, 
	 dia_size_min, dia_size_max, small_dias, small_dia_num, small_dia_carat, mounting_type, main_dia_num, main_dia_size, 
	 video_link, images, text, online, verified, in_stock, featured, price, stock_quantity, profitable,
	 totally_scanned, free_acc, last_scan_at,offline_at
	 FROM jewelrys WHERE stock_id = '%s' AND online = 'YES'`, strings.ToUpper(c.PostForm("ref")))

	class := strings.ToUpper(c.Query("class"))
	needDiamond := ""
	if class == "NOMOUNTING" {
		needDiamond = "NO"
	}
	if class == "MOUNTING" {
		needDiamond = "YES"
	}
	if needDiamond != "" {
		q = fmt.Sprintf("%s need_diamond= '%s'", q, needDiamond)
	}
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
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
	class := c.Query("class")
	var querys []string
	needDiamond := "NO"
	if class == "mounting" {
		needDiamond = "YES"
	}
	querys = append(querys, fmt.Sprintf("need_diamond='%s'", needDiamond))
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

		smallDias := "NO"
		if c.PostForm("sds") == "YES" {
			smallDias = "YES"
			querys = append(querys, fmt.Sprintf("small_dias='%s'", smallDias))
		}
	}

	if c.PostForm("mounting_type") != "" {
		querys = append(querys, fmt.Sprintf("mounting_type='%s'", strings.ToUpper(c.PostForm("mounting_type"))))
	}

	if c.PostForm("diashape") != "" {
		querys = append(querys, fmt.Sprintf("dia_shape LIKE '%s'", strings.ToUpper(c.PostForm("diashape"))))
	}

	if c.PostForm("stock_quantity") != "" {
		querys = append(querys, fmt.Sprintf("stock_quantity > '%s'", strings.ToUpper(c.PostForm("stock_quantity"))))
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
	q := fmt.Sprintf(`SELECT id, stock_id, category, unit_number, dia_shape, material, metal_weight, need_diamond, name, 
	 dia_size_min, dia_size_max, small_dias, small_dia_num, small_dia_carat, mounting_type, main_dia_num, main_dia_size, 
	 video_link, images, text, online, verified, in_stock, featured, price, stock_quantity, profitable,
	 totally_scanned, free_acc, last_scan_at,offline_at
	 FROM jewelrys WHERE (%s) GROUP BY name ORDER BY online DESC, stock_quantity DESC, created_at DESC %s`,
		strings.Join(querys, ") AND ("), limit)
	util.Println(q)
	return q, nil
}
