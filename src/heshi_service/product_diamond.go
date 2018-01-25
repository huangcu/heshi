package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"util"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type diamond struct {
	ID                    string  `json:"id"`
	DiamondID             string  `json:"diamond_id"`
	StockRef              string  `json:"stock_ref"`
	Shape                 string  `json:"shape"`
	Carat                 float64 `json:"carat"`
	CaratStr              string  `json:"-"`
	Color                 string  `json:"color"`
	Clarity               string  `json:"clarity"`
	GradingLab            string  `json:"grading_lab"`
	CertificateNumber     string  `json:"certificate_number"`
	CutGrade              string  `json:"cut_grade"`
	Polish                string  `json:"polish"`
	Symmetry              string  `json:"symmetry"`
	FluorescenceIntensity string  `json:"fluorescence_intensity"`
	Country               string  `json:"country"`
	Supplier              string  `json:"supplier"`
	PriceNoAddedValue     float64 `json:"price_no_added_value"`
	PriceNoAddedValueStr  string  `json:"-"`
	PriceRetail           float64 `json:"price_retail"`
	PriceRetailStr        string  `json:"-"`
	CertificateLink       string  `json:"certificate_link"`
	Featured              string  `json:"featured"`
	RecommandWords        string  `json:"recommand_words"`
	ExtraWords            string  `json:"extra_words"`
	Status                string  `json:"status"`
	OrderedBy             string  `json:"ordered_by"`
	PickedUp              string  `json:"picked_up"`
	SoldPrice             float64 `json:"sold_price"`
	SoldPriceStr          string  `json:"-"`
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
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if ds == nil {
		vemsgNotExist.Message = fmt.Sprintf("Fail to find diamond with id: %s", c.Param("id"))
		c.JSON(http.StatusOK, vemsgNotExist)
		return
	}
	c.JSON(http.StatusOK, ds)
}

func newDiamond(c *gin.Context) {
	d := diamond{
		ID:                    uuid.NewV4().String(),
		DiamondID:             strings.ToUpper(c.PostForm("diamond_id")),
		StockRef:              strings.ToUpper(c.PostForm("stock_ref")),
		Shape:                 strings.ToUpper(c.PostForm("shape")),
		CaratStr:              strings.ToUpper(c.PostForm("carat")),
		Color:                 strings.ToUpper(c.PostForm("color")),
		Clarity:               strings.ToUpper(c.PostForm("clarity")),
		GradingLab:            strings.ToUpper(c.PostForm("grading_lab")),
		CertificateNumber:     strings.ToUpper(c.PostForm("certificate_number")),
		CutGrade:              strings.ToUpper(c.PostForm("cut_grade")),
		Polish:                strings.ToUpper(c.PostForm("polish")),
		Symmetry:              strings.ToUpper(c.PostForm("symmetry")),
		FluorescenceIntensity: strings.ToUpper(c.PostForm("fluorescence_intensity")),
		Country:               strings.ToUpper(c.PostForm("country")),
		Supplier:              strings.ToUpper(c.PostForm("supplier")),
		PriceNoAddedValueStr:  c.PostForm("price_no_added_value"),
		PriceRetailStr:        c.PostForm("price_retail"),
		Featured:              strings.ToUpper(c.PostForm("featured")),
		RecommandWords:        c.PostForm("recommand_words"),
		ExtraWords:            c.PostForm("extra_words"),
		Status:                strings.ToUpper(c.PostForm("status")),
		Profitable:            strings.ToUpper(c.PostForm("profitable")),
	}
	if vemsg, err := d.validateDiamondReq(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsg) != 0 {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	q := d.composeInsertQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, d.ID)
}

func updateDiamond(c *gin.Context) {
	d := diamond{
		ID:                    c.Param("id"),
		DiamondID:             strings.ToUpper(c.PostForm("diamond_id")),
		StockRef:              strings.ToUpper(c.PostForm("stock_ref")),
		Shape:                 strings.ToUpper(c.PostForm("shape")),
		CaratStr:              strings.ToUpper(c.PostForm("carat")),
		Color:                 strings.ToUpper(c.PostForm("color")),
		Clarity:               strings.ToUpper(c.PostForm("clarity")),
		GradingLab:            strings.ToUpper(c.PostForm("grading_lab")),
		CertificateNumber:     strings.ToUpper(c.PostForm("certificate_number")),
		CutGrade:              strings.ToUpper(c.PostForm("cut_grade")),
		Polish:                strings.ToUpper(c.PostForm("polish")),
		Symmetry:              strings.ToUpper(c.PostForm("symmetry")),
		FluorescenceIntensity: strings.ToUpper(c.PostForm("fluorescence_intensity")),
		Country:               strings.ToUpper(c.PostForm("country")),
		Supplier:              strings.ToUpper(c.PostForm("supplier")),
		PriceNoAddedValueStr:  c.PostForm("price_no_added_value"),
		PriceRetailStr:        c.PostForm("price_retail"),
		Featured:              strings.ToUpper(c.PostForm("featured")),
		RecommandWords:        c.PostForm("recommand_words"),
		ExtraWords:            c.PostForm("extra_words"),
		Status:                strings.ToUpper(c.PostForm("status")),
		PickedUp:              strings.ToUpper(c.PostForm("picked_up")),
		Profitable:            strings.ToUpper(c.PostForm("profitable")),
	}
	if vemsg, err := d.validateDiamondUpdateReq(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsg) != 0 {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	q := d.composeUpdateQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, d.ID)
}

func composeDiamond(rows *sql.Rows) ([]diamond, error) {
	var id, diamondID, stockRef, shape, color, country, supplier, gradingLab string
	var clarity, certificateNumber, cutGrade, polish, symmetry, fluorescenceIntensity string
	var featured, status, profitable string
	var recommandWords, extraWords, orderedBy, pickedUp sql.NullString
	var soldPrice sql.NullFloat64
	var carat, priceNoAddedValue, priceRetail float64

	var ds []diamond
	for rows.Next() {
		if err := rows.Scan(&id, &diamondID, &stockRef, &shape, &carat, &color, &clarity, &gradingLab, &certificateNumber,
			&cutGrade, &polish, &symmetry, &fluorescenceIntensity, &country, &supplier, &priceNoAddedValue, &priceRetail,
			&featured, &recommandWords, &extraWords, &status, &orderedBy, &pickedUp, &soldPrice, &profitable); err != nil {
			return nil, err
		}
		d := diamond{
			ID:                    id,
			DiamondID:             diamondID,
			StockRef:              stockRef,
			Shape:                 shape,
			Carat:                 carat,
			Color:                 color,
			Clarity:               clarity,
			GradingLab:            gradingLab,
			CertificateNumber:     certificateNumber,
			CertificateLink:       composeCertifcateLink(gradingLab, certificateNumber),
			CutGrade:              cutGrade,
			Polish:                polish,
			Symmetry:              symmetry,
			FluorescenceIntensity: fluorescenceIntensity,
			Country:               country,
			Supplier:              supplier,
			PriceNoAddedValue:     priceNoAddedValue,
			PriceRetail:           priceRetail,
			Featured:              featured,
			RecommandWords:        recommandWords.String,
			ExtraWords:            extraWords.String,
			Status:                status,
			OrderedBy:             orderedBy.String,
			PickedUp:              pickedUp.String,
			SoldPrice:             soldPrice.Float64,
			Profitable:            profitable,
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func selectDiamondQuery(id string) string {
	q := `SELECT id, diamond_id, stock_ref, shape, carat, color, clarity, grading_lab, 
	certificate_number, cut_grade, polish, symmetry, fluorescence_intensity, country, 
	supplier, price_no_added_value, price_retail, featured, recommand_words, extra_words, 
	status, ordered_by, picked_up, sold_price, profitable FROM diamonds`

	if id != "" {
		q = fmt.Sprintf("%s WHERE id='%s'", q, id)
	}
	return q
}

//TODO
func processDiamonds(c *gin.Context) {
	id := c.MustGet("id").(string)
	headers := make(map[string]string)
	headers["diamond_id"] = c.PostForm("diamond_id")
	headers["stock_ref"] = c.PostForm("stock_ref")
	headers["shape"] = c.PostForm("shape")
	headers["carat"] = c.PostForm("carat")
	headers["color"] = c.PostForm("color")
	headers["clarity"] = c.PostForm("clarity")
	headers["grading_lab"] = c.PostForm("grading_lab")
	headers["certificate_number"] = c.PostForm("certificate_number")
	headers["cut_grade"] = c.PostForm("cut_grade")
	headers["polish"] = c.PostForm("polish")
	headers["symmetry"] = c.PostForm("symmetry")
	headers["fluorescence_intensity"] = c.PostForm("fluorescence_intensity")
	headers["country"] = c.PostForm("country")
	headers["supplier"] = c.PostForm("supplier")
	headers["price_no_added_value"] = c.PostForm("price_no_added_value")
	headers["price_retail"] = c.PostForm("price_retail")

	vmsg := []string{}
	for k, v := range headers {
		if v == "" {
			vmsg = append(vmsg, k+" has no mapped column\n")
		}
	}
	if len(vmsg) != 0 {
		c.JSON(http.StatusBadRequest, strings.Join(vmsg, ""))
		return
	}

	file := filepath.Join(os.TempDir(), id, c.PostForm("filename"))
	if !util.PathExists(file) {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to find products file :%s", file))
		return
	}

	records, err := util.ParseCSVToArrays(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to parse uploaded file :%s", file))
		return
	}
	ignoredrows, err := importDiamondsCustomizeHeaders(headers, records)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to process uploaded file :%s", file))
		return
	}
	if len(ignoredrows) != 0 {
		// var msg string
		// for i := 0; i < len(ignoredrows); i++ {
		// 	msg := msg + strings.Join(ignoredrows[i], ",") + "\n"
		// }
		c.JSON(http.StatusOK, gin.H{"IngoredRows": ignoredrows})
		return
	}
	c.JSON(http.StatusOK, "success processed uploaded file!")
}
