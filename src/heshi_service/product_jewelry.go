package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"
	"sql_patch"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type jewelry struct {
	ID               string    `json:"id"`
	StockID          string    `json:"stock_id"`
	Name             string    `json:"name"`
	NeedDiamond      string    `json:"need_diamond"`
	Category         string    `json:"category"`
	MountingType     string    `json:"mounting_type"`
	Material         string    `json:"material"`
	MetalWeight      float64   `json:"metal_weight"`
	MetalWeightStr   string    `json:"-"`
	DiaShape         string    `json:"dia_shape"`
	UnitNumber       string    `json:"unit_number"`
	DiaSizeMin       float64   `json:"dia_size_min"`
	DiaSizeMinStr    string    `json:"-"`
	DiaSizeMax       float64   `json:"dia_size_max"`
	DiaSizeMaxStr    string    `json:"-"`
	MainDiaNum       int64     `json:"main_dia_num"`
	MainDiaNumStr    string    `json:"-"`
	MainDiaSize      float64   `json:"main_dia_size"`
	MainDiaSizeStr   string    `json:"-"`
	SmallDias        string    `json:"small_dias"`
	SmallDiaNum      int64     `json:"small_dia_num"`
	SmallDiaNumStr   string    `json:"-"`
	SmallDiaCarat    float64   `json:"small_dia_carat"`
	SmallDiaCaratStr string    `json:"-"`
	Price            float64   `json:"price"`
	PriceStr         string    `json:"-"`
	VideoLink        string    `json:"video_link"`
	Images           []string  `json:"images"`
	Text             string    `json:"text"`
	Status           string    `json:"status"`
	Verified         string    `json:"verified"`
	Featured         string    `json:"featured"`
	StockQuantity    int       `json:"stock_quantity"`
	StockQuantityStr string    `json:"-"`
	Profitable       string    `json:"profitable"`
	TotallyScanned   int       `json:"totally_scanned"`
	FreeAcc          string    `json:"free_acc"`
	LastScanAt       time.Time `json:"last_scan_at"`
	OfflineAt        time.Time `json:"offline_at"`
	promotion
}

func getAllJewelrys(c *gin.Context) {
	q := selectJewelryQuery("")
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ds, err := composeJewelry(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ds)
}

func getJewelry(c *gin.Context) {
	q := selectJewelryQuery(c.Param("id"))
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ds, err := composeJewelry(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if ds == nil {
		vemsgNotExist.Message = fmt.Sprintf("Fail to find jewelry with id: %s", c.Param("id"))
		c.JSON(http.StatusOK, vemsgNotExist)
		return
	}
	c.JSON(http.StatusOK, ds)
}

func newJewelry(c *gin.Context) {
	fileHeader, _ := c.FormFile("video")
	filename, vemsg, err := validateUploadedSingleFile(fileHeader, "jewelry", "video", int64(videoSizeLimit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if vemsg != (errors.HSMessage{}) {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	imageFileNames, vemsg, err := validateUploadedMultipleFile(c, "jewelry", "image", int64(imageSizeLimit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if vemsg != (errors.HSMessage{}) {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	j := jewelry{
		ID:               newV4(),
		StockID:          strings.ToUpper(c.PostForm("stock_id")),
		Name:             c.PostForm("name"),
		Category:         c.PostForm("category"),
		NeedDiamond:      strings.ToUpper(c.PostForm("need_diamond")),
		MountingType:     strings.ToUpper(c.PostForm("mounting_type")),
		Material:         strings.ToUpper(c.PostForm("material")),
		MetalWeightStr:   c.PostForm("metal_weight"),
		DiaShape:         FormatInputString(c.PostForm("dia_shape")),
		UnitNumber:       c.PostForm("unit_number"),
		DiaSizeMinStr:    c.PostForm("dia_size_min"),
		DiaSizeMaxStr:    c.PostForm("dia_size_max"),
		MainDiaNumStr:    c.PostForm("main_dia_num"),
		MainDiaSizeStr:   c.PostForm("main_dia_size"),
		SmallDias:        strings.ToUpper(c.PostForm("small_dias")),
		SmallDiaNumStr:   c.PostForm("small_dia_num"),
		SmallDiaCaratStr: c.PostForm("small_dia_carat"),
		PriceStr:         c.PostForm("price"),
		VideoLink:        filename,
		Images:           imageFileNames,
		Text:             c.PostForm("text"),
		Status:           strings.ToUpper(c.PostForm("status")),
		Verified:         strings.ToUpper(c.PostForm("verified")),
		Featured:         strings.ToUpper(c.PostForm("featured")),
		Profitable:       strings.ToUpper(c.PostForm("profitable")),
		FreeAcc:          strings.ToUpper(c.PostForm("free_acc")),
		StockQuantityStr: c.PostForm("stock_quantity"),
	}
	if vemsg, err := j.validateJewelryReq(false); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsg) != 0 {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	if err := saveUploadedSingleFile(c, "jewelry", "video", filename); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if err := saveUploadedMultipleFile(c, "jewelry", "image", imageFileNames); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	q := j.composeInsertQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, j.ID)
}

func updateJewelry(c *gin.Context) {
	uid := c.MustGet("id").(string)
	jid := c.Param("id")
	if exist, err := isJewelryExistByID(jid); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if !exist {
		c.JSON(http.StatusBadRequest, "Item doesn't exist")
		return
	}
	fileHeader, _ := c.FormFile("video")
	filename, vemsg, err := validateUploadedSingleFile(fileHeader, "jewelry", "video", int64(videoSizeLimit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if vemsg != (errors.HSMessage{}) {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	imageFileNames, vemsg, err := validateUploadedMultipleFile(c, "jewelry", "image", int64(imageSizeLimit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if vemsg != (errors.HSMessage{}) {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	j := jewelry{
		ID:               jid,
		StockID:          strings.ToUpper(c.PostForm("stock_id")),
		Name:             c.PostForm("name"),
		Category:         c.PostForm("category"),
		NeedDiamond:      strings.ToUpper(c.PostForm("need_diamond")),
		MountingType:     strings.ToUpper(c.PostForm("mounting_type")),
		Material:         strings.ToUpper(c.PostForm("material")),
		MetalWeightStr:   c.PostForm("metal_weight"),
		DiaShape:         FormatInputString(c.PostForm("dia_shape")),
		UnitNumber:       c.PostForm("unit_number"),
		DiaSizeMinStr:    c.PostForm("dia_size_min"),
		DiaSizeMaxStr:    c.PostForm("dia_size_max"),
		MainDiaNumStr:    c.PostForm("main_dia_num"),
		MainDiaSizeStr:   c.PostForm("main_dia_size"),
		SmallDias:        strings.ToUpper(c.PostForm("small_dias")),
		SmallDiaNumStr:   c.PostForm("small_dia_num"),
		SmallDiaCaratStr: c.PostForm("small_dia_carat"),
		PriceStr:         c.PostForm("price"),
		VideoLink:        filename,
		Images:           imageFileNames,
		Text:             c.PostForm("text"),
		Status:           strings.ToUpper(c.PostForm("status")),
		Verified:         strings.ToUpper(c.PostForm("verified")),
		Featured:         strings.ToUpper(c.PostForm("featured")),
		Profitable:       strings.ToUpper(c.PostForm("profitable")),
		FreeAcc:          strings.ToUpper(c.PostForm("free_acc")),
		StockQuantityStr: c.PostForm("stock_quantity"),
	}
	if vemsg, err := j.validateJewelryReq(true); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsg) != 0 {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	if err := saveUploadedSingleFile(c, "jewelry", "video", filename); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if err := saveUploadedMultipleFile(c, "jewelry", "image", imageFileNames); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	q := j.composeUpdateQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, j.ID)
	go newHistoryRecords(uid, "jewelrys", j.ID, j.parmsKV())
}

func composeJewelry(rows *sql.Rows) ([]jewelry, error) {
	var id, stockID, category, needDiamond, name, status, verified, featured, profitable, freeAcc string
	var unitNumber, diaShape, material, smallDias, mountingType, videoLink, images, text sql.NullString
	var metalWeight, mainDiaSize, diaSizeMin, diaSizeMax, smallDiaCarat, price sql.NullFloat64
	var mainDiaNum, smallDiaNum sql.NullInt64
	var stockQuantity, totallyScanned int
	var lastScanAt time.Time
	var offlineAt sql_patch.NullTime

	var pid, promType, pstatus sql.NullString
	var promPrice sql.NullFloat64
	var promDiscount sql.NullInt64
	var beginAt, endAt sql_patch.NullTime
	var ds []jewelry
	for rows.Next() {
		if err := rows.Scan(&id, &stockID, &category, &unitNumber, &diaShape, &material, &metalWeight,
			&needDiamond, &name, &diaSizeMin, &diaSizeMax, &smallDias, &smallDiaNum, &smallDiaCarat,
			&mountingType, &mainDiaNum, &mainDiaSize, &videoLink, &images, &text, &status, &verified,
			&featured, &price, &stockQuantity, &profitable, &totallyScanned, &freeAcc, &lastScanAt, &offlineAt,
			&pid, &promType, &promDiscount, &promPrice, &beginAt, &endAt, &pstatus); err != nil {
			return nil, err
		}
		d := jewelry{
			ID:             id,
			StockID:        stockID,
			NeedDiamond:    needDiamond,
			Name:           name,
			Category:       category,
			MountingType:   mountingType.String,
			Material:       material.String,
			MetalWeight:    metalWeight.Float64,
			DiaShape:       diaShape.String,
			UnitNumber:     unitNumber.String,
			DiaSizeMin:     diaSizeMin.Float64,
			DiaSizeMax:     diaSizeMax.Float64,
			MainDiaNum:     mainDiaNum.Int64,
			MainDiaSize:    mainDiaSize.Float64,
			SmallDias:      smallDias.String,
			SmallDiaCarat:  smallDiaCarat.Float64,
			SmallDiaNum:    smallDiaNum.Int64,
			Price:          price.Float64,
			VideoLink:      videoLink.String,
			Text:           text.String,
			Status:         status,
			Verified:       verified,
			Featured:       featured,
			StockQuantity:  stockQuantity,
			Profitable:     profitable,
			TotallyScanned: totallyScanned,
			FreeAcc:        freeAcc,
			LastScanAt:     lastScanAt,
			OfflineAt:      offlineAt.Time,
		}
		if images.String != "" {
			for _, image := range strings.Split(images.String, ";") {
				d.Images = append(d.Images, image)
			}
		}
		if pid.String != "" && pstatus.String == "ACTIVE" && endAt.Time.After(beginAt.Time) && endAt.Time.After(time.Now().UTC()) && beginAt.Time.Before(time.Now()) {
			b := beginAt.Time
			e := endAt.Time
			d.PromType = promType.String
			d.PromDiscount = int(promDiscount.Int64)
			d.PromPrice = promPrice.Float64
			d.BeginAt = &b
			d.EndAt = &e
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func selectJewelryQuery(id string) string {
	q := `SELECT jewelrys.id, stock_id, category, unit_number, dia_shape, material, metal_weight, 
	need_diamond, name, dia_size_min, dia_size_max, small_dias, small_dia_num, small_dia_carat, 
	mounting_type, main_dia_num, main_dia_size, video_link, images, text, jewelrys.status, verified, 
	featured, price, stock_quantity, profitable, totally_scanned, free_acc, last_scan_at,offline_at, 
	promotions.id, prom_type, prom_discount, prom_price, begin_at, end_at, promotions.status 
	FROM jewelrys 
	LEFT JOIN promotions ON jewelrys.promotion_id=promotions.id 
	WHERE jewelrys.status IN ('AVAILABLE','OFFLINE')`

	if id != "" {
		q = fmt.Sprintf("%s AND jewelrys.id='%s'", q, id)
	}
	return q
}
