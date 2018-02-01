package main

import (
	"fmt"
	"strings"
	"util"
)

func customerLevelByAmount(amount float64) string {
	if amount < 5000 {
		return LEVEL0
	} else if amount >= 5000 && amount < 10000 {
		return LEVEL1
	} else if amount >= 5000 && amount < 10000 {
		return LEVEL2
	}
	return LEVEL3
}

func priceForCustomer(level string, price float64) float64 {
	switch level {
	case LEVEL0:
		return price * 0.99
	case LEVEL1:
		return price * 0.98
	case LEVEL2:
		return price * 0.95
	case LEVEL3:
		return price * 0.93
	case LEVEL4:
		return price * 0.90
	case LEVEL5:
		return price * 0.85
	case LEVEL6:
		return price * 0.90
	default:
		return price
	}
}

func priceForAgent(level string, price float64) float64 {
	switch level {
	case LEVEL0:
		return price
	case LEVEL1:
		return price * 0.9
	case LEVEL2:
		return price * 0.85
	case LEVEL3:
		return price * 0.83
	default:
		return price
	}
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
