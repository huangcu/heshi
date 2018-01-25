package main

import (
	"fmt"
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
