package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"
	"sql_patch"
	"time"

	"github.com/gin-gonic/gin"
)

type jewelry struct {
	ID             string    `json:"id"`
	Category       int       `json:"category"`
	UnitNumber     string    `json:"unit_number"`
	DiaShape       string    `json:"dia_shape"`
	Material       string    `json:"material"`
	MetalWeight    float64   `json:"metal_weight"`
	GradingLab     int       `json:"grading_lab"`
	NeedDiamond    string    `json:"need_diamond"`
	Name           string    `json:"name "`
	NameSuffix     int64     `json:"name_suffix"`
	DiaSizeMin     float64   `json:"dia_size_min"`
	DiaSizeMax     float64   `json:"dia_size_max"`
	SmallDias      string    `json:"small_dias"`
	SmallDiaNum    int64     `json:"small_dia_num"`
	SmallDiaCarat  float64   `json:"small_dia_carat"`
	MountingType   string    `json:"mounting_type"`
	MainDiaNum     int       `json:"main_dia_num"`
	MainDiaSize    string    `json:"main_dia_size"`
	VideoLink      string    `json:"video_link"`
	Text           string    `json:"text"`
	Online         string    `json:"online"`
	Verified       string    `json:"verified"`
	InStock        string    `json:"in_stock"`
	Featured       string    `json:"featured"`
	Price          float64   `json:"price"`
	StockQuantity  int       `json:"stock_quantity"`
	Profitable     string    `json:"profitable"`
	ClarityNumber  string    `json:"clarity_number"`
	TotallyScanned int       `json:"totally_scanned"`
	FreeAcc        string    `json:"free_acc"`
	LastScanAt     time.Time `json:"last_scan_at"`
	OfflineAt      time.Time `json:"offline_at"`
}

func getAllJewelrys(c *gin.Context) {
	q := selectJewelryQuery("")
	rows, err := db.Query(q)
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
	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ds, err := composeJewelry(rows)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, fmt.Sprintf("Fail to find product with id: %s", c.Param("id")))
			return
		}
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ds)
}

func newJewelry(c *gin.Context) {
	q := selectJewelryQuery(c.Param("id"))
	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ds, err := composeJewelry(rows)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, fmt.Sprintf("Fail to find product with id: %s", c.Param("id")))
			return
		}
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ds)
}

func composeJewelry(rows *sql.Rows) ([]jewelry, error) {
	var id, unitNumber, needDiamond, name, online, verified, inStock, featured, profitable, clarityNumber, freeAcc string
	var diaShape, material, smallDias, mountingType, mainDiaSize, videoLink, text sql.NullString
	var metalWeight, price sql.NullFloat64
	var nameSuffix, smallDiaNum sql.NullInt64
	var diaSizeMin, diaSizeMax, smallDiaCarat float64
	var category, mainDiaNum, stockQuantity, totallyScanned int
	var lastScanAt time.Time
	var offlineAt sql_patch.NullTime

	var ds []jewelry
	for rows.Next() {
		if err := rows.Scan(&id, &category, &unitNumber, &diaShape, &material, &metalWeight, &needDiamond, &name, &nameSuffix,
			&diaSizeMin, &diaSizeMax, &smallDias, &smallDiaNum, &smallDiaCarat, &mountingType, &mainDiaNum, &mainDiaSize,
			&videoLink, &text, &online, &verified, &inStock, &featured, &price, &stockQuantity, &profitable, &clarityNumber,
			&totallyScanned, &freeAcc, &lastScanAt, &offlineAt); err != nil {
			return nil, err
		}
		d := jewelry{ID: id, Category: category, UnitNumber: unitNumber, DiaShape: diaShape.String,
			Material: material.String, MetalWeight: metalWeight.Float64, NeedDiamond: needDiamond, Name: name,
			NameSuffix: nameSuffix.Int64, DiaSizeMin: diaSizeMin, DiaSizeMax: diaSizeMax, SmallDias: smallDias.String,
			SmallDiaCarat: smallDiaCarat, SmallDiaNum: smallDiaNum.Int64, MountingType: mountingType.String,
			MainDiaNum: mainDiaNum, MainDiaSize: mainDiaSize.String, VideoLink: videoLink.String, Text: text.String,
			Online: online, Verified: verified, InStock: inStock, Featured: featured, Price: price.Float64,
			StockQuantity: stockQuantity, Profitable: profitable, ClarityNumber: clarityNumber, TotallyScanned: totallyScanned,
			FreeAcc: freeAcc, LastScanAt: lastScanAt.Local(), OfflineAt: offlineAt.Time}
		ds = append(ds, d)
	}
	return ds, nil
}

func selectJewelryQuery(id string) string {
	q := `SELECT id, category, unit_number, dia_shape, material, metal_weight, need_diamond, name, name_suffix,
	 dia_size_min, dia_size_max, small_dias, small_dia_num, small_dia_carat, mounting_type, main_dia_num, main_dia_size, 
	 video_link, text, online, verified, in_stock, featured, price, stock_quantity, profitable, clarity_number,
	 totally_scanned, free_acc, last_scan_at,offline_at FROM jewelrys`

	if id != "" {
		q = fmt.Sprintf("%s WHERE id='%s'", q, id)
	}
	return q
}
