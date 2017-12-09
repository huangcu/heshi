package main

import (
	"fmt"
	"regexp"
	"strings"
	"util"

	"github.com/asaskevich/govalidator"
)

/*preValidateNewUser validate new user Params*/
func (u *User) preValidateNewUser() string {
	vmsg := u.requiredField()
	if vmsg != "" {
		return vmsg
	}
	var vmsgs []string
	if vmsg = u.validUserName(); vmsg != "" {
		vmsgs = append(vmsgs, vmsg)
	}
	if u.Email != "" {
		if govalidator.IsEmail(u.Email) {
			vmsgs = append(vmsgs, VEMSG_USER_EMAIL_NOT_VALID)
		}
	}

	if vmsg = u.validPhone(); vmsg != "" {
		vmsgs = append(vmsgs, vmsg)
	}

	if !util.IsInArrayString(u.UserType, VALID_USERTYPE) {
		vmsgs = append(vmsgs, VEMSG_USER_USERTYPE_NOT_VALID)
	}

	if vmsg = u.validRecommnadedBy(); vmsg != "" {
		vmsgs = append(vmsgs, vmsg)
	}

	if len(vmsgs) != 0 {
		return strings.Join(vmsgs, " ")
	}
	return ""
}

func (u *User) requiredField() string {
	var vmsg string
	if u.Cellphone == "" && u.Email == "" {
		vmsg = VEMSG_USER_CELLPHONE_EMAIL_EMPTY
	}
	if u.Password == "" {
		vmsg = vmsg + VEMSG_USER_PASSWORD_EMPTY
	}
	if u.UserType == "" {
		vmsg = vmsg + VEMSG_USER_USERTYPE_EMPTY
	}
	return vmsg
}

func (u *User) validUserName() string {
	var vmsg string
	if u.Username != "" {
		if !govalidator.StringLength(u.Username, "6", "40") {
			vmsg = VEMSG_USER_USERNAME_ERROR1
		}
		if !govalidator.Matches(u.Username, "^[a-zA-Z0-9]*$") {
			vmsg = vmsg + VEMSG_USER_USERNAME_ERROR2
		}
	}
	return vmsg
}

// TODO refine phone number valdiation
func (u *User) validPhone() string {
	// regex := regexp.MustCompile("^(\\+\\d{1,3}[- ]?)?\\d{14}$")
	// regex := regexp.MustCompile("^\\+{0,1}0{0,1}62[0-9]+$")
	regex := regexp.MustCompile("^[0-9]*$")
	if !regex.MatchString(u.Cellphone) {
		return VEMSG_USER_CELLPHONE_NOT_VALID
	}
	return ""
}

func (u *User) validRecommnadedBy() string {
	if u.RecommendedBy != "" {
		var count int
		q := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE invitation_code=%s", u.RecommendedBy)
		if err := dbQueryRow(q).Scan(&count); err != nil {
			return VEMSG_USER_ERROR_RECOMMAND_CODE
		}
		if count == 0 {
			return VEMSG_USER_ERROR_RECOMMAND_CODE
		}
	}

	return ""
}
