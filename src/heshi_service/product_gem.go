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

type gem struct {
	ID             string    `json:"id"`
	StockID        string    `json:"stock_id"`
	Shape          string    `json:"shape"`
	Material       string    `json:"material"`
	Name           string    `json:"name "`
	Size           float64   `json:"size"`
	Text           string    `json:"text"`
	Online         string    `json:"online"`
	Verified       string    `json:"verified"`
	InStock        string    `json:"in_stock"`
	Featured       string    `json:"featured"`
	Price          float64   `json:"price"`
	StockQuantity  int       `json:"stock_quantity"`
	Profitable     string    `json:"profitable"`
	TotallyScanned int       `json:"totally_scanned"`
	FreeAcc        string    `json:"free_acc"`
	LastScanAt     time.Time `json:"last_scan_at"`
	OfflineAt      time.Time `json:"offline_at"`
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
		VEMSG_NOT_EXIST.Message = fmt.Sprintf("Fail to find gem with id: %s", c.Param("id"))
		c.JSON(http.StatusOK, VEMSG_NOT_EXIST)
		return
	}
	c.JSON(http.StatusOK, ds)
}

func composeGem(rows *sql.Rows) ([]gem, error) {
	var id, stockID, shape, online, material, name, text, verified, inStock, featured, profitable, freeAcc string
	var size, price float64
	var stockQuantity, totallyScanned int
	var lastScanAt time.Time
	var offlineAt sql_patch.NullTime

	var gs []gem
	for rows.Next() {
		if err := rows.Scan(&id, &stockID, &shape, &material, &size, &name,
			&text, &online, &verified, &inStock, &featured, &price, &stockQuantity, &profitable,
			&totallyScanned, &freeAcc, &lastScanAt, &offlineAt); err != nil {
			return nil, err
		}
		g := gem{ID: id, StockID: stockID, Shape: shape,
			Material: material, Size: size, Name: name,
			Text: text, Online: online, Verified: verified, InStock: inStock, Featured: featured, Price: price,
			StockQuantity: stockQuantity, Profitable: profitable, TotallyScanned: totallyScanned,
			FreeAcc: freeAcc, LastScanAt: lastScanAt.Local(), OfflineAt: offlineAt.Time}
		gs = append(gs, g)
	}
	return gs, nil
}

func selectGemQuery(id string) string {
	q := `SELECT id, stock_id, shape, material, size, name, text, 
	online, verified, in_stock, featured, price, stock_quantity, profitable,
	 totally_scanned, free_acc, last_scan_at,offline_at FROM gems`

	if id != "" {
		q = fmt.Sprintf("%s WHERE id='%s'", q, id)
	}
	return q
}
