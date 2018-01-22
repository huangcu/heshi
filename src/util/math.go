package util

import (
	"math"
	"strconv"
	"strings"
)

func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func StringToFloat(v string) (float64, error) {
	cValue, err := strconv.ParseFloat(strings.Replace(v, ",", "", -1), 64)
	if err != nil {
		return 0, err
	}
	return math.Abs(cValue), nil
}
