package main

import (
	"fmt"
	"heshi/errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"util"

	"github.com/gin-gonic/gin"
)

//TODO search
func searchProducts(c *gin.Context) {
	category := c.Param("category")
	if !util.IsInArrayString(category, []string{"diamonds", "jewelrys", "gems"}) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if category == "diamonds" {
		ds, err := searchDiamonds(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, ds)
		return
	}
	if category == "jewelrys" {
		js, err := searchJewelrys(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, js)
		return
	}
	if category == "gems" {
		gs, err := searchGems(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, gs)
		return
	}
}

func filterProducts(c *gin.Context) {
	category := c.Param("category")
	if !util.IsInArrayString(category, []string{"diamonds", "jewelrys", "gems"}) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if category == "diamonds" {
		ds, err := filterDiamonds(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, ds)
		return
	}
	if category == "jewelrys" {
		js, err := filterJewelrys(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, js)
		return
	}
	if category == "gems" {
		gs, err := filterGems(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, gs)
		return
	}
}

func searchDiamonds(c *gin.Context) ([]diamond, error) {
	ref := strings.ToUpper(c.PostForm("ref"))
	q := fmt.Sprintf(`SELECT id, diamond_id, stock_ref, shape, carat, color, clarity, grading_lab, 
		certificate_number, cut_grade, polish, symmetry, fluorescence_intensity, country, supplier, 
		price_no_added_value, price_retail, featured, recommend_words, images, extra_words, status,
		 ordered_by, picked_up, sold_price, profitable 
	 FROM diamonds WHERE stock_ref='%s' OR certificate_number='%s'`,
		ref, ref)
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ds, err := composeDiamond(rows)
	if err != nil {
		return nil, err
	}
	return ds, nil
}

func filterDiamonds(c *gin.Context) ([]diamond, error) {
	q, err := composeFilterDiamondsQuery(c)
	if err != nil {
		return nil, err
	}
	rows, err := dbQuery(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ds, err := composeDiamond(rows)
	if err != nil {
		return nil, err
	}
	return ds, nil
}

// shape: shape = "BR"  OR  shape = "PS"
// color:
// clarity:
// cut:
// polish:
// sym:
// fluo:
// place:
// certi:
// weight_from:
// weight_to:
// price_from:
// price_to:
// sorting:
// sorting_direction:ASC
// crr_page:1
// vat_choice:NO /
func composeFilterDiamondsQuery(c *gin.Context) (string, error) {
	var querys []string
	if c.PostForm("shape") != "" {
		querys = append(querys, strings.ToUpper(c.PostForm("shape")))
	}

	if c.PostForm("color") != "" {
		querys = append(querys, strings.ToUpper(c.PostForm("color")))
	}

	if c.PostForm("clarity") != "" {
		querys = append(querys, strings.ToUpper(c.PostForm("clarity")))
	}
	if c.PostForm("cut") != "" {
		querys = append(querys, strings.ToUpper(c.PostForm("cut")))
	}

	if c.PostForm("polish") != "" {
		querys = append(querys, strings.ToUpper(c.PostForm("polish")))
	}

	if c.PostForm("sym") != "" {
		querys = append(querys, strings.ToUpper(c.PostForm("sym")))
	}

	if c.PostForm("fluo") != "" {
		querys = append(querys, strings.ToUpper(c.PostForm("fluo")))
	}

	if c.PostForm("place") != "" {
		querys = append(querys, strings.ToUpper(c.PostForm("place")))
	}

	if c.PostForm("certi") != "" {
		querys = append(querys, strings.ToUpper(c.PostForm("certi")))
	}

	var caratFrom, caratTo float64
	caratFrom = 0
	caratTo = 100
	if c.PostForm("weight_from") != "" {
		cValue, err := strconv.ParseFloat(c.PostForm("weight_from"), 64)
		if err != nil {
			return "", err
		}
		if cValue == 0 {
			caratFrom = math.Abs(cValue) - 0.01
		}
	}

	if c.PostForm("weight_to") != "" {
		cValue, err := strconv.ParseFloat(c.PostForm("weight_from"), 64)
		if err != nil {
			return "", err
		}
		if cValue == 0 {
			caratTo = math.Abs(cValue) + 0.01
		}
	}
	querys = append(querys, fmt.Sprintf("carat>= %f AND carat<= %f", caratFrom, caratTo))

	var priceFrom, priceTo float64
	priceFrom = 0
	priceTo = 99999
	if c.PostForm("price_from") != "" {
		cValue, err := strconv.ParseFloat(c.PostForm("weight_from"), 64)
		if err != nil {
			return "", err
		}
		if cValue == 0 {
			caratTo = cValue - 0.01
		}
		priceFrom = math.Abs(cValue)
	}

	if c.PostForm("price_to") != "" {
		cValue, err := strconv.ParseFloat(c.PostForm("weight_from"), 64)
		if err != nil {
			return "", err
		}
		if cValue == 0 {
			caratTo = cValue - 0.01
		}
		priceTo = math.Abs(cValue)
	}

	//TODO tax(contains or not)
	if c.PostForm("vat_choice") != "" {
		if strings.ToUpper(c.PostForm("vat_choice")) == "YES" {
			priceFrom = math.Floor(priceFrom / 1.2)
			priceTo = math.Ceil(priceTo / 1.2)
		}
	}
	querys = append(querys, fmt.Sprintf("price_retail between %f AND %f", caratFrom, caratTo))

	//current page
	//The SQL query below says "return only 10 records, start on record 16 (OFFSET 15)":
	//$sql = "SELECT * FROM Orders LIMIT 10 OFFSET 15";
	var limit string
	currentPage := 1
	if c.PostForm("crr_page") != "" {
		var err error
		currentPage, err = strconv.Atoi(c.PostForm("crr_page"))
		if err != nil {
			return "", err
		}
		//28 records per page
		limit = fmt.Sprintf("LIMIT 28 OFFSET %d", util.AbsInt(currentPage-1)*28)
	}

	direction := "ASC"
	if c.PostForm("sorting_direction") != "" {
		direction = strings.ToUpper(c.PostForm("sorting_direction"))
	}

	sort := sortDiamondsByQuery(c.PostForm("sorting"), direction)
	q := fmt.Sprintf(`SELECT id, diamond_id, stock_ref, shape, carat, color, clarity, grading_lab, 
		certificate_number, cut_grade, polish, symmetry, fluorescence_intensity, country, supplier, 
		price_no_added_value, price_retail, featured, recommend_words, images, extra_words, status,
		 ordered_by, picked_up, sold_price, profitable 
	 FROM diamonds WHERE (%s) %s %s`, strings.Join(querys, ") AND ("), limit, sort)
	util.Traceln(q)
	return q, nil
}

func sortDiamondsByQuery(sortBy, direction string) string {
	switch sortBy {
	case "weight":
		return fmt.Sprintf("ORDER BY carat %s, supplier ASC, price_retail ASC", direction)

	case "color":
		return fmt.Sprintf("ORDER BY color %s, supplier ASC, price_retail ASC", direction)

	case "clarity":
		clarityFields := strings.Join(VALID_CLARITY, "','")
		return fmt.Sprintf("ORDER BY Field(clarity, '%s')  %s, supplier ASC, price_retail ASC", clarityFields, direction)
	case "price":
		return fmt.Sprintf("ORDER BY price_retail %s, supplier ASC", direction)

	default:
		return "ORDER BY supplier ASC, price_retail ASC"
	}
}
