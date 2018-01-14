package main

import (
	"database/sql"
	"encoding/json"
	"heshi/errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

// Free, Reliable Currency Converter API（免费账户，一天更新一次，一个月免费1000次，免费账户不支持https）
// https://currencylayer.com/

// Forex API: Realtime Forex Quotes（免费账户，每分钟更新一次，一天免费1000次，支持https）
// https://1forge.com/forex-data-api

// Open Exchange Rates(每月1000次免费调用；收费用户最高可小时更新) 一般来说，法币汇率比较稳定，可以一天获取一次数据，如果你的应用需要依赖高实时性的，那就花些钱，或者是去外汇交易所拉去实时数据。
// https://openexchangerates.org
// email: huang372923@hotmail.com
// password: hxxxxxxxxxx?1
// app_id: 16358f5c12884da89c2ae2c767d183d5
// https://openexchangerates.org/api/latest.json?app_id=16358f5c12884da89c2ae2c767d183d5
// {
//   "disclaimer": "Usage subject to terms: https://openexchangerates.org/terms",
//   "license": "https://openexchangerates.org/license",
//   "timestamp": 1512799203,
//   "base": "USD",
//   "rates": {
//     "AED": 3.673097,
//     "AFN": 68.693,
//     "ALL": 113.595865,
//     "AMD": 483.71,
//     "ANG": 1.785348,
//     "AOA": 165.9235,
//     "ARS": 17.256,
//     "AUD": 1.33209,
//     "AWG": 1.786821,
//     "AZN": 1.7,
//     "BAM": 1.661116,
//     "BBD": 2,
//     "BDT": 82.163042,
//     "BGN": 1.6635,
//     "BHD": 0.376759,
//     "BIF": 1770,
//     "BMD": 1,
//     "BND": 1.35309,
//     "BOB": 6.888858,
//     "BRL": 3.2916,
//     "BSD": 1,
//     "BTC": 0.000064435797,
//     "BTN": 64.423471,
//     "BWP": 10.320053,
//     "BYN": 2.02926,
//     "BZD": 2.010404,
//     "CAD": 1.284855,
//     "CDF": 1596,
//     "CHF": 0.99255,
//     "CLF": 0.02424,
//     "CLP": 655.6,
//     "CNH": 6.624005,
//     "CNY": 6.615946,
//     "COP": 3002,
//     "CRC": 565.92,
//     "CUC": 1,
//     "CUP": 25.5,
//     "CVE": 93.75,
//     "CZK": 21.71135,
//     "DJF": 178.57,
//     "DKK": 6.321895,
//     "DOP": 48.143492,
//     "DZD": 115.448,
//     "EGP": 17.8,
//     "ERN": 15.228279,
//     "ETB": 27.095551,
//     "EUR": 0.850051,
//     "FJD": 2.077803,
//     "FKP": 0.747105,
//     "GBP": 0.747105,
//     "GEL": 2.688101,
//     "GGP": 0.747105,
//     "GHS": 4.486158,
//     "GIP": 0.747105,
//     "GMD": 47.5,
//     "GNF": 9000,
//     "GTQ": 7.35168,
//     "GYD": 207.982754,
//     "HKD": 7.80517,
//     "HNL": 23.572307,
//     "HRK": 6.4163,
//     "HTG": 63.551356,
//     "HUF": 266.7245,
//     "IDR": 13544.307977,
//     "ILS": 3.52521,
//     "IMP": 0.747105,
//     "INR": 64.4866,
//     "IQD": 1188.4,
//     "IRR": 35094.293807,
//     "ISK": 104.549272,
//     "JEP": 0.747105,
//     "JMD": 125.019767,
//     "JOD": 0.709002,
//     "JPY": 113.475,
//     "KES": 103.11792,
//     "KGS": 69.735896,
//     "KHR": 4025.6,
//     "KMF": 418.45,
//     "KPW": 900,
//     "KRW": 1092.2,
//     "KWD": 0.301652,
//     "KYD": 0.833541,
//     "KZT": 333.971541,
//     "LAK": 8300.5,
//     "LBP": 1511.35,
//     "LKR": 152.73096,
//     "LRD": 125.5,
//     "LSL": 13.578701,
//     "LYD": 1.363406,
//     "MAD": 9.468123,
//     "MDL": 17.246807,
//     "MGA": 3199.65,
//     "MKD": 52.30053,
//     "MMK": 1354.199087,
//     "MNT": 2436.956458,
//     "MOP": 8.042383,
//     "MRO": 354.138571,
//     "MUR": 33.982,
//     "MVR": 15.409873,
//     "MWK": 725.54,
//     "MXN": 18.923425,
//     "MYR": 4.081134,
//     "MZN": 60.003288,
//     "NAD": 13.578701,
//     "NGN": 358.9,
//     "NIO": 30.738893,
//     "NOK": 8.300822,
//     "NPR": 103.059806,
//     "NZD": 1.46185,
//     "OMR": 0.384963,
//     "PAB": 1,
//     "PEN": 3.236,
//     "PGK": 3.207749,
//     "PHP": 50.475,
//     "PKR": 105.9,
//     "PLN": 3.56611,
//     "PYG": 5653.2,
//     "QAR": 3.700015,
//     "RON": 3.937483,
//     "RSD": 101.610561,
//     "RUB": 59.0638,
//     "RWF": 855.255338,
//     "SAR": 3.7496,
//     "SBD": 7.768106,
//     "SCR": 14.000747,
//     "SDG": 6.675,
//     "SEK": 8.44565,
//     "SGD": 1.35225,
//     "SHP": 0.747105,
//     "SLL": 7660.273512,
//     "SOS": 576.641453,
//     "SRD": 7.448,
//     "SSP": 130.2634,
//     "STD": 20778.698538,
//     "SVC": 8.752125,
//     "SYP": 514.96999,
//     "SZL": 13.578004,
//     "THB": 32.5975,
//     "TJS": 8.801761,
//     "TMT": 3.499986,
//     "TND": 2.496507,
//     "TOP": 2.30228,
//     "TRY": 3.838831,
//     "TTD": 6.709363,
//     "TWD": 30.008335,
//     "TZS": 2241.1,
//     "UAH": 27.116,
//     "UGX": 3603.9,
//     "USD": 1,
//     "UYU": 29.074238,
//     "UZS": 8071.5,
//     "VEF": 10.72835,
//     "VND": 22717.035659,
//     "VUV": 107.938209,
//     "WST": 2.551251,
//     "XAF": 557.59691,
//     "XAG": 0.06310745,
//     "XAU": 0.00080115,
//     "XCD": 2.70255,
//     "XDR": 0.708986,
//     "XOF": 557.59691,
//     "XPD": 0.00099308,
//     "XPF": 101.438068,
//     "XPT": 0.00112617,
//     "YER": 250.325,
//     "ZAR": 13.663355,
//     "ZMW": 10.243477,
//     "ZWL": 322.355011
//   }
// }

type currency struct {
	ID        string   `json:"-"`
	Note      string   `json:"note"`
	Base      string   `json:"base"`
	Timestamp int64    `json:"timestamp"`
	Rates     Rate     `json:"rates"`
	RatesFluc RateFluc `json:"-"`
}

type Rate struct {
	USD float64 `json:"USD"` //USD
	CNY float64 `json:"CNY"` //RMB
	EUR float64 `json:"EUR"` //Euro
	CAD float64 `json:"CAD"` //Canadian Dollar
	AUD float64 `json:"AUD"` //Australian Dollar
	CHF float64 `json:"CHF"` //Swiss Franc
	RUB float64 `json:"RUB"` //Russian Ruble
	NZD float64 `json:"NZD"` //New Zealand Dollar
}

type RateFluc struct {
	USDFluc float64 `json:"USD_Fluc"` //USD
	CNYFluc float64 `json:"CNY_Fluc"` //RMB
	EURFluc float64 `json:"EUR_Fluc"` //Euro
	CADFluc float64 `json:"CAD_Fluc"` //Canadian Dollar
	AUDFluc float64 `json:"AUD_Fluc"` //Australian Dollar
	CHFFluc float64 `json:"CHF_Fluc"` //Swiss Franc
	RUBFluc float64 `json:"RUB_Fluc"` //Russian Ruble
	NZDFluc float64 `json:"NZD_Fluc"` //New Zealand Dollar
}

var exchangeRateURL = "https://openexchangerates.org/api/latest.json?app_id=16358f5c12884da89c2ae2c767d183d5"

func getLatestRates() error {
	req, err := http.NewRequest("GET", exchangeRateURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Newf("fail to get exchange rate from %s", exchangeRateURL)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var c currency
	if err := json.Unmarshal(b, &c); err != nil {
		return err
	}
	c.Note = "FROM Open Exchange Rates API"
	c.ID = uuid.NewV4().String()

	q := c.composeInsertQuery()
	_, err = dbExec(q)
	return err
}

func getCurrencyRate(c *gin.Context) {
	currencyRate, err := getAcitveCurrencyRate()
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, VEMSG_EXCHANGE_RATE_NOT_EXIST)
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, currencyRate)
}

func newCurrencyRate(c *gin.Context) {
	currencyRate := c.MustGet("currency").(*currency)
	q := currencyRate.composeInsertQuery()
	if _, err := db.Exec(q); err != nil {
		c.String(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}

	c.String(http.StatusOK, currencyRate.ID)
}

func getAcitveCurrencyRate() (*currency, error) {
	var base, note string
	var usd, cny, eur, cad, aud, chf, rub, nzd float64
	var createdAt time.Time
	q := `SELECT base,note,usd,cny,eur,cad,aud,chf,rub,nzd,created_at FROM currency_exchange_rates ORDER BY created_at DESC LIMIT 1`
	if err := db.QueryRow(q).Scan(&base, &note, &usd, &cny, &eur, &cad, &aud, &chf, &rub, &nzd, &createdAt); err != nil {
		return nil, err
	}
	currencyRate := &currency{
		Base:      base,
		Note:      note,
		Timestamp: createdAt.Unix(),
		Rates: Rate{
			USD: usd,
			CNY: cny,
			EUR: eur,
			CAD: cad,
			AUD: aud,
			CHF: chf,
			RUB: rub,
			NZD: nzd,
		},
	}
	return currencyRate, nil
}
