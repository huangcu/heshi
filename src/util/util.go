package util

import (
	"math/rand"
	"strings"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func IsInArrayString(item string, items []string) bool {
	inFlag := false
	for _, v := range items {
		if v == item {
			inFlag = true
			break
		}
	}
	return inFlag
}

func IsInArrayInt(item int, items []int) bool {
	inFlag := false
	for _, v := range items {
		if v == item {
			inFlag = true
			break
		}
	}
	return inFlag
}

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func IsIn(items1 []string, items2 []string) bool {
	for _, v1 := range items1 {
		inFlag := false
		for _, v2 := range items2 {
			if v1 == v2 {
				inFlag = true
				break
			}
		}
		if !inFlag {
			return false
		}
	}
	return true
}

func ItemsNotInArray(item string, items []string) []string {
	itemStr := FormatInputString(item)
	var notIn []string
	for _, v := range strings.Split(itemStr, ",") {
		if !IsInArrayString(v, items) {
			notIn = append(notIn, v)
		}
	}
	return notIn
}

func FormatInputString(input string) string {
	return strings.ToUpper(strings.Replace(input, " ", "", -1))
}
