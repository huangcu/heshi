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
	Online           string    `json:"online"`
	Verified         string    `json:"verified"`
	InStock          string    `json:"in_stock"`
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
}

func newGems(c *gin.Context) {
	imageFileNames, vemsg, err := validateUploadedMultipleFile(c, "gem", "image", 1*1024*1024)
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
		Online:           strings.ToUpper(c.PostForm("online")),
		Verified:         strings.ToUpper(c.PostForm("verified")),
		InStock:          strings.ToUpper(c.PostForm("in_stock")),
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
	imageFileNames, vemsg, err := validateUploadedMultipleFile(c, "gem", "image", 1*1024*1024)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if vemsg != (errors.HSMessage{}) {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	id := c.Param("id")
	g := gem{
		ID:               id,
		StockID:          strings.ToUpper(c.PostForm("stock_id")),
		Shape:            FormatInputString(c.PostForm("shape")),
		Material:         strings.ToUpper(c.PostForm("material")),
		Name:             strings.ToUpper(c.PostForm("name")),
		SizeStr:          c.PostForm("size"),
		Text:             c.PostForm("text"),
		Images:           imageFileNames,
		Certificate:      strings.ToUpper(c.PostForm("certificate")),
		Online:           strings.ToUpper(c.PostForm("online")),
		Verified:         strings.ToUpper(c.PostForm("verified")),
		InStock:          strings.ToUpper(c.PostForm("in_stock")),
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
	q := g.composeUpdateQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, g.ID)
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
	var id, stockID, shape, online, material, name, text, certificate, verified, inStock, featured, profitable, freeAcc string
	var images sql.NullString
	var size, price float64
	var stockQuantity, totallyScanned int
	var lastScanAt time.Time
	var offlineAt sql_patch.NullTime

	var gs []gem
	for rows.Next() {
		if err := rows.Scan(&id, &stockID, &shape, &material, &size, &name,
			&text, &images, &certificate, &online, &verified, &inStock, &featured, &price, &stockQuantity, &profitable,
			&totallyScanned, &freeAcc, &lastScanAt, &offlineAt); err != nil {
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
			Online:         online,
			Verified:       verified,
			InStock:        inStock,
			Featured:       featured,
			Price:          price,
			StockQuantity:  stockQuantity,
			Profitable:     profitable,
			TotallyScanned: totallyScanned,
			FreeAcc:        freeAcc,
			LastScanAt:     lastScanAt.Local(),
			OfflineAt:      offlineAt.Time,
		}
		if images.String != "" {
			for _, image := range strings.Split(images.String, ";") {
				g.Images = append(g.Images, "image/gem/"+image)
			}
		}
		gs = append(gs, g)
	}
	return gs, nil
}

func selectGemQuery(id string) string {
	q := `SELECT id, stock_id, shape, material, size, name, text, images, certificate, 
	online, verified, in_stock, featured, price, stock_quantity, profitable,
	 totally_scanned, free_acc, last_scan_at,offline_at FROM gems`

	if id != "" {
		q = fmt.Sprintf("%s WHERE id='%s'", q, id)
	}
	return q
}
