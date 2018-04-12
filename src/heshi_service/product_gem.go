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

type gem struct {
	ID               string    `json:"id"`
	StockID          string    `json:"stock_id"`
	Shape            string    `json:"shape"`
	Material         string    `json:"material"`
	Name             string    `json:"name "`
	Size             float64   `json:"size"`
	SizeStr          string    `json:"-"`
	Text             string    `json:"text"`
	Images           []string  `json:"images"`
	Certificate      string    `json:"certificate"`
	Status           string    `json:"status"`
	Verified         string    `json:"verified"`
	Featured         string    `json:"featured"`
	Price            float64   `json:"price"`
	PriceStr         string    `json:"-"`
	StockQuantity    int       `json:"stock_quantity"`
	StockQuantityStr string    `json:"-"`
	Profitable       string    `json:"profitable"`
	TotallyScanned   int       `json:"totally_scanned"`
	FreeAcc          string    `json:"free_acc"`
	LastScanAt       time.Time `json:"last_scan_at"`
	OfflineAt        time.Time `json:"offline_at"`
	promotion
}

func newGems(c *gin.Context) {
	imageFileNames, vemsg, err := validateUploadedMultipleFile(c, "gem", "image", int64(imageSizeLimit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if vemsg != (errors.HSMessage{}) {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	g := gem{
		ID:               newV4(),
		StockID:          strings.ToUpper(c.PostForm("stock_id")),
		Shape:            FormatInputString(c.PostForm("shape")),
		Material:         strings.ToUpper(c.PostForm("material")),
		Name:             strings.ToUpper(c.PostForm("name")),
		SizeStr:          c.PostForm("size"),
		Text:             c.PostForm("text"),
		Images:           imageFileNames,
		Certificate:      strings.ToUpper(c.PostForm("certificate")),
		Status:           strings.ToUpper(c.PostForm("status")),
		Verified:         strings.ToUpper(c.PostForm("verified")),
		Featured:         strings.ToUpper(c.PostForm("featured")),
		PriceStr:         c.PostForm("price"),
		StockQuantityStr: c.PostForm("stock_quantity"),
		Profitable:       strings.ToUpper(c.PostForm("profitable")),
		FreeAcc:          strings.ToUpper(c.PostForm("free_acc")),
	}

	if vemsg, err := g.validateGemReq(false); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsg) != 0 {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	if err := saveUploadedMultipleFile(c, "gem", "image", imageFileNames); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	q := g.composeInsertQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, g.ID)
}

//TODO what to update; stockid???
func updateGems(c *gin.Context) {
	uid := c.MustGet("id").(string)
	gid := c.Param("id")
	if exist, err := isGemExistByID(gid); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if !exist {
		c.JSON(http.StatusBadRequest, "Item doesn't exist")
		return
	}
	imageFileNames, vemsg, err := validateUploadedMultipleFile(c, "gem", "image", int64(imageSizeLimit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if vemsg != (errors.HSMessage{}) {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	g := gem{
		ID:               gid,
		StockID:          strings.ToUpper(c.PostForm("stock_id")),
		Shape:            FormatInputString(c.PostForm("shape")),
		Material:         strings.ToUpper(c.PostForm("material")),
		Name:             strings.ToUpper(c.PostForm("name")),
		SizeStr:          c.PostForm("size"),
		Text:             c.PostForm("text"),
		Images:           imageFileNames,
		Certificate:      strings.ToUpper(c.PostForm("certificate")),
		Status:           strings.ToUpper(c.PostForm("status")),
		Verified:         strings.ToUpper(c.PostForm("verified")),
		Featured:         strings.ToUpper(c.PostForm("featured")),
		PriceStr:         c.PostForm("price"),
		StockQuantityStr: c.PostForm("stock_quantity"),
		Profitable:       strings.ToUpper(c.PostForm("profitable")),
		FreeAcc:          strings.ToUpper(c.PostForm("free_acc")),
	}
	if vemsg, err := g.validateGemReq(true); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsg) != 0 {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	if err := saveUploadedMultipleFile(c, "gem", "image", imageFileNames); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	q := g.composeUpdateQueryTrack(uid)
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, g.ID)
	// go newHistoryRecords(uid, "gems", g.ID, g.parmsKV())
}

func getAllGems(c *gin.Context) {
	q := selectGemQuery("")
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ds, err := composeGem(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ds)
}

func getGem(c *gin.Context) {
	q := selectGemQuery(c.Param("id"))
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ds, err := composeGem(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if ds == nil {
		vemsgNotExist.Message = fmt.Sprintf("Fail to find gem with id: %s", c.Param("id"))
		c.JSON(http.StatusOK, vemsgNotExist)
		return
	}
	c.JSON(http.StatusOK, ds)
}

func composeGem(rows *sql.Rows) ([]gem, error) {
	var id, stockID, shape, status, material, name, text, certificate, verified, featured, profitable, freeAcc string
	var images sql.NullString
	var size, price float64
	var stockQuantity, totallyScanned int
	var lastScanAt time.Time
	var offlineAt sql_patch.NullTime
	var pid, promType, pstatus sql.NullString
	var promPrice sql.NullFloat64
	var promDiscount sql.NullInt64
	var beginAt, endAt sql_patch.NullTime

	var gs []gem
	for rows.Next() {
		if err := rows.Scan(&id, &stockID, &shape, &material, &size, &name,
			&text, &images, &certificate, &status, &verified, &featured, &price, &stockQuantity, &profitable,
			&totallyScanned, &freeAcc, &lastScanAt, &offlineAt,
			&pid, &promType, &promDiscount, &promPrice, &beginAt, &endAt, &pstatus); err != nil {
			return nil, err
		}
		g := gem{
			ID:             id,
			StockID:        stockID,
			Shape:          shape,
			Material:       material,
			Size:           size,
			Name:           name,
			Text:           text,
			Certificate:    certificate,
			Status:         status,
			Verified:       verified,
			Featured:       featured,
			Price:          price,
			StockQuantity:  stockQuantity,
			Profitable:     profitable,
			TotallyScanned: totallyScanned,
			FreeAcc:        freeAcc,
			LastScanAt:     lastScanAt,
			OfflineAt:      offlineAt.Time,
		}
		if images.String != "" {
			for _, image := range strings.Split(images.String, ";") {
				g.Images = append(g.Images, image)
			}
		}
		if pid.String != "" && pstatus.String == "ACTIVE" && endAt.Time.After(beginAt.Time) && endAt.Time.After(time.Now().UTC()) && beginAt.Time.Before(time.Now()) {
			b := beginAt.Time
			e := endAt.Time
			g.PromType = promType.String
			g.PromDiscount = int(promDiscount.Int64)
			g.PromPrice = promPrice.Float64
			g.BeginAt = &b
			g.EndAt = &e
		}
		gs = append(gs, g)
	}
	return gs, nil
}

func selectGemQuery(id string) string {
	q := `SELECT gems.id, stock_id, shape, material, size, name, text, images, certificate, 
	gems.status, verified, featured, price, stock_quantity, profitable, 
	totally_scanned, free_acc, last_scan_at,offline_at, 
	promotions.id, prom_type, prom_discount, prom_price, begin_at, end_at, promotions.status 
	FROM gems 
	LEFT JOIN promotions ON gems.promotion_id=promotions.id 
	WHERE gems.status IN ('AVAILABLE','OFFLINE')`

	if id != "" {
		q = fmt.Sprintf("%s AND gems.id='%s'", q, id)
	}
	return q
}
