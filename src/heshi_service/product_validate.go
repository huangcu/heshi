package main

import (
	"bytes"
	"fmt"
	"strings"
)

type ValidationErrors map[string]string

func (ve ValidationErrors) Error() string {
	var buff bytes.Buffer
	for key, msg := range ve {
		buff.WriteString(fmt.Sprintf("%s %s\n", key, msg))
	}
	return strings.TrimSpace(buff.String())
}

func (d *diamond) validateDiamondReq() string {
	// govalidator.StringLength()
	return ""
}
func (d *jewelry) validateJewelryReq() string {
	// govalidator.StringLength()
	return ""
}
func (d *smallDiamond) validateSmallDiamondReq() string {
	// govalidator.StringLength()
	return ""
}
