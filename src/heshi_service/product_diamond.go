package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"
	"os"
	"path/filepath"
	"sql_patch"
	"strings"
	"time"
	"util"

	"github.com/gin-gonic/gin"
)

type diamond struct {
	ID                    string   `json:"id"`
	DiamondID             string   `json:"diamond_id"`
	StockRef              string   `json:"stock_ref"`
	Shape                 string   `json:"shape"`
	Carat                 float64  `json:"carat"`
	CaratStr              string   `json:"-"`
	Color                 string   `json:"color"`
	Clarity               string   `json:"clarity"`
	GradingLab            string   `json:"grading_lab"`
	CertificateNumber     string   `json:"certificate_number"`
	CutGrade              string   `json:"cut_grade"`
	Polish                string   `json:"polish"`
	Symmetry              string   `json:"symmetry"`
	FluorescenceIntensity string   `json:"fluorescence_intensity"`
	Country               string   `json:"country"`
	Supplier              string   `json:"supplier"`
	PriceNoAddedValue     float64  `json:"price_no_added_value"`
	PriceNoAddedValueStr  string   `json:"-"`
	PriceRetail           float64  `json:"price_retail"`
	PriceRetailStr        string   `json:"-"`
	CertificateLink       string   `json:"certificate_link"`
	Featured              string   `json:"featured"`
	RecommendWords        string   `json:"recommend_words"`
	ExtraWords            string   `json:"extra_words"`
	Images                []string `json:"images"`
	Status                string   `json:"status"`
	OrderedBy             string   `json:"ordered_by"`
	PickedUp              string   `json:"picked_up"`
	SoldPrice             float64  `json:"sold_price"`
	SoldPriceStr          string   `json:"-"`
	Profitable            string   `json:"profitable"`
	promotion
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
	imageFileNames, vemsg, err := validateUploadedMultipleFile(c, "diamond", "image", int64(imageSizeLimit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if vemsg != (errors.HSMessage{}) {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	d := diamond{
		ID:                    newV4(),
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
		RecommendWords:        c.PostForm("recommend_words"),
		ExtraWords:            c.PostForm("extra_words"),
		Images:                imageFileNames,
		Status:                strings.ToUpper(c.PostForm("status")),
		Profitable:            strings.ToUpper(c.PostForm("profitable")),
	}
	if vemsg, err := d.validateDiamondReq(false); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsg) != 0 {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	if err := saveUploadedMultipleFile(c, "diamond", "image", imageFileNames); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
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
	uid := c.MustGet("id").(string)
	did := c.Param("id")
	if exist, err := isDiamondExistByID(did); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if !exist {
		c.JSON(http.StatusBadRequest, "Item doesn't exist")
		return
	}
	imageFileNames, vemsg, err := validateUploadedMultipleFile(c, "diamond", "image", int64(imageSizeLimit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if vemsg != (errors.HSMessage{}) {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	d := diamond{
		ID:                    did,
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
		RecommendWords:        c.PostForm("recommend_words"),
		ExtraWords:            c.PostForm("extra_words"),
		Status:                strings.ToUpper(c.PostForm("status")),
		Images:                imageFileNames,
		PickedUp:              strings.ToUpper(c.PostForm("picked_up")),
		Profitable:            strings.ToUpper(c.PostForm("profitable")),
	}
	if vemsg, err := d.validateDiamondReq(true); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsg) != 0 {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	if err := saveUploadedMultipleFile(c, "diamond", "image", imageFileNames); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	q := d.composeUpdateQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, d.ID)
	go newHistoryRecords(uid, "diamonds", d.ID, d.parmsKV())
}

func composeDiamond(rows *sql.Rows) ([]diamond, error) {
	var id, diamondID, stockRef, shape, color, country, supplier, gradingLab string
	var clarity, certificateNumber, cutGrade, polish, symmetry, fluorescenceIntensity string
	var featured, status, profitable string
	var images, recommendWords, extraWords, orderedBy, pickedUp sql.NullString
	var soldPrice sql.NullFloat64
	var carat, priceNoAddedValue, priceRetail float64

	var pid, promType, pstatus sql.NullString
	var promPrice sql.NullFloat64
	var promDiscount sql.NullInt64
	var beginAt, endAt sql_patch.NullTime
	var ds []diamond
	for rows.Next() {
		if err := rows.Scan(&id, &diamondID, &stockRef, &shape, &carat, &color, &clarity, &gradingLab, &certificateNumber,
			&cutGrade, &polish, &symmetry, &fluorescenceIntensity, &country, &supplier, &priceNoAddedValue, &priceRetail,
			&featured, &recommendWords, &extraWords, &images, &status, &orderedBy, &pickedUp, &soldPrice, &profitable,
			&pid, &promType, &promDiscount, &promPrice, &beginAt, &endAt, &pstatus); err != nil {
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
			RecommendWords:        recommendWords.String,
			ExtraWords:            extraWords.String,
			Status:                status,
			OrderedBy:             orderedBy.String,
			PickedUp:              pickedUp.String,
			SoldPrice:             soldPrice.Float64,
			Profitable:            profitable,
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

func selectDiamondQuery(id string) string {
	q := `SELECT diamonds.id, diamond_id, stock_ref, shape, carat, color, clarity, grading_lab, 
	certificate_number, cut_grade, polish, symmetry, fluorescence_intensity, country, 
	supplier, price_no_added_value, price_retail, featured, recommend_words, extra_words, images,
	diamonds.status, ordered_by, picked_up, sold_price, profitable, 
	promotions.id, prom_type, prom_discount, prom_price, begin_at, end_at, promotions.status 
	FROM diamonds 
	LEFT JOIN promotions ON diamonds.promotion_id=promotions.id 
	WHERE diamonds.status IN ('AVAILABLE','OFFLINE')`

	if id != "" {
		q = fmt.Sprintf("%s AND id='%s'", q, id)
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
