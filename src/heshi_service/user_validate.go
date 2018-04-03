package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"regexp"
	"util"

	"github.com/asaskevich/govalidator"
)

func (u *User) validNewUser() ([]errors.HSMessage, error) {
	if vmsg := u.preValidateNewUser(); len(vmsg) != 0 {
		return vmsg, nil
	}

	if vmsg, err := u.validUniqueKey(); err != nil {
		return nil, err
	} else if len(vmsg) != 0 {
		return vmsg, nil
	}

	if msg, err := u.validRecommendedBy(); err != nil {
		return nil, err
	} else if msg != (errors.HSMessage{}) {
		var vmsg []errors.HSMessage
		vmsg = append(vmsg, msg)
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

func (u *User) validRecommendedBy() (errors.HSMessage, error) {
	if u.RecommendedBy != "" {
		var id string
		q := fmt.Sprintf("SELECT id FROM users WHERE invitation_code=%s", u.RecommendedBy)
		if err := dbQueryRow(q).Scan(&id); err != nil {
			if err != sql.ErrNoRows {
				return errors.HSMessage{}, err
			}
			return vemsgUserErrorRecommendCode, nil
		}
		u.RecommendedBy = id
	}
	return errors.HSMessage{}, nil
}

func (u *User) validUniqueKey() ([]errors.HSMessage, error) {
	var vemsgs []errors.HSMessage
	if u.Username != "" {
		if exist, err := u.isUserExistByUserName(); err != nil {
			return nil, err
		} else if exist {
			vemsgs = append(vemsgs, vemsgUserUsernameDuplicate)
		}
	}
	if u.Cellphone != "" {
		if exist, err := u.isUserExistByCellphone(); err != nil {
			return nil, err
		} else if exist {
			vemsgs = append(vemsgs, vemsgUserCellphoneDuplicate)
		}
	}
	if u.Email != "" {
		if exist, err := u.isUserExistByEmail(); err != nil {
			return nil, err
		} else if exist {
			vemsgs = append(vemsgs, vemsgUserEmailDuplicate)
		}
	}

	return vemsgs, nil
}

func (u *User) isUserExistByUserName() (bool, error) {
	return isItemExistInDbByProperty("users", "username", u.Username)
}

func (u *User) isUserExistByCellphone() (bool, error) {
	return isItemExistInDbByProperty("users", "cellphone", u.Cellphone)
}

func (u *User) isUserExistByEmail() (bool, error) {
	return isItemExistInDbByProperty("users", "email", u.Email)
}
