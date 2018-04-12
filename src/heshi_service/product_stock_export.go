package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"heshi/errors"
	"net/http"
	"os"
	"path/filepath"
	"sql_patch"
	"strings"
	"time"
	"util"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/gin-gonic/gin"
)

type onlineOfflineProduct struct {
	ItemID       string `json:"item_id"`
	ItemCategory string `json:"item_category"`
}

func onlineOfflineProducts(c *gin.Context) {
	updatedBy := c.MustGet("id").(string)
	action := c.Param("action")
	var products []onlineOfflineProduct
	if err := json.Unmarshal([]byte(c.PostForm("ids")), &products); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}
	switch strings.ToLower(action) {
	case "offline":
		if err := onlineOffline("OFFLINE", updatedBy, products); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	case "online":
		if err := onlineOffline("AVAILABLE", updatedBy, products); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
	default:
		c.JSON(http.StatusNotFound, "page not found")
		return
	}
	c.JSON(http.StatusOK, "SUCCESS")
}

func onlineOffline(status, updatedBy string, products []onlineOfflineProduct) error {
	//can re-online, offline and deleted product
	cstatus := `'OFFLINE','DELETED'`
	if status == "OFFLINE" {
		//offline only these online
		cstatus = `'AVAILABLE'`
	}
	smap := make(map[string]interface{})
	smap["status"] = status
	for _, product := range products {
		switch strings.ToUpper(product.ItemCategory) {
		case DIAMOND:
			q := fmt.Sprintf(`UPDATE diamonds SET status='%s' WHERE id='%s' AND status IN (%s)`,
				status, product.ItemID, cstatus)
			_, err := dbExec(q)
			if err != nil {
				return err
			}
			// NOT Track product status change
			// c, err := r.RowsAffected()
			// if err != nil {
			// 	return err
			// }
			// if int(c) == 1 {
			// 	go newHistoryRecords(updatedBy, "diamonds", product.ItemID, smap)
			// }
		case JEWELRY:
			q := fmt.Sprintf(`UPDATE jewelrys SET status='%s' WHERE id='%s' AND status IN (%s)`,
				status, product.ItemID, cstatus)
			_, err := dbExec(q)
			if err != nil {
				return err
			}
			// NOT Track product status change
			// c, err := r.RowsAffected()
			// if err != nil {
			// 	return err
			// }
			// if int(c) == 1 {
			// 	go newHistoryRecords(updatedBy, "jewelrys", product.ItemID, smap)
			// }
		case GEM:
			q := fmt.Sprintf(`UPDATE gems SET status='%s' WHERE id='%s' AND status IN (%s)`,
				status, product.ItemID, cstatus)
			_, err := dbExec(q)
			if err != nil {
				return err
			}
			// NOT Track product status change
			// c, err := r.RowsAffected()
			// if err != nil {
			// 	return err
			// }
			// if int(c) == 1 {
			// 	go newHistoryRecords(updatedBy, "gems", product.ItemID, smap)
			// }
		default:
			return errors.Newf("Item category: not right", product.ItemCategory)
		}
	}
	return nil
}

func exportProduct(c *gin.Context) {
	uid := c.MustGet("id").(string)
	category := strings.ToUpper(c.PostForm("category"))
	jewelrySubCategory := strings.ToUpper(c.PostForm("jewelryCategory"))
	if category == "" {
		c.JSON(http.StatusBadRequest, "must specify a product category")
		return
	}
	if jewelrySubCategory != "" {
		if !util.IsInArrayString(jewelrySubCategory, VALID_CATEGORY) {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("%s not a valid jewelry category", jewelrySubCategory))
			return
		}
	}
	if !util.IsInArrayString(category, VALID_PRODUCTS) {
		vemsgUploadProductsCategoryNotValid.Message = fmt.Sprintf("%s is not valid product type", category)
		c.JSON(http.StatusOK, vemsgUploadProductsCategoryNotValid)
		return
	}
	serveFile, err := exportProductFromDB(uid, category, jewelrySubCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	// TODO redirect to download Or need do it from webpage, just return the link
	c.JSON(http.StatusOK, serveFile)
	return
}

func exportProductFromDB(uid, category, jewelrySubCategory string) (string, error) {
	switch strings.ToUpper(category) {
	case DIAMOND:
		servePath, err := exportDiamondProducts(uid)
		if err != nil {
			return "", err
		}
		return servePath, nil
	case SMALLDIAMOND:

	case JEWELRY:
		servePath, err := exportJewelryProducts(uid, jewelrySubCategory)
		if err != nil {
			return "", err
		}
		return servePath, nil
	case GEM:
		servePath, err := exportGemProducts(uid)
		if err != nil {
			return "", err
		}
		return servePath, nil
	}
	return "", nil
}

func exportDiamondProducts(uid string) (string, error) {
	q := `SELECT id, diamond_id, stock_ref, shape, carat, color, clarity, grading_lab, 
	certificate_number, cut_grade, polish, symmetry, fluorescence_intensity, country, 
	supplier, price_no_added_value, price_retail, featured, recommend_words, extra_words, images,
	status, profitable, updated_at, created_at 
	FROM diamonds ORDER BY updated_at DESC`
	rows, err := dbQuery(q)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Get column names
	columnNames, err := rows.Columns()
	if err != nil {
		return "", err
	}
	xlsx := excelize.NewFile()
	// Create a new sheet.and set active
	xlsx.SetActiveSheet(xlsx.NewSheet("Sheet1"))
	xlsx.InsertRow(strings.Join(columnNames, ","), 0)
	fmt.Println(strings.Join(columnNames, ","))
	columns := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y"}
	for i, column := range columns {
		xlsx.SetCellValue("Sheet1", column+"1", columnNames[i])
	}

	var id, diamondID, stockRef, shape, color, country, supplier, gradingLab string
	var clarity, certificateNumber, cutGrade, polish, symmetry, fluorescenceIntensity string
	var featured, status, profitable string
	var images, recommendWords, extraWords sql.NullString
	var carat, priceNoAddedValue, priceRetail float64
	var updatedAt, createdAt time.Time
	index := 1
	for rows.Next() {
		if err := rows.Scan(&id, &diamondID, &stockRef, &shape, &carat, &color, &clarity, &gradingLab, &certificateNumber,
			&cutGrade, &polish, &symmetry, &fluorescenceIntensity, &country, &supplier, &priceNoAddedValue, &priceRetail,
			&featured, &recommendWords, &extraWords, &images, &status, &profitable,
			&updatedAt, &createdAt); err != nil {
			return "", err
		}
		index++
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "A", index), id)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "B", index), diamondID)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "C", index), stockRef)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "D", index), shape)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "E", index), carat)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "E", index), color)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "G", index), clarity)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "H", index), gradingLab)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "I", index), certificateNumber)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "J", index), cutGrade)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "K", index), polish)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "L", index), symmetry)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "M", index), fluorescenceIntensity)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "N", index), country)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "O", index), supplier)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "P", index), priceNoAddedValue)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "Q", index), priceRetail)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "R", index), featured)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "S", index), recommendWords.String)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "T", index), extraWords.String)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "U", index), images.String)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "V", index), status)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "W", index), profitable)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "X", index), updatedAt.Format(timeFormat))
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "Y", index), createdAt.Format(timeFormat))
	}

	dst := filepath.Join(UPLOADFILEDIR, DIAMOND, uid, "export"+time.Now().Format(timeFormat)+".xlsx")
	if err := os.MkdirAll(filepath.Join(UPLOADFILEDIR, DIAMOND, uid), 0755); err != nil {
		return "", err
	}
	servePath := strings.TrimLeft(dst, UPLOADFILEDIR+"/")
	if err := xlsx.SaveAs(dst); err != nil {
		return "", err
	}
	go func() {
		//here to track, who demand an export of product and filename saved on disk
		p := productStockHandleRecord{
			ID:             newV4(),
			UserID:         uid,
			Category:       DIAMOND,
			Action:         "EXPORT STOCK",
			Filename:       "",
			FileNameOnDisk: servePath,
		}
		p.newProductStockHanldeRecords()
	}()
	return servePath, nil
}

func exportJewelryProducts(uid, jewelrySubCategory string) (string, error) {
	q := `SELECT id, stock_id, category, unit_number, dia_shape, material, metal_weight, need_diamond, name, 
	 dia_size_min, dia_size_max, small_dias, small_dia_num, small_dia_carat, mounting_type, main_dia_num, 
	 main_dia_size, video_link, images, text, status, verified, featured, price, stock_quantity, 
	 profitable, totally_scanned, free_acc, last_scan_at,offline_at, updated_at, created_at
	 FROM jewelrys`
	if jewelrySubCategory != "" {
		q = fmt.Sprintf("%s WHERE category='%s'", q, jewelrySubCategory)
	}
	q = fmt.Sprintf("%s ORDER BY updated_at DESC ", q)

	rows, err := dbQuery(q)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Get column names
	columnNames, err := rows.Columns()
	if err != nil {
		return "", err
	}
	xlsx := excelize.NewFile()
	// Create a new sheet.and set active
	xlsx.SetActiveSheet(xlsx.NewSheet("Sheet1"))
	xlsx.InsertRow(strings.Join(columnNames, ","), 0)
	fmt.Println(strings.Join(columnNames, ","))
	columns := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N",
		"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF"}
	for i, column := range columns {
		xlsx.SetCellValue("Sheet1", column+"1", columnNames[i])
	}
	var id, stockID, category, needDiamond, name, status, verified, featured, profitable, freeAcc string
	var unitNumber, diaShape, material, smallDias, mountingType, videoLink, images, text sql.NullString
	var metalWeight, mainDiaSize, diaSizeMin, diaSizeMax, smallDiaCarat, price sql.NullFloat64
	var mainDiaNum, smallDiaNum sql.NullInt64
	var stockQuantity, totallyScanned int
	var lastScanAt, updatedAt, createdAt time.Time
	var offlineAt sql_patch.NullTime

	index := 1
	for rows.Next() {
		if err := rows.Scan(&id, &stockID, &category, &unitNumber, &diaShape, &material, &metalWeight, &needDiamond, &name,
			&diaSizeMin, &diaSizeMax, &smallDias, &smallDiaNum, &smallDiaCarat, &mountingType, &mainDiaNum,
			&mainDiaSize, &videoLink, &images, &text, &status, &verified, &featured, &price, &stockQuantity,
			&profitable, &totallyScanned, &freeAcc, &lastScanAt, &offlineAt, &updatedAt, &createdAt); err != nil {
			return "", err
		}
		index++
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "A", index), id)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "B", index), stockID)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "C", index), category)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "D", index), unitNumber.String)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "E", index), diaShape.String)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "E", index), material.String)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "G", index), metalWeight.Float64)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "H", index), needDiamond)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "I", index), name)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "J", index), diaSizeMin.Float64)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "K", index), diaSizeMax.Float64)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "L", index), smallDias.String)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "M", index), smallDiaNum.Int64)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "N", index), smallDiaCarat.Float64)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "O", index), mountingType.String)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "P", index), mainDiaNum.Int64)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "Q", index), mainDiaSize.Float64)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "R", index), videoLink.String)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "S", index), images.String)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "T", index), text.String)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "U", index), status)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "V", index), verified)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "W", index), featured)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "X", index), price.Float64)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "Y", index), stockQuantity)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "Z", index), profitable)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "AA", index), totallyScanned)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "AB", index), freeAcc)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "AC", index), lastScanAt.Format(timeFormat))
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "AD", index), offlineAt.Time.Format(timeFormat))
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "AE", index), updatedAt.Format(timeFormat))
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "AF", index), createdAt.Format(timeFormat))
	}

	dst := filepath.Join(UPLOADFILEDIR, JEWELRY, uid, "export"+time.Now().Format(timeFormat)+".xlsx")
	if err := os.MkdirAll(filepath.Join(UPLOADFILEDIR, JEWELRY, uid), 0755); err != nil {
		return "", err
	}
	servePath := strings.TrimLeft(dst, UPLOADFILEDIR+"/")
	if err := xlsx.SaveAs(dst); err != nil {
		return "", err
	}
	go func() {
		//here to track, who demand an export of product and filename saved on disk
		p := productStockHandleRecord{
			ID:             newV4(),
			UserID:         uid,
			Category:       JEWELRY,
			Action:         "EXPORT STOCK",
			Filename:       "",
			FileNameOnDisk: servePath,
		}
		p.newProductStockHanldeRecords()
	}()
	return servePath, nil
}

func exportGemProducts(uid string) (string, error) {
	q := `SELECT id, stock_id, shape, material, size, name, text, images, certificate, 
	 status, verified, featured, price, stock_quantity, profitable,
	 totally_scanned, free_acc, last_scan_at,offline_at,updated_at, created_at 
	 FROM gems ORDER BY updated_at DESC`
	rows, err := dbQuery(q)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Get column names
	columnNames, err := rows.Columns()
	if err != nil {
		return "", err
	}
	xlsx := excelize.NewFile()
	// Create a new sheet.and set active
	xlsx.SetActiveSheet(xlsx.NewSheet("Sheet1"))
	xlsx.InsertRow(strings.Join(columnNames, ","), 0)
	fmt.Println(strings.Join(columnNames, ","))
	columns := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N",
		"O", "P", "Q", "R", "S", "T", "U"}
	for i, column := range columns {
		xlsx.SetCellValue("Sheet1", column+"1", columnNames[i])
	}
	var id, stockID, shape, status, material, name, text, certificate, verified, inStock string
	var featured, profitable, freeAcc string
	var images sql.NullString
	var size, price float64
	var stockQuantity, totallyScanned int
	var lastScanAt, updatedAt, createdAt time.Time
	var offlineAt sql_patch.NullTime

	index := 1
	for rows.Next() {
		if err := rows.Scan(&id, &stockID, &shape, &material, &size, &name, &text, &images, &certificate,
			&status, &verified, &inStock, &featured, &price, &stockQuantity,
			&profitable, &totallyScanned, &freeAcc, &lastScanAt, &offlineAt, &updatedAt, &createdAt); err != nil {
			return "", err
		}
		index++
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "A", index), id)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "B", index), stockID)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "C", index), shape)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "D", index), material)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "E", index), size)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "E", index), name)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "G", index), text)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "H", index), images.String)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "I", index), certificate)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "J", index), status)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "K", index), verified)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "L", index), featured)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "M", index), price)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "N", index), stockQuantity)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "O", index), profitable)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "P", index), totallyScanned)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "Q", index), freeAcc)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "R", index), lastScanAt.Format(timeFormat))
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "S", index), offlineAt.Time.Format(timeFormat))
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "T", index), updatedAt.Format(timeFormat))
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", "U", index), createdAt.Format(timeFormat))
	}

	dst := filepath.Join(UPLOADFILEDIR, GEM, uid, "export"+time.Now().Format(timeFormat)+".xlsx")
	if err := os.MkdirAll(filepath.Join(UPLOADFILEDIR, GEM, uid), 0755); err != nil {
		return "", err
	}
	servePath := strings.TrimLeft(dst, UPLOADFILEDIR+"/")
	if err := xlsx.SaveAs(dst); err != nil {
		return "", err
	}
	go func() {
		//here to track, who demand an export of product and filename saved on disk
		p := productStockHandleRecord{
			ID:             newV4(),
			UserID:         uid,
			Category:       GEM,
			Action:         "EXPORT STOCK",
			Filename:       "",
			FileNameOnDisk: servePath,
		}
		p.newProductStockHanldeRecords()
	}()
	return servePath, nil
}
