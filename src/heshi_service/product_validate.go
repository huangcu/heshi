package main

import (
	"bytes"
	"fmt"
	"strings"
)

type validationErrors map[string]string

func (ve validationErrors) Error() string {
	var buff bytes.Buffer
	for key, msg := range ve {
		buff.WriteString(fmt.Sprintf("%s %s\n", key, msg))
	}
	return strings.TrimSpace(buff.String())
}

func (d *smallDiamond) validateSmallDiamondReq() string {
	// govalidator.StringLength()
	return ""
}
