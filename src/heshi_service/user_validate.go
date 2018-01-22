package main

import (
	"fmt"
	"heshi/errors"
	"regexp"
	"util"

	"github.com/asaskevich/govalidator"
)

func (u *User) validNewUser() ([]errors.HSMessage, error) {
	if vmsg := u.preValidateNewUser(); len(vmsg) != 0 {
		return vmsg, nil
	} else if vmsg, err := u.validUniqueKey(); err != nil {
		return nil, err
	} else if len(vmsg) != 0 {
		return vmsg, nil
	}
	return nil, nil
}

/*preValidateNewUser validate new user Params*/
func (u *User) preValidateNewUser() []errors.HSMessage {
	vmsg := u.requiredField()
	if len(vmsg) != 0 {
		return vmsg
	}
	var vmsgs []errors.HSMessage
	if vmsg = u.validUserName(); len(vmsg) != 0 {
		vmsgs = append(vmsgs, vmsg...)
	}
	if u.Email != "" {
		if !govalidator.IsEmail(u.Email) {
			vmsgs = append(vmsgs, vemsgUserEmailNotValid)
		}
	}

	if vmsg := u.validPhone(); vmsg != (errors.HSMessage{}) {
		vmsgs = append(vmsgs, vmsg)
	}

	if !util.IsInArrayString(u.UserType, VALID_USERTYPE) {
		vmsgs = append(vmsgs, vemsgUserUsertypeNotValid)
	}

	if vmsg := u.validRecommnadedBy(); vmsg != (errors.HSMessage{}) {
		vmsgs = append(vmsgs, vmsg)
	}

	if len(vmsgs) != 0 {
		return vmsg
	}
	return nil
}

func (u *User) requiredField() []errors.HSMessage {
	var vmsg []errors.HSMessage
	if u.Cellphone == "" && u.Email == "" {
		vmsg = append(vmsg, vemsgUserCellphoneEmailEmpty)
	}
	if u.Password == "" {
		vmsg = append(vmsg, vemsgUserPasswordEmpty)
	}
	return vmsg
}

func (u *User) validUserName() []errors.HSMessage {
	var vmsg []errors.HSMessage
	if u.Username != "" {
		if !govalidator.StringLength(u.Username, "6", "40") {
			vmsg = append(vmsg, vemsgUserUsernameError1)
		}
		if !govalidator.Matches(u.Username, "^[a-zA-Z0-9]*$") {
			vmsg = append(vmsg, vemsgUserUsernameError2)
		}
	}
	return vmsg
}

// TODO refine phone number valdiation
func (u *User) validPhone() errors.HSMessage {
	// regex := regexp.MustCompile("^(\\+\\d{1,3}[- ]?)?\\d{14}$")
	// regex := regexp.MustCompile("^\\+{0,1}0{0,1}62[0-9]+$")
	regex := regexp.MustCompile("^[0-9]*$")
	if !regex.MatchString(u.Cellphone) {
		return vemsgUserCellphoneNotValid
	}
	return errors.HSMessage{}
}

func (u *User) validRecommnadedBy() errors.HSMessage {
	if u.RecommendedBy != "" {
		var count int
		q := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE invitation_code=%s", u.RecommendedBy)
		if err := dbQueryRow(q).Scan(&count); err != nil {
			return vemsgUserErrorRecommandCode
		}
		if count == 0 {
			return vemsgUserErrorRecommandCode
		}
	}

	return errors.HSMessage{}
}

func (u *User) validUniqueKey() ([]errors.HSMessage, error) {
	var vemsgs []errors.HSMessage
	if exist, err := u.isUserExistByUserName(); err != nil {
		return nil, nil
	} else if exist {
		vemsgs = append(vemsgs, vemsgUserUsernameDuplicate)
	}
	if exist, err := u.isUserExistByCellphone(); err != nil {
		return nil, nil
	} else if exist {
		vemsgs = append(vemsgs, vemsgUserCellphoneDuplicate)
	}

	if exist, err := u.isUserExistByEmail(); err != nil {
		return nil, nil
	} else if exist {
		vemsgs = append(vemsgs, vemsgUserEmailDuplicate)
	}

	return vemsgs, nil
}

func (u *User) isUserExistByUserName() (bool, error) {
	if u.Username == "" {

		return false, nil
	}
	var count int
	q := fmt.Sprintf("SELECT count(*) FROM users WHERE username='%s'", u.Username)
	if err := dbQueryRow(q).Scan(&count); err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (u *User) isUserExistByCellphone() (bool, error) {
	var count int
	q := fmt.Sprintf("SELECT count(*) FROM users WHERE cellphone='%s'", u.Cellphone)
	if err := dbQueryRow(q).Scan(&count); err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (u *User) isUserExistByEmail() (bool, error) {
	var count int
	q := fmt.Sprintf("SELECT count(*) FROM users WHERE email='%s'", u.Email)
	if err := dbQueryRow(q).Scan(&count); err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
