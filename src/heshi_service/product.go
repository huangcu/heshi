package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
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

var PRODUCTS = []string{"diamond", "small_diamond", "jewelry"}
var diamondHeaders = []string{"diamond_id", "stock_ref", "shape", "carat", "color", "clarity", "grading_lab",
	"certificate_number", "cut_grade", "polish", "symmetry", "fluorescence_intensity", "country",
	"supplier", "price_no_added_value", "price_retail", "clarity_number", "cut_number"}

type product struct {
	Diamond      []diamond      `json:"diamond"`
	Jewelry      []jewelry      `json:"jewelry"`
	SmallDiamond []smallDiamond `json:"small_diamond"`
}

func getAllProducts(c *gin.Context) {

}

func uploadProducts(c *gin.Context) {
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

func processUploadedProducts(c *gin.Context) {
	id := c.MustGet("id").(string)
	product := c.PostForm("product")
	if !util.IsInArrayString(product, PRODUCTS) {
		c.JSON(http.StatusOK, fmt.Sprintf("%s is not valid product type", product))
	}
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
	return nil, nil
}
func importSmallDiamondProducts(file string) ([][]string, error) {
	return nil, nil
}
func importDiamondProducts(file string) ([][]string, error) {
	originalHeaders := []string{}
	records, err := util.ParseCSVToArrays(file)
	if err != nil {
		return nil, err
	}
	ignoredRows := [][]string{}
	//get headers
	for index := 0; index < len(records); index++ {
		if index == 0 {
			originalHeaders = records[0]
		}
	}
	for index := 0; index < len(records); index++ {
		//process records
		if index != 0 {
			ignored := false
			d := diamond{}
			record := records[index]
			fmt.Println("processsing " + strconv.Itoa(index))
			for _, header := range originalHeaders {
				for i := 0; i < len(originalHeaders); i++ {
					switch header {
					case "diamond_id":
						d.DiamondID = record[i]
					case "stock_ref":
						d.StockRef = record[i]
					case "shape":
						d.Shape = record[i]
					case "carat":
						cValue, err := strconv.ParseFloat(record[i], 64)
						if err != nil {
							ignoredRows = append(ignoredRows, record)
							ignored = true
						}
						if cValue == 0 {
							ignored = true
						}
						d.Carat = cValue
					case "color":
						d.Color = record[i]
					case "clarity":
						d.Clarity = record[i]
					case "grading_lab":
						cValue, err := strconv.Atoi(record[i])
						if err != nil {
							ignoredRows = append(ignoredRows, record)
							ignored = true
						}
						if cValue == 0 {
							ignored = true
						}
						d.GradingLab = cValue
					case "certificate_number":
						d.CertificateNumber = record[i]
					case "cut_grade":
						d.CutGrade = record[i]
					case "polish":
						d.Polish = record[i]
					case "symmetry":
						d.Symmetry = record[i]
					case "fluorescence_intensity":
						fmt.Println(record[i])
						d.FluorescenceIntensity = record[i]
					case "country":
						d.Country = record[i]
					case "supplier":
						d.Supplier = record[i]
					case "price_no_added_value":
						fmt.Println(record[i])
						cValue, err := strconv.ParseFloat(strings.Replace(record[i], ",", "", -1), 64)
						if err != nil {
							ignoredRows = append(ignoredRows, record)
							ignored = true
						}
						if cValue == 0 {
							ignored = true
						}
						d.PriceNoAddedValue = cValue
					case "price_retail":
						cValue, err := strconv.ParseFloat(strings.Replace(record[i], ",", "", -1), 64)
						if err != nil {
							ignoredRows = append(ignoredRows, record)
							ignored = true
						}
						if cValue == 0 {
							ignored = true
						}
						d.PriceRetail = cValue
					case "clarity_number":
						d.ClarityNumber = record[i]
					case "cut_number":
						d.CutNumber = record[i]
					}
				}
			}
			//insert into db
			if !ignored {
				fmt.Printf("%#v", d)
				var id string
				if err := db.QueryRow(fmt.Sprintf("SELECT id FROM diamonds WHERE stock_ref='%s'", d.StockRef)).Scan(&id); err != nil {
					if err == sql.ErrNoRows {
						d.ID = uuid.NewV4().String()
						q := d.composeInsertQuery()
						fmt.Println(q)
						if _, err := db.Exec(q); err != nil {
							return nil, err
							// ignoredRows = append(ignoredRows, record)
						}
					} else {
						// ignoredRows = append(ignoredRows, record)
						return nil, err
					}
				}

				d.ID = id
				q := d.composeUpdateQuery()
				if _, err := db.Exec(q); err != nil {
					// ignoredRows = append(ignoredRows, record)
					return nil, err
				}
			}
			fmt.Println("finish process")
		}
	}
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

func validateDiamondHeaders(headers []string) []string {
	var missingHeaders []string
	for _, header := range diamondHeaders {
		if !util.IsInArrayString(header, headers) {
			missingHeaders = append(missingHeaders, header)
		}
	}
	return missingHeaders
}

func validateSmallDiamondHeaders(headers []string) []string {
	var missingHeaders []string
	smallDiamondHeaders := []string{"diamond_id", "stock_ref", "shape", "carat", "color", "clarity", "grading_lab",
		"certificate_number", "cut_grade", "polish", "symmetry", "fluorescence_intensity", "country",
		"supplier", "price_no_added_value", "price_retail", "clarity_number", "cut_number"}

	for _, header := range smallDiamondHeaders {
		if !util.IsInArrayString(header, headers) {
			missingHeaders = append(missingHeaders, header)
		}
	}
	return missingHeaders
}

func validateJewelryHeaders(headers []string) []string {
	var missingHeaders []string
	jewelryHeaders := []string{"diamond_id", "stock_ref", "shape", "carat", "color", "clarity", "grading_lab",
		"certificate_number", "cut_grade", "polish", "symmetry", "fluorescence_intensity", "country",
		"supplier", "price_no_added_value", "price_retail", "clarity_number", "cut_number"}

	for _, header := range jewelryHeaders {
		if !util.IsInArrayString(header, headers) {
			missingHeaders = append(missingHeaders, header)
		}
	}
	return missingHeaders
}
