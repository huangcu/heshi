package util

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
