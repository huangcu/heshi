package main

import (
	"bytes"
	"fmt"
	"strings"
	"util"

	"github.com/asaskevich/govalidator"
)

type ValidationErrors map[string]string

func (ve ValidationErrors) Error() string {
	var buff bytes.Buffer
	for key, msg := range ve {
		buff.WriteString(fmt.Sprintf("%s %s\n", key, msg))
	}
	return strings.TrimSpace(buff.String())
}

/*preValidateNewUser validate new user Params*/
func preValidateNewUser(nu User) string {
	var vmsg []string
	if _, err := govalidator.ValidateStruct(nu); err != nil {
		vmsg = append(vmsg, err.Error())
	}
	if nu.Cellphone == "" && nu.Email == "" {
		vmsg = append(vmsg, VEMSG_CELLPHONE_EMAIL)
	}
	if !util.IsInArrayString(nu.UserType, VALID_USERTYPE) {
		vmsg = append(vmsg, VEMSG_VALUE_NOT_VALID)
	}

	if len(vmsg) != 0 {
		return strings.Join(vmsg, ";")
	}
	return ""
}
