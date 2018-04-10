package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"
	"strings"
	"util"

	"github.com/gin-gonic/gin"
)

func customizeProduct(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "diasize":
		customizeDiamondSize(c)
	case "jewelrycategory":
		customizeJewelryCategory(c)
	case "jewelryitem":
		customizeJewelryItems(c)
	case "diaquality":
		customizeJewelryDiamondsQualityFirst(c)
	case "diamaxcarat":
		customizeJewelryDiamondsMaxCarat(c)
	default:
		c.AbortWithStatus(404)
		return
	}
}

// EuroToDollar PRICE
func customizeDiamondSize(c *gin.Context) {
	budget, err := util.StringToFloat(c.PostForm("budget"))
	if err != nil {
		c.JSON(http.StatusBadGateway, fmt.Sprintf("Input budget:%s is not valid", c.PostForm("budget")))
		return
	}
	q := fmt.Sprintf(`SELECT carat 
		FROM diamonds 
		WHERE shape='BR' 
		AND price_retail <= '%f' 
		AND status IN ('AVAILABLE', 'OFFLINE') 
		 ORDER BY carat DESC LIMIT 1`, budget-60)
	var carat float64
	if err := dbQueryRow(q).Scan(&carat); err != nil {
		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, "NO-DIAMOND-FOUND")
		return
	}
	c.JSON(http.StatusOK, carat)
}

func customizeJewelryCategory(c *gin.Context) {
	// 	$budget=$_POST['budget'];
	budget, err := util.StringToFloat(c.PostForm("budget"))
	if err != nil {
		vemsgNotValid.Message = fmt.Sprintf("Input budget:%s is not valid", c.PostForm("budget"))
		c.JSON(http.StatusOK, vemsgNotValid)
		return
	}
	diaSize, err := util.StringToFloat(c.PostForm("dia_size"))
	if err != nil {
		vemsgNotValid.Message = fmt.Sprintf("Input diamond size:%s is not valid", c.PostForm("dia_size"))
		c.JSON(http.StatusOK, vemsgNotValid)
		return
	}
	q := fmt.Sprintf(`SELECT price_retail 
		FROM diamonds 
		WHERE shape='BR' 
		AND carat >= '%f' 
		AND status IN ('AVAILABLE', 'OFFLINE') 
		 ORDER BY price_retail DESC LIMIT 1`, diaSize)

	var diamondRetailPrice float64

	if err := dbQueryRow(q).Scan(&diamondRetailPrice); err != nil {
		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, "no diamond found")
		return
	}
	//DollarToEuro
	budgetJewelry := budget - diamondRetailPrice
	q = fmt.Sprintf(`SELECT DISTINCT category 
	FROM jewelrys 
	WHERE need_diamond='YES' 
	AND status IN ('AVAILABLE', 'OFFLINE') 
	AND price <= '%f' 
	AND dia_size_min < '%f' 
	AND dia_size_max > '%f'`, budgetJewelry, diaSize, diaSize)

	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	var categorys []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		categorys = append(categorys, category)
	}
	c.JSON(http.StatusOK, categorys)
}

func customizeJewelryItems(c *gin.Context) {
	budget, err := util.StringToFloat(c.PostForm("budget"))
	if err != nil {
		vemsgNotValid.Message = fmt.Sprintf("Input budget:%s is not valid", c.PostForm("budget"))
		c.JSON(http.StatusOK, vemsgNotValid)
		return
	}
	diaSize, err := util.StringToFloat(c.PostForm("dia_size"))
	if err != nil {
		vemsgNotValid.Message = fmt.Sprintf("Input diamond size:%s is not valid", c.PostForm("dia_size"))
		c.JSON(http.StatusOK, vemsgNotValid)
		return
	}
	category := c.PostForm("category")
	if !util.IsInArrayString(category, VALID_CATEGORY) {
		c.JSON(http.StatusOK, vemsgUploadProductsCategoryNotValid)
		return
	}

	q := fmt.Sprintf(`SELECT price_retail 
		FROM diamonds 
		WHERE shape='BR' 
		AND carat >= '%f' 
		AND status IN ('AVAILABLE', 'OFFLINE') 
		 ORDER BY price_retail DESC LIMIT 1`, diaSize)

	var diamondRetailPrice float64

	if err := dbQueryRow(q).Scan(&diamondRetailPrice); err != nil {
		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, "no diamond found")
		return
	}
	//DollarToEuro
	budgetJewelry := budget - diamondRetailPrice
	q = fmt.Sprintf(`SELECT jewelrys.id, stock_id, category, unit_number, dia_shape, material, metal_weight, 
		need_diamond, name, dia_size_min, dia_size_max, small_dias, small_dia_num, small_dia_carat, 
	 mounting_type, main_dia_num, main_dia_size, video_link, text, jewelrys.status, verified, 
	 featured, price, stock_quantity, profitable, totally_scanned, free_acc, last_scan_at,offline_at,
	 promotions.id, prom_type, prom_discount, prom_price, begin_at, end_at, promotions.status 
	FROM jewelrys 
	LEFT JOIN promotions ON jewelrys.promotion_id=promotions.id 
	WHERE need_diamond='YES' 
	AND jewelrys.status IN ('AVAILABLE', 'OFFLINE') 
	AND category = '%s'
	AND price <= '%f' 
	AND dia_size_min < '%f' 
	AND dia_size_max > '%f'`, category, budgetJewelry, diaSize, diaSize)

	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()
	//TODO price - > agent price, customer price base on level
	// 	require_once('../../_includes/functions/currency_calculator.php');
	// require_once('../../_includes/functions/accountprice.php');
	// require_once('../../_includes/functions/agentprice.php');
	// 	if(isset($agent) && $agent=='YES'){
	// 	$jew_price_f = round(priceforagent_jewelry($agent_level, $row['price']), 2);
	// }else{
	// 	$jew_price_f = round(priceforaccount_jewelry($accountlevel, $row['price']), 2);
	// }
	js, err := composeJewelry(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, js)
}

func customizeJewelryDiamondsQualityFirst(c *gin.Context) {
	budget, err := util.StringToFloat(c.PostForm("budget"))
	if err != nil {
		vemsgNotValid.Message = fmt.Sprintf("Input budget:%s is not valid", c.PostForm("budget"))
		c.JSON(http.StatusOK, vemsgNotValid)
		return
	}
	diaSize, err := util.StringToFloat(c.PostForm("dia_size"))
	if err != nil {
		vemsgNotValid.Message = fmt.Sprintf("Input diamond size:%s is not valid", c.PostForm("dia_size"))
		c.JSON(http.StatusOK, vemsgNotValid)
		return
	}
	jewelryID := c.PostForm("jewelry_id")

	q := fmt.Sprintf(`SELECT category, dia_shape, dia_size_min, dia_size_max, price 
		FROM jewelrys WHERE id = '%s'`, jewelryID)

	var category, diaShape string
	var diaSizeMin, diaSizeMax, price float64
	if err := dbQueryRow(q).Scan(&category, &diaShape, &diaSizeMin, &diaSizeMax, &price); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	var ss []string
	for _, v := range strings.Split(diaShape, ",") {
		s := fmt.Sprintf("shape='%s'", v)
		ss = append(ss, s)
	}
	if (diaSize + 0.3) < diaSizeMax {
		diaSizeMax = diaSize + 0.3
	}
	if (diaSize - 0.3) < diaSizeMin {
		diaSizeMin = diaSize - 0.3
	}

	clarityFields := strings.Join(VALID_CLARITY, "','")
	//DollarToEuro round(EuroToDollar($budget-priceforaccount_jewelry($accountlevel, $jew_price))*1.02)
	budgetDiamond := budget - price
	if category != "EARRING" {
		q := fmt.Sprintf(`SELECT * FROM diamonds
		WHERE (%s)
		AND carat>='%f'
		AND carat<='%f'
		AND (color = 'F' OR color = 'E' OR color = 'D' OR color = 'G') 
		AND (clarity = 'VVS1' OR clarity = 'FL' OR clarity = 'IF' OR clarity = 'VVS2') 
		AND cut_grade = 'EX' 
		AND polish = 'EX' 
		AND symmetry = 'EX' 
		AND fluorescence_intensity = 'NONE' 
		AND (grading_lab = 'HRD' OR grading_lab = 'IGI' OR grading_lab='GIA') 
		AND price_retail <= '%f' 
		AND status IN ('AVAILABLE', 'OFFLINE') 
		ORDER BY color ASC, Field(clarity, %s) ASC, price_retail DESC LIMIT 2`,
			strings.Join(ss, "OR"), diaSizeMin, diaSizeMax, budgetDiamond, clarityFields)
		fmt.Println(q)
	} else {
		q := fmt.Sprintf(`SELECT COUNT(*) AS howmanySamesize, carat FROM diamonds 
		WHERE (%s) 
		AND carat>='%f'
		AND carat<='%f'
		AND (color = "F" OR color = "E" OR color = "D" OR color = "G" OR color = "H") 
		AND (clarity = "VVS1" OR clarity = "FL" OR clarity = "IF" OR clarity = "VVS2" OR clarity = "VS1") 
		AND (cut_grade = "EX" OR cut_grade = "VG") 
		AND (polish = "EX" OR polish = "VG") 
		AND (symmetry = "EX" OR symmetry = "VG") 
		AND fluorescence_intensity = "NONE" 
		AND (grading_lab = "HRD" OR grading_lab = "IGI" OR grading_lab="GIA") 
		AND price_retail <= '%f' 
		AND status IN ('AVAILABLE', 'OFFLINE') 
		GROUP BY carat 
		ORDER BY color ASC, Field(clarity, %s) ASC, price_retail DESC LIMIT 200`,
			strings.Join(ss, "OR"), diaSizeMin, diaSizeMax, budgetDiamond, clarityFields)
		fmt.Println(q)
		rows, err := dbQuery(q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		defer rows.Close()

		caratSameNumber := make(map[float64]int)
		for rows.Next() {
			var caratSize float64
			var number int
			if err := rows.Scan(&number, &caratSize); err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			caratSameNumber[caratSize] = number
		}
		var dds []diamond
		for caratSize, number := range caratSameNumber {
			if number < 1 {
				continue
			}

			q = fmt.Sprintf(`SELECT * FROM diamonds
				WHERE (%s)
				AND carat>='%f'
				AND carat<='%f'
				AND (color = 'F' OR color = 'E' OR color = 'D' OR color = 'G' OR color = 'H') 
				AND (clarity = 'VVS1' OR clarity = 'FL' OR clarity = 'IF' OR clarity = 'VVS2') 
				AND (cut_grade = 'EX' OR cut_grade = 'VG')
				AND (polish = 'EX' OR polish = 'VG')
				AND (symmetry = 'EX' OR symmetry = 'VG')
				AND fluorescence_intensity = 'NONE' 
				AND (grading_lab = 'HRD' OR grading_lab = 'IGI' OR grading_lab='GIA') 
				AND price_retail <= '%f' 
				AND status IN ('AVAILABLE', 'OFFLINE') 
				ORDER BY color ASC, Field(clarity, %s) ASC, price_retail DESC LIMIT 2`,
				strings.Join(ss, "OR"), caratSize-0.01, caratSize+0.01, budgetDiamond, clarityFields)

			rows, err := dbQuery(q)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			defer rows.Close()

			ds, err := composeDiamond(rows)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			dds = append(dds, ds...)
		}
		c.JSON(http.StatusOK, dds)
	}
}

func customizeJewelryDiamondsMaxCarat(c *gin.Context) {
	budget, err := util.StringToFloat(c.PostForm("budget"))
	if err != nil {
		vemsgNotValid.Message = fmt.Sprintf("Input budget:%s is not valid", c.PostForm("budget"))
		c.JSON(http.StatusOK, vemsgNotValid)
		return
	}
	diaSize, err := util.StringToFloat(c.PostForm("dia_size"))
	if err != nil {
		vemsgNotValid.Message = fmt.Sprintf("Input diamond size:%s is not valid", c.PostForm("dia_size"))
		c.JSON(http.StatusOK, vemsgNotValid)
		return
	}
	jewelryID := c.PostForm("jewelry_id")

	q := fmt.Sprintf(`SELECT category, dia_shape, dia_size_min, dia_size_max, price 
		FROM jewelrys WHERE id = '%s'`, jewelryID)

	var category, diaShape string
	var diaSizeMin, diaSizeMax, price float64
	if err := dbQueryRow(q).Scan(&category, &diaShape, &diaSizeMin, &diaSizeMax, &price); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	var ss []string
	for _, v := range strings.Split(diaShape, ",") {
		s := fmt.Sprintf("shape='%s'", v)
		ss = append(ss, s)
	}
	if (diaSize + 0.3) < diaSizeMax {
		diaSizeMax = diaSize + 0.3
	}
	if (diaSize - 0.3) < diaSizeMin {
		diaSizeMin = diaSize - 0.3
	}

	//DollarToEuro round(EuroToDollar($budget-priceforaccount_jewelry($accountlevel, $jew_price))*1.02)
	budgetDiamond := budget - price
	if category != "EARRING" {
		q := fmt.Sprintf(`SELECT * FROM diamonds
		WHERE (%s)
		AND carat>='%f'
		AND carat<='%f'
	  AND (color='I' OR color='J' OR color='K' OR color='L' OR color='M' OR color='N') 
	  AND (clarity = 'VS1' OR clarity = 'SI1' OR clarity = 'SI2' OR clarity = 'VS2') 
		AND cut_grade = 'EX' 
		AND polish = 'EX' 
		AND symmetry = 'EX' 
		AND fluorescence_intensity = 'NONE' 
		AND (grading_lab = 'HRD' OR grading_lab = 'IGI' OR grading_lab='GIA') 
		AND price_retail <= '%f' 
		AND status IN ('AVAILABLE', 'OFFLINE') 
		ORDER BY carat DESC LIMIT 2`,
			strings.Join(ss, "OR"), diaSizeMin, diaSizeMax, budgetDiamond)
		fmt.Println(q)
	} else {
		q := fmt.Sprintf(`SELECT COUNT(*) AS howmanySamesize, carat FROM diamonds 
		WHERE (%s) 
		AND carat>='%f'
		AND carat<='%f'
	  AND (color='I' OR color='J' OR color='K' OR color='L' OR color='M' OR color='N') 
	  AND (clarity = 'VS1' OR clarity = 'SI1' OR clarity = 'SI2' OR clarity = 'VS2') 
		AND (cut_grade = "EX" OR cut_grade = "VG") 
		AND (polish = "EX" OR polish = "VG") 
		AND (symmetry = "EX" OR symmetry = "VG") 
		AND fluorescence_intensity = "NONE" 
		AND (grading_lab = "HRD" OR grading_lab = "IGI" OR grading_lab="GIA") 
		AND price_retail <= '%f' 
		AND status IN ('AVAILABLE', 'OFFLINE') 
		GROUP BY carat 
		ORDER BY carat DESC LIMIT 200`,
			strings.Join(ss, "OR"), diaSizeMin, diaSizeMax, budgetDiamond)
		fmt.Println(q)
		rows, err := dbQuery(q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		defer rows.Close()

		caratSameNumber := make(map[float64]int)
		for rows.Next() {
			var caratSize float64
			var number int
			if err := rows.Scan(&number, &caratSize); err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			caratSameNumber[caratSize] = number
		}
		var dds []diamond
		for caratSize, number := range caratSameNumber {
			if number < 1 {
				continue
			}

			q = fmt.Sprintf(`SELECT * FROM diamonds
				WHERE (%s)
				AND carat>='%f'
				AND carat<='%f'
	  		AND (color='I' OR color='J' OR color='K' OR color='L' OR color='M' OR color='N') 
	  		AND (clarity = 'VS1' OR clarity = 'SI1' OR clarity = 'SI2' OR clarity = 'VS2') 
				AND (cut_grade = 'EX' OR cut_grade = 'VG')
				AND (polish = 'EX' OR polish = 'VG')
				AND (symmetry = 'EX' OR symmetry = 'VG')
				AND fluorescence_intensity = 'NONE' 
				AND (grading_lab = 'HRD' OR grading_lab = 'IGI' OR grading_lab='GIA') 
				AND price_retail <= '%f' 
				AND status IN ('AVAILABLE', 'OFFLINE') 
				ORDER BY price_retail ASC LIMIT 2`,
				strings.Join(ss, "OR"), caratSize-0.01, caratSize+0.01, budgetDiamond)

			rows, err := dbQuery(q)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			defer rows.Close()
			ds, err := composeDiamond(rows)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
				return
			}
			dds = append(dds, ds...)
		}
		c.JSON(http.StatusOK, dds)
	}
}
