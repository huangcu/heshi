package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

func (c *currency) composeInsertQuery() string {
	params := c.parmsKV()
	q := `INSERT INTO currency_exchange_rates (id`
	va := fmt.Sprintf(`VALUES ('%s'`, c.ID)
	for k, v := range params {
		q = fmt.Sprintf("%s, %s", q, k)
		switch v.(type) {
		case string:
			va = fmt.Sprintf("%s, '%s'", va, v.(string))
		case float64:
			va = fmt.Sprintf("%s, '%f'", va, v.(float64))
		case int:
			va = fmt.Sprintf("%s, '%d'", va, v.(int))
		case int64:
			va = fmt.Sprintf("%s, '%d'", va, v.(int64))
		}
	}
	q = fmt.Sprintf("%s) %s)", q, va)
	return q
}

// 	params := make(map[string]interface{})
func (c *currency) parmsKV() map[string]interface{} {
	params := make(map[string]interface{})
	params["base"] = c.Base
	params["note"] = c.Note
	params["usd"] = c.Rates.USD
	params["cny"] = c.Rates.CNY
	params["eur"] = c.Rates.EUR
	params["cad"] = c.Rates.CAD
	params["aud"] = c.Rates.AUD
	params["chf"] = c.Rates.CHF
	params["nzd"] = c.Rates.NZD
	params["rub"] = c.Rates.RUB
	return params
}

func currencyRateReqValidator(h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// if !util.IsInArrayString(c.PostForm("base"), VALID_CURRENCY_SYMBOL) {
		// 	c.String(http.StatusOK, VEMSG_CURRENCY_NOT_VALID_SYMBOL)
		// 	return
		// }
		//NOTES: so far currency exchange rate base is USD.
		if c.PostForm("base") != "USD" {
			c.JSON(http.StatusOK, vemsgCurrencyBaseNotValid)
			return
		}
		//TODO
		cny, err := strconv.ParseFloat(c.PostForm("cny"), 6)
		if err != nil {
			c.JSON(http.StatusOK, vemsgCurrencyRateCNYNotValid)
			return
		}
		eur, err := strconv.ParseFloat(c.PostForm("eur"), 6)
		if err != nil {
			c.JSON(http.StatusOK, vemsgCurrencyRateEURNotValid)
			return
		}
		//NOTES: so far only need cny & eur
		cad, aud, chf, rub, nzd := 0.0, 0.0, 0.0, 0.0, 0.0
		// aud, err := strconv.ParseFloat(c.PostForm("aud"), 6)
		// if err != nil {
		// 	c.String(http.StatusBadRequest, err.Error())
		// 	return
		// }
		// cad, err := strconv.ParseFloat(c.PostForm("cad"), 6)
		// if err != nil {
		// 	c.String(http.StatusBadRequest, err.Error())
		// 	return
		// }
		// chf, err := strconv.ParseFloat(c.PostForm("chf"), 6)
		// if err != nil {
		// 	c.String(http.StatusBadRequest, err.Error())
		// 	return
		// }
		// rub, err := strconv.ParseFloat(c.PostForm("rub"), 6)
		// if err != nil {
		// 	c.String(http.StatusBadRequest, err.Error())
		// 	return
		// }
		// nzd, err := strconv.ParseFloat(c.PostForm("nzd"), 6)
		// if err != nil {
		// 	c.String(http.StatusBadRequest, err.Error())
		// 	return
		// }

		currencyRate := currency{
			ID:   uuid.NewV4().String(),
			Base: c.PostForm("base"),
			Note: "manual input",
			Rates: Rate{
				USD: 1,
				CNY: cny,
				EUR: eur,
				CAD: cad,
				AUD: aud,
				CHF: chf,
				RUB: rub,
				NZD: nzd,
			},
		}
		c.Set("currency", currencyRate)
		h(c)
	}
}
