package util

import (
	"strconv"
	"time"
)

func ParseStringTimestamp(duration string) time.Time {
	//duration:="1512799203"
	i, err := strconv.ParseInt(duration, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm
}
