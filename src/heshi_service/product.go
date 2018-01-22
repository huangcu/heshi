package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"util"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type product struct {
	Diamond      []diamond      `json:"diamond"`
	Jewelry      []jewelry      `json:"jewelry"`
	SmallDiamond []smallDiamond `json:"small_diamond"`
}

//TODO
func getAllProducts(c *gin.Context) {

}

//TODO customize header
func uploadAndGetFileHeaders(c *gin.Context) {
	id := c.MustGet("id").(string)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}
	// Upload the file to specific dst.
	filename := file.Filename + time.Now().Format("20060102150405")
	dst := filepath.Join(os.TempDir(), id, filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	headers, err := util.GetCSVHeaders(dst)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{filename: headers})
}

func uploadAndProcessProducts(c *gin.Context) {
	id := c.MustGet("id").(string)
	product := c.PostForm("product")
	if !util.IsInArrayString(product, VALID_PRODUCTS) {
		vemsgUploadProductsCategoryNotValid.Message = fmt.Sprintf("%s is not valid product type", product)
		c.JSON(http.StatusOK, vemsgUploadProductsCategoryNotValid)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}
	// Upload the file to specific dst.
	filename := file.Filename + time.Now().Format("20060102150405")

	if err := os.MkdirAll(filepath.Join(UPLOADFILEDIR, id), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	dst := filepath.Join(UPLOADFILEDIR, id, filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	headers, err := util.GetCSVHeaders(dst)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	missingHeaders := validateHeaders(product, headers)
	if len(missingHeaders) != 0 {
		c.JSON(http.StatusOK, gin.H{"missing-headers": missingHeaders})
		return
	}

	importProducts(product, dst)
}

func importProducts(product, file string) ([][]string, error) {
	switch product {
	case "diamond":
		return importDiamondProducts(file)
	case "small_diamond":
		return importSmallDiamondProducts(file)
	case "jewelry":
		return importJewelryProducts(file)
	default:
		return nil, nil
	}
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
			case "category":
				sValue, err := strconv.Atoi(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.Category = util.AbsInt(sValue)
			case "stock_id":
				j.StockID = record[i]
			case "unit_number":
				j.UnitNumber = record[i]
			case "dia_shape":
				j.DiaShape = record[i]
			case "need_diamond":
				j.NeedDiamond = record[i]
			case "metal_weight":
				cValue, err := strconv.ParseFloat(strings.Replace(record[i], ",", "", -1), 64)
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				if cValue == 0 {
					ignored = true
				}
				j.MetalWeight = cValue
			case "material":
				j.Material = record[i]
			case "name":
				j.Name = record[i]
			case "name_suffix":
				sValue, err := strconv.Atoi(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.NameSuffix = int64(util.AbsInt(sValue))
			case "dia_size_min":
				sValue, err := strconv.ParseFloat(record[i], 64)
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.DiaSizeMin = math.Abs(sValue)
			case "dia_size_max":
				sValue, err := strconv.ParseFloat(record[i], 64)
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.DiaSizeMax = math.Abs(sValue)
			case "small_dias":
				j.SmallDias = record[i]
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
				sValue, err := strconv.ParseFloat(record[i], 64)
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.SmallDiaCarat = math.Abs(sValue)
			case "mounting_type":
				j.MountingType = record[i]
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
				sValue, err := strconv.ParseFloat(record[i], 64)
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.MainDiaSize = math.Abs(sValue)
			case "video_link":
				j.VideoLink = record[i]
			case "text":
				j.Text = record[i]
			case "online":
				j.Online = record[i]
			case "verified":
				j.Verified = record[i]
			case "in_stock":
				j.InStock = record[i]
			case "featured":
				j.Featured = record[i]
			case "price":
				sValue, err := strconv.ParseFloat(record[i], 64)
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.Price = math.Abs(sValue)
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
				j.Profitable = record[i]
			case "totally_scanned":
				sValue, err := strconv.Atoi(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				j.TotallyScanned = util.AbsInt(sValue)
			case "free_acc":
				j.FreeAcc = record[i]
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

func importSmallDiamondProducts(file string) ([][]string, error) {
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
		sd := smallDiamond{}
		record := records[index]
		util.Printf("processsing row: %d, %s", index, record)
		for i, header := range originalHeaders {
			switch header {
			case "size_from":
				sValue, err := strconv.ParseFloat(record[i], 64)
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				sd.SizeFrom = math.Abs(sValue)
			case "size_to":
				sValue, err := strconv.ParseFloat(record[i], 64)
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				sd.SizeTo = math.Abs(sValue)
			case "price":
				sValue, err := strconv.ParseFloat(record[i], 64)
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				sd.Price = math.Abs(sValue)
			case "quantity":
				sValue, err := strconv.Atoi(record[i])
				if err != nil {
					ignoredRows = append(ignoredRows, record)
					ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					ignored = true
				}
				sd.Quantity = util.AbsInt(sValue)
			}
		}

		//insert into db
		if !ignored {
			q := `INSERT INTO small_diamonds (id, size_from, size_to, price, quantity) VALUSE('%s', '%f', '%f', '%f', '%d')`
			if _, err := dbExec(fmt.Sprintf(q, uuid.NewV4().String()), sd.SizeFrom, sd.SizeTo, sd.Price, sd.Quantity); err != nil {
				return nil, err
			}
		}
	}
	util.Println("finish process small diamond")
	return ignoredRows, nil
}

func validateHeaders(product string, headers []string) []string {
	switch product {
	case "diamond":
		return validateDiamondHeaders(headers)
	case "small_diamond":
		return validateSmallDiamondHeaders(headers)
	case "jewelry":
		return validateJewelryHeaders(headers)
	default:
		return nil
	}
}

func validateSmallDiamondHeaders(headers []string) []string {
	var missingHeaders []string
	for _, header := range smallDiamondHeaders {
		if !util.IsInArrayString(header, headers) {
			missingHeaders = append(missingHeaders, header)
		}
	}
	return missingHeaders
}

func validateJewelryHeaders(headers []string) []string {
	var missingHeaders []string
	for _, header := range jewelryHeaders {
		if !util.IsInArrayString(header, headers) {
			missingHeaders = append(missingHeaders, header)
		}
	}
	return missingHeaders
}
