package main

import (
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
	filename := file.Filename + time.Now().Format(timeFormatFileName)
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
	uid := c.MustGet("id").(string)
	category := strings.ToUpper(c.PostForm("category"))
	jewelrySubCategory := strings.ToUpper(c.PostForm("jewelryCategory"))
	if !util.IsInArrayString(category, VALID_PRODUCTS) {
		vemsgUploadProductsCategoryNotValid.Message = fmt.Sprintf("%s is not valid product type", category)
		c.JSON(http.StatusOK, vemsgUploadProductsCategoryNotValid)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}
	// Upload the file to specific dst.
	if err := os.MkdirAll(filepath.Join(UPLOADFILEDIR, category, uid), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	// TODO must be excel
	var filename string
	exts := strings.SplitN(file.Filename, ".", 2)
	if len(exts) == 2 {
		filename = exts[0] + time.Now().Format(timeFormatFileName) + exts[1]
	} else {
		filename = file.Filename + time.Now().Format(timeFormatFileName)
	}
	dst := filepath.Join(UPLOADFILEDIR, category, uid, filename)
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

	missingHeaders := validateHeaders(category, headers)
	if len(missingHeaders) != 0 {
		c.JSON(http.StatusOK, gin.H{"missing-headers": missingHeaders})
		return
	}

	go func() {
		//here to track, who uploaded which file, and and filename saved on disk
		p := productStockHandleRecord{
			ID:             newV4(),
			UserID:         uid,
			Category:       category,
			Action:         "UPLOAD STOCK",
			Filename:       file.Filename,
			FileNameOnDisk: filename,
		}
		p.newProductStockHanldeRecords()
	}()
	importProducts(uid, category, dst, jewelrySubCategory)
}

func importProducts(uid, category, file, cate string) ([]util.Row, error) {
	switch strings.ToUpper(category) {
	case DIAMOND:
		return importDiamondProducts(uid, file)
	case SMALLDIAMOND:
		return importSmallDiamondProducts(file)
	case JEWELRY:
		return importJewelryProducts(uid, file, cate)
	case GEM:
		return importGemProducts(uid, file)
	default:
		return nil, errors.Newf("product category:%s not right", category)
	}
}

func importSmallDiamondProducts(file string) ([]util.Row, error) {
	records, err := util.ParseCSVToStruct(file)
	if err != nil {
		return nil, err
	}
	if len(records) < 1 {
		return nil, errors.New("uploaded file has no records")
	}
	ignoredRows := []util.Row{}
	//get headers
	originalHeaders := records[0]

	//process records
	util.Println("start process small diamond")
	for index := 1; index < len(records); index++ {
		sd := smallDiamond{}
		row := records[index]
		record := row.Value
		util.Printf("processsing row: %d, %s", index, record)
		for i, header := range originalHeaders.Value {
			switch header {
			case "size_from":
				sValue, err := util.StringToFloat(record[i])
				if err != nil {
					row.Message = append(row.Message, errors.GetMessage(err))
					row.Ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					row.Message = append(row.Message, "size from cannot be 0")
					row.Ignored = true
				}
				sd.SizeFrom = sValue
			case "size_to":
				sValue, err := util.StringToFloat(record[i])
				if err != nil {
					row.Message = append(row.Message, errors.GetMessage(err))
					row.Ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					row.Message = append(row.Message, "size to cannot be 0")
					row.Ignored = true
				}
				sd.SizeTo = sValue
			case "price":
				sValue, err := util.StringToFloat(record[i])
				if err != nil {
					row.Message = append(row.Message, errors.GetMessage(err))
					row.Ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					row.Message = append(row.Message, "price cannot be 0")
					row.Ignored = true
				}
				sd.Price = sValue
			case "quantity":
				sValue, err := strconv.Atoi(record[i])
				if err != nil {
					row.Message = append(row.Message, errors.GetMessage(err))
					row.Ignored = true
				}
				//value cannot be 0
				if sValue == 0 {
					row.Message = append(row.Message, "quantity cannot be 0")
					row.Ignored = true
				}
				sd.Quantity = util.AbsInt(sValue)
			}
		}

		//insert into db
		if !row.Ignored {
			q := `INSERT INTO small_diamonds (id, size_from, size_to, price, quantity) VALUSE('%s', '%f', '%f', '%f', '%d')`
			if _, err := dbExec(fmt.Sprintf(q, newV4()), sd.SizeFrom, sd.SizeTo, sd.Price, sd.Quantity); err != nil {
				return nil, err
			}
		}
	}
	util.Println("finish process small diamond")
	return ignoredRows, nil
}

func validateHeaders(product string, headers []string) []string {
	switch product {
	case DIAMOND:
		return validateDiamondHeaders(headers)
	case SMALLDIAMOND:
		return validateSmallDiamondHeaders(headers)
	case JEWELRY:
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
