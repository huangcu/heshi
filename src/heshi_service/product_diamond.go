package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type diamond struct {
	ID                    string  `json:"id"`
	StockRef              string  `json:"stock_ref"`
	Shape                 string  `json:"shape"`
	Carat                 float64 `json:"carat"`
	Color                 string  `json:"color"`
	Clarity               string  `json:"clarity"`
	GradingLab            int     `json:"grading_lab"`
	CertificateNumber     string  `json:"certificate_number"`
	CutGrade              string  `json:"cut_grade"`
	Polish                string  `json:"polish"`
	Symmetry              string  `json:"symmetry"`
	FluorescenceIntensity string  `json:"fluorescence_intensity"`
	Country               string  `json:"country"`
	Supplier              string  `json:"supplier"`
	PriceNoAddedValue     float64 `json:"price_no_added_value"`
	PriceRetail           float64 `json:"price_retail"`
	CertificateLink       string  `json:"certificate_link"`
	ClarityNumber         string  `json:"clarity_number"`
	CutNumber             string  `json:"cut_number"`
	Featured              string  `json:"featured"`
	RecommandWords        string  `json:"recommand_words"`
	ExtraWords            string  `json:"extra_words"`
	Status                string  `json:"status"`
	OrderedBy             int64   `json:"ordered_by"`
	PickedUp              string  `json:"picked_up"`
	Sold                  string  `json:"sold"`
	SoldPrice             float64 `json:"sold_price"`
	Profitable            string  `json:"profitable"`
}

func getAllDiamonds(c *gin.Context) {
	q := selectDiamondQuery("")
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ds, err := composeDiamond(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, ds)
}

func getDiamond(c *gin.Context) {
	q := selectDiamondQuery(c.Param("id"))
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	ds, err := composeDiamond(rows)
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

//TODO
func newDiamond(c *gin.Context) {
	cValue, err := strconv.ParseFloat(c.PostForm("carat"), 64)
	if err != nil {
		c.JSON(http.StatusOK, "invalid carat input")
		return
	}
	gValue, err := strconv.Atoi(c.PostForm("grading_lab"))
	if err != nil {
		c.JSON(http.StatusOK, "invalid grading_lab input")
		return
	}
	pValue, err := strconv.ParseFloat(c.PostForm("price_no_added_value"), 64)
	if err != nil {
		c.JSON(http.StatusOK, "invalid carat input")
		return
	}
	prValue, err := strconv.ParseFloat(c.PostForm("price_retail"), 64)
	if err != nil {
		c.JSON(http.StatusOK, "invalid price_retail input")
		return
	}
	oValue, err := strconv.ParseInt(c.PostForm("ordered_by"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, "invalid ordered_by input")
		return
	}
	sValue, err := strconv.ParseFloat(c.PostForm("sold_price"), 64)
	if err != nil {
		c.JSON(http.StatusOK, "invalid sold_price input")
		return
	}

	d := diamond{
		ID:                    uuid.NewV4().String(),
		StockRef:              c.PostForm("stock_ref"),
		Shape:                 c.PostForm("shape"),
		Carat:                 cValue,
		Color:                 c.PostForm("color"),
		Clarity:               c.PostForm("clarity"),
		GradingLab:            gValue,
		CertificateNumber:     c.PostForm("certificate_number"),
		CutGrade:              c.PostForm("cut_grade"),
		Polish:                c.PostForm("polish"),
		Symmetry:              c.PostForm("symmetry"),
		FluorescenceIntensity: c.PostForm("fluorescence_intensity"),
		Country:               c.PostForm("country"),
		Supplier:              c.PostForm("supplier"),
		PriceNoAddedValue:     pValue,
		PriceRetail:           prValue,
		CertificateLink:       c.PostForm("certificate_link"),
		ClarityNumber:         c.PostForm("clarity_number"),
		CutNumber:             c.PostForm("cut_number"),
		Featured:              c.PostForm("featured"),
		RecommandWords:        c.PostForm("recommand_words"),
		ExtraWords:            c.PostForm("extra_words"),
		Status:                c.PostForm("status"),
		OrderedBy:             oValue,
		PickedUp:              c.PostForm("picked_up"),
		Sold:                  c.PostForm("sold"),
		SoldPrice:             sValue,
		Profitable:            c.PostForm("profitable"),
	}
	vmsg := d.validateDiamondReq()
	if vmsg != "" {
		c.JSON(http.StatusOK, vmsg)
		return
	}
}

func composeDiamond(rows *sql.Rows) ([]diamond, error) {
	var id, stockRef, country, supplier, clarityNumber, cutNumber, featured, status, pickedUp, sold, profitable string
	var shape, color, clarity, certificateNumber, cutGrade, polish, symmetry, fluorescenceIntensity sql.NullString
	var certificateLink, recommandWords, extraWords sql.NullString
	var soldPrice sql.NullFloat64
	var orderedBy sql.NullInt64
	var carat, priceNoAddedValue, priceRetail float64
	var gradingLab int

	var ds []diamond
	for rows.Next() {
		if err := rows.Scan(&id, &stockRef, &shape, &carat, &color, &clarity, &gradingLab, &certificateNumber,
			&cutGrade, &polish, &symmetry, &fluorescenceIntensity, &country, &supplier, &priceNoAddedValue, &priceRetail,
			&certificateLink, &clarityNumber, &cutNumber, &featured, &recommandWords, &extraWords, &status, &orderedBy, &pickedUp,
			&sold, &soldPrice, &profitable); err != nil {
			return nil, err
		}
		d := diamond{id, stockRef, shape.String, carat, color.String, clarity.String, gradingLab, certificateNumber.String,
			cutGrade.String, polish.String, symmetry.String, fluorescenceIntensity.String, country, supplier, priceNoAddedValue,
			priceRetail, certificateLink.String, clarityNumber, cutNumber, featured, recommandWords.String, extraWords.String,
			status, orderedBy.Int64, pickedUp, sold, soldPrice.Float64, profitable}
		ds = append(ds, d)
	}
	return ds, nil
}

func selectDiamondQuery(id string) string {
	q := `SELECT id, stock_ref, shape, carat, color, clarity, grading_lab, certificate_number, cut_grade,
	 polish, symmetry, fluorescence_intensity, country, supplier, price_no_added_value, price_retail, 
	 certificate_link, clarity_number, cut_number, featured, recommand_words, extra_words, status, ordered_by, picked_up, 
	 sold, sold_price, profitable FROM diamonds`

	if id != "" {
		q = fmt.Sprintf("%s WHERE id='%s'", q, id)
	}
	return q
}
