package main

import (
	"fmt"
	"strings"
	"util"
)

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
