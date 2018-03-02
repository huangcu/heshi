package main

import (
	"fmt"
	"strings"
	"util"

	uuid "github.com/satori/go.uuid"
)

func downPayment(price float64) float64 {
	return floatToFixed2(price * 0.3)
}

func customerLevelByAmount(amount float64) (string, error) {
	q := fmt.Sprintf("SELECT level FROM configs WHERE type='%s' AND amount < '%f' order by level DESC",
		CUSTOMER, amount)
	rows, err := dbQuery(q)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		var level string
		if err := rows.Scan(&level); err != nil {
			return level, nil
		}
	}
	return LEVEL1, nil
	// if amount < 5000 {
	// 	return LEVEL0
	// } else if amount >= 5000 && amount < 10000 {
	// 	return LEVEL1
	// } else if amount >= 10000 && amount < 20000 {
	// 	return LEVEL2
	// }
	// return LEVEL3
}

func agentLevelByAmountAndPieces(amount float64, pieces int) (string, error) {
	q := fmt.Sprintf(`SELECT level FROM configs 
		WHERE type='%s' 
		AND (amount < '%f' OR pieces < '%d') order by level DESC`,
		AGENT, amount, pieces)
	rows, err := dbQuery(q)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		var level string
		if err := rows.Scan(&level); err != nil {
			return level, nil
		}
	}
	return LEVEL1, nil
}

func priceForCustomer(level string, price float64) (float64, error) {
	q := fmt.Sprintf("SELECT discount FROM configs WHERE type='%s' AND level = '%s' limit 1",
		CUSTOMER, level)
	var discount float64
	if err := dbQueryRow(q).Scan(&discount); err != nil {
		return 0, err
	}
	return floatToFixed2(price * discount), nil
	// switch level {
	// case LEVEL0:
	// 	return floatToFixed2(price * 0.99)
	// case LEVEL1:
	// 	return floatToFixed2(price * 0.98)
	// case LEVEL2:
	// 	return floatToFixed2(price * 0.95)
	// case LEVEL3:
	// 	return floatToFixed2(price * 0.93)
	// case LEVEL4:
	// 	return floatToFixed2(price * 0.90)
	// case LEVEL5:
	// 	return floatToFixed2(price * 0.85)
	// case LEVEL6:
	// 	return floatToFixed2(price * 0.90)
	// default:
	// 	return price
	// }
}

//TODO agentLevelDiscount
func priceForAgent(level string, price float64) (float64, error) {
	q := fmt.Sprintf("SELECT discount FROM configs WHERE type='%s' AND level = '%s' limit 1",
		AGENT, level)
	var discount float64
	if err := dbQueryRow(q).Scan(&discount); err != nil {
		return 0, err
	}
	return floatToFixed2(price * discount), nil
	// switch level {
	// case LEVEL0:
	// 	return price
	// case LEVEL1:
	// 	return floatToFixed2(price * 0.9)
	// case LEVEL2:
	// 	return floatToFixed2(price * 0.85)
	// case LEVEL3:
	// 	return floatToFixed2(price * 0.83)
	// default:
	// 	return price
	// }
}

func euroToYuan(priceEuro float64) float64 {
	return (priceEuro / activeCurrencyRate.Rates.EUR) * activeCurrencyRate.Rates.CNY
}

func euroToDollar(priceEuro float64) float64 {
	return priceEuro / activeCurrencyRate.Rates.EUR
}

func yuanToEuro(priceYuan float64) float64 {
	return (priceYuan / activeCurrencyRate.Rates.CNY) * activeCurrencyRate.Rates.EUR
}

func yuanToDollar(priceYuan float64) float64 {
	return priceYuan / activeCurrencyRate.Rates.CNY
}

func dollarToEuro(priceDollar float64) float64 {
	return priceDollar * activeCurrencyRate.Rates.EUR
}

func dollarToYuan(priceDollar float64) float64 {
	return priceDollar * activeCurrencyRate.Rates.CNY
}

func isItemExistInDbByProperty(dbName, property, propertyValue string) (bool, error) {
	var count int
	q := fmt.Sprintf("SELECT count(*) FROM %s WHERE %s='%s'", dbName, property, propertyValue)
	if err := dbQueryRow(q).Scan(&count); err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func isItemExistInDbByPropertyWithDifferentID(dbName, property, propertyValue, id string) (bool, error) {
	var count int
	q := fmt.Sprintf("SELECT count(*) FROM %s WHERE %s='%s' AND id!='%s'", dbName, property, propertyValue, id)
	if err := dbQueryRow(q).Scan(&count); err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func ItemsNotInArray(item string, items []string) []string {
	itemStr := FormatInputString(item)
	var notIn []string
	for _, v := range strings.Split(itemStr, ",") {
		if !util.IsInArrayString(v, items) {
			notIn = append(notIn, v)
		}
	}
	return notIn
}

func FormatInputString(input string) string {
	return strings.ToUpper(strings.Replace(input, " ", "", -1))
}

func floatToFixed2(v float64) float64 {
	return float64(int(v*100)) / 100
}

func newV4() string {
	v4, _ := uuid.NewV4()
	return v4.String()
}
