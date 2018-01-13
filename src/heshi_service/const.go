package main

import "heshi/errors"

var (
	//GENERAL
	VEMSG_SHOULD_BE_JSON      = errors.HSMessage{Code: 2000, Message: "should be JSON"}
	VEMSG_SHOULD_NOT_BE_EMPTY = errors.HSMessage{Code: 2000, Message: "should not be empty"}
	VEMSG_ALREADY_EXIST       = errors.HSMessage{Code: 2000, Message: "already exists."}
	VEMSG_SERVER_ERROR        = errors.HSMessage{Code: 2000, Message: "something is wrong, please try later"}
	// VEMSG_ERROR_RECOMMAND_CODE   =errors.HSMessage{Code: 2000, Message: "您的推荐码不正确，请核实"}

	VEMSG_ALREADY_RECOMMANDED = errors.HSMessage{Code: 2000, Message: "您以前已经输入过一次推荐码，不需要再用其他推荐码了"}
	VEMSG_NONEED_RECOMMANDED  = errors.HSMessage{Code: 2000, Message: "您的用户级别已经很高，不需要再被别人推荐了"}
	VEMSG_CANNOT_RECOMMAND    = errors.HSMessage{Code: 2000, Message: "被您推荐的人不能再推荐您"}
	VEMSG_RECOMMAND_TOOTHER   = errors.HSMessage{Code: 2000, Message: "请将该页发送给你的朋友"}

	// //User Login (20-29)
	VEMSG_LOGIN_ERROR_USERNAME = errors.HSMessage{Code: 20020, Message: "Error: login info not correct, wrong user or password"}
	// VEMSG_LOGIN_ERROR_EMAIL     = errors.HSMessage{Code: 20021, Message: "邮箱或密码错误，请重试"}
	// VEMSG_LOGIN_ERROR_CELLPHONE = errors.HSMessage{Code: 20022, Message: "电话号码或密码错误，请重试"}

	//fail to find in db
	VEMSG_USER_NOT_EXIST          = errors.HSMessage{Code: 20023, Message: "user not exist"}
	VEMSG_DISCOUNT_NOT_EXIST      = errors.HSMessage{Code: 20023, Message: "discount not exist"}
	VEMSG_EXCHANGE_RATE_NOT_EXIST = errors.HSMessage{Code: 20023, Message: "exchange rate not exist"}
	VEMSG_NOT_EXIST               = errors.HSMessage{Code: 20023, Message: "not exist"}

	//User Register(01-11)
	VEMSG_USER_CELLPHONE_EMAIL_EMPTY = errors.HSMessage{Code: 20001, Message: "you must input cellphone or email;"}
	VEMSG_USER_PASSWORD_EMPTY        = errors.HSMessage{Code: 20002, Message: "password can not be empty;"}
	VEMSG_USER_PASSWORD_WARNING      = errors.HSMessage{Code: 20003, Message: "密码请使用英文字母或数字组合"}
	VEMSG_USER_USERNAME_DUPLICATE    = errors.HSMessage{Code: 20004, Message: "该帐户已存在"}
	VEMSG_USER_USERNAME_ERROR1       = errors.HSMessage{Code: 20004, Message: "user name length is 6 to 40;"}
	VEMSG_USER_USERNAME_ERROR2       = errors.HSMessage{Code: 20005, Message: "user name should only contain characher and number;"}
	VEMSG_USER_EMAIL_NOT_VALID       = errors.HSMessage{Code: 20006, Message: "email input is not a valid email address;"}
	VEMSG_USER_EMAIL_DUPLICATE       = errors.HSMessage{Code: 20007, Message: "email already register!"}
	VEMSG_USER_CELLPHONE_NOT_VALID   = errors.HSMessage{Code: 20008, Message: "cellphone input is not a valid cellphone number;"}
	VEMSG_USER_CELLPHONE_DUPLICATE   = errors.HSMessage{Code: 20009, Message: "cellphone already register!"}
	VEMSG_USER_USERTYPE_NOT_VALID    = errors.HSMessage{Code: 20010, Message: "user_type value is not valid;"}
	VEMSG_USER_ERROR_RECOMMAND_CODE  = errors.HSMessage{Code: 20011, Message: "your invitation code is not correct, please verify;"}

	//AGENT (11-19)
	VEMSG_AGENT_LEVEL_NOT_VALID    = errors.HSMessage{Code: 20012, Message: "agent level is notvalid"}
	VEMSG_AGENT_DISCOUNT_NOT_VALID = errors.HSMessage{Code: 20013, Message: "discount is not valid"}

	//Currency(90-99)
	VEMSG_CURRENCY_SYMBOL_NOT_VALID   = errors.HSMessage{Code: 20090, Message: "input is not a valid currency symbol;"}
	VEMSG_CURRENCY_RATE_NOT_VALID     = errors.HSMessage{Code: 20091, Message: "currency exchange rate should be float;"}
	VEMSG_CURRENCY_BASE_NOT_VALID     = errors.HSMessage{Code: 20092, Message: "currency exchange rate base can only be USD for now!;"}
	VEMSG_CURRENCY_RATE_EUR_NOT_VALID = errors.HSMessage{Code: 20093, Message: "EUR currency exchange rate not valid!;"}
	VEMSG_CURRENCY_RATE_CNY_NOT_VALID = errors.HSMessage{Code: 20094, Message: "CNY currency exchange rate not valid;"}
)

var (
	letterRunes           = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	VALID_USERTYPE        = []string{"admin", "customer", "agent"}
	VALID_AGENTLEVEL      = []string{"1", "2", "3"}
	VALID_CURRENCY_SYMBOL = []string{"USD", "CNY", "EUR", "CAD", "AUD", "CHF", "RUB", "NZD"}
	USER_SESSION_KEY      = "hs_sessionuserid"
	ADMIN_KEY             = "hs_sessionadmin"
)

const (
	CUSTOMER = "customer"
	AGENT    = "agent"
	ADMIN    = "admin"
)

const (
	LEVEL1 = "1"
	LEVEL2 = "2"
	LEVEL3 = "3"
)
