package main

import "heshi/errors"

var (
	//GENERAL
	vemsgShouldBeJSON       = errors.HSMessage{Code: 2000, Message: "should be JSON"}
	vemsgShouldNotBeEmpty   = errors.HSMessage{Code: 2000, Message: "should not be empty"}
	vemsgAlreadyExist       = errors.HSMessage{Code: 2000, Message: "already exists."}
	vemsgServerError        = errors.HSMessage{Code: 500, Message: "something is wrong, please try later"}
	vemsgAlreadyRecommended = errors.HSMessage{Code: 2000, Message: "您以前已经输入过一次推荐码，不需要再用其他推荐码了"}
	vemsgNoNeedRecommended  = errors.HSMessage{Code: 2000, Message: "您的用户级别已经很高，不需要再被别人推荐了"}
	vemsgCannotRecommend    = errors.HSMessage{Code: 2000, Message: "被您推荐的人不能再推荐您"}

	// //User Login (20-29)
	vemsgLoginErrorUserName   = errors.HSMessage{Code: 20020, Message: "wrong username or password"}
	vemsgLoginErrorPassword   = errors.HSMessage{Code: 20020, Message: "wrong password"}
	vemsgLoginMissingUserName = errors.HSMessage{Code: 20021, Message: "missing username or password"}

	//fail to find in db
	vemsgUserNotExist         = errors.HSMessage{Code: 20023, Message: "user not exist"}
	vemsgDiscountNotExist     = errors.HSMessage{Code: 20023, Message: "discount not exist"}
	vemsgExchangeRateNotExist = errors.HSMessage{Code: 20023, Message: "exchange rate not exist"}
	vemsgNotExist             = errors.HSMessage{Code: 20023, Message: "not exist"}

	//User Register(01-11)
	vemsgUserCellphoneEmailEmpty = errors.HSMessage{Code: 20001, Message: "you must input cellphone or email;"}
	vemsgUserPasswordEmpty       = errors.HSMessage{Code: 20002, Message: "password can not be empty;"}
	vemsgUserPasswordWarning     = errors.HSMessage{Code: 20003, Message: "密码请使用英文字母或数字组合"}
	vemsgUserUsernameDuplicate   = errors.HSMessage{Code: 20004, Message: "该帐户已存在"}
	vemsgUserUsernameError1      = errors.HSMessage{Code: 20004, Message: "user name length is 6 to 40;"}
	vemsgUserUsernameError2      = errors.HSMessage{Code: 20005, Message: "user name should only contain characher and number;"}
	vemsgUserEmailNotValid       = errors.HSMessage{Code: 20006, Message: "email input is not a valid email address;"}
	vemsgUserEmailDuplicate      = errors.HSMessage{Code: 20007, Message: "email already register!"}
	vemsgUserCellphoneNotValid   = errors.HSMessage{Code: 20008, Message: "cellphone input is not a valid cellphone number;"}
	vemsgUserCellphoneDuplicate  = errors.HSMessage{Code: 20009, Message: "cellphone already register!"}
	vemsgUserUsertypeNotValid    = errors.HSMessage{Code: 20010, Message: "user_type value is not valid;"}
	vemsgUserErrorRecommendCode  = errors.HSMessage{Code: 20011, Message: "your invitation code is not correct, please verify;"}

	//AGENT (11-19)
	vemsgAdminLevelNotValid         = errors.HSMessage{Code: 20015, Message: "admin level is not valid"}
	vemsgAgentLevelNotValid         = errors.HSMessage{Code: 20012, Message: "agent level is not valid"}
	vemsgAgentDiscountNotValid      = errors.HSMessage{Code: 20013, Message: "agent discount is not valid"}
	vemsgNotValid                   = errors.HSMessage{Code: 20014, Message: "input is not valid"}
	vemsgEmpty                      = errors.HSMessage{Code: 20014, Message: "input is cannot be empty"}
	vemsgDiamondCaratEmpty          = errors.HSMessage{Code: 20014, Message: "diamond carat size input cannot be empty"}
	vemsgDiamondCaratNotValid       = errors.HSMessage{Code: 20014, Message: "diamond carat size input is not valid"}
	vemsgDiamondRawPriceEmpty       = errors.HSMessage{Code: 20014, Message: "diamond raw price input cannot be empty"}
	vemsgDiamondRawPriceNotValid    = errors.HSMessage{Code: 20014, Message: "diamond raw price input is not valid"}
	vemsgDiamondRetailPriceEmpty    = errors.HSMessage{Code: 20014, Message: "diamond retail price input cannot be empty"}
	vemsgDiamondRetailPriceNotValid = errors.HSMessage{Code: 20014, Message: "diamond retail price input is not valid"}
	vemsgGemSizeEmpty               = errors.HSMessage{Code: 20014, Message: "gem size input cannot be empty"}
	vemsgGemSizeNotValid            = errors.HSMessage{Code: 20014, Message: "gem size input is not valid"}
	vemsgGemPriceEmpty              = errors.HSMessage{Code: 20014, Message: "gem price input cannot be empty"}
	vemsgGemPriceNotValid           = errors.HSMessage{Code: 20014, Message: "gem price input is not valid"}
	vemsgStockQuantityEmpty         = errors.HSMessage{Code: 20014, Message: "gem stock quantity cannot be empty"}
	vemsgStockQuantityNotValid      = errors.HSMessage{Code: 20014, Message: "gem stock quantity is not valid"}
	vemsgMetalWeightEmpty           = errors.HSMessage{Code: 20014, Message: "jewelry metal weight cannot be empty"}
	vemsgMetalWeightNotValid        = errors.HSMessage{Code: 20014, Message: "jewelry metal weight is not valid"}
	vemsgDiaSizeMinNotValid         = errors.HSMessage{Code: 20014, Message: "jewelry diamond size min is not valid"}
	vemsgDiaSizeMaxNotValid         = errors.HSMessage{Code: 20014, Message: "jewelry diamond size max is not valid"}
	vemsgMainDiaNumNotValid         = errors.HSMessage{Code: 20014, Message: "jewelry main diamond number is not valid"}
	vemsgMainDiaSizeNotValid        = errors.HSMessage{Code: 20014, Message: "jewelry main diamond size is not valid"}
	vemsgSmallDiaNumNotValid        = errors.HSMessage{Code: 20014, Message: "jewelry small diamond number is not valid"}
	vemsgSmallDiaCaratNotValid      = errors.HSMessage{Code: 20014, Message: "jewelry small diamond carat is not valid"}
	vemsgStockQuantityNotValidJ     = errors.HSMessage{Code: 20014, Message: "jewelry stock quantity is not valid"}
	vemsgPriceEmpty                 = errors.HSMessage{Code: 20014, Message: "jewelry price cannot be empty"}
	vemsgPriceNotValid              = errors.HSMessage{Code: 20014, Message: "jewelry price is not valid"}
	vemsgOrderPriceNotValid         = errors.HSMessage{Code: 20014, Message: "order price is not valid"}
	vemsgOrderDownPaymentNotValid   = errors.HSMessage{Code: 20014, Message: "order downpaymemnt is not valid"}
	vemsgOrderStatusNotValid        = errors.HSMessage{Code: 20014, Message: "order status input is not valid"}

	//UPLOAD PRODUCTS(30-39)
	vemsgUploadProductsCategoryNotValid = errors.HSMessage{Code: 20090, Message: "product category not valid"}
	vemsgSupplierNotValid               = errors.HSMessage{Code: 20080, Message: "supplier is not valid"}

	//Supplier
	vemsgSupplierNameDuplicate   = errors.HSMessage{Code: 20081, Message: "the supplier name already exist!"}
	vemsgSupplierPrefixDuplicate = errors.HSMessage{Code: 20082, Message: "the supplier prefix already exist!"}

	//Currency(90-99)
	vemsgCurrencySymbolNotValid  = errors.HSMessage{Code: 20090, Message: "input is not a valid currency symbol;"}
	vemsgCurrencyRateNotValid    = errors.HSMessage{Code: 20091, Message: "currency exchange rate should be float;"}
	vemsgCurrencyBaseNotValid    = errors.HSMessage{Code: 20092, Message: "currency exchange rate base can only be USD for now!;"}
	vemsgCurrencyRateEURNotValid = errors.HSMessage{Code: 20093, Message: "EUR currency exchange rate not valid!;"}
	vemsgCurrencyRateCNYNotValid = errors.HSMessage{Code: 20094, Message: "CNY currency exchange rate not valid;"}
)

var (
	VALID_USERTYPE        = []string{CUSTOMER, AGENT, ADMIN}
	VALID_AGENTLEVEL      = []string{LEVEL1, LEVEL2, LEVEL3}
	VALID_ADMINLEVEL      = []string{LEVEL1, LEVEL2, LEVEL3, LEVEL4, LEVEL5, LEVEL6}
	VALID_CURRENCY_SYMBOL = []string{"USD", "CNY", "EUR", "CAD", "AUD", "CHF", "RUB", "NZD"}
	VALID_ORDER_STATUS_M  = []string{MCANCELLED, MDOWNPAYMENT, MPAID, MDELIVERED, MRECEIVED}
	VALID_ORDER_STATUS_A  = []string{ORDERED, CANCELLED, DOWNPAYMENT, PAID, DELIVERED, RECEIVED}
	USER_SESSION_KEY      = "hs_sessionuserid"
	ADMIN_KEY             = "hs_sessionadmin"
	AGENT_KEY             = "hs_sessionagent"

	UPLOADFILEDIR        = ".uploaded"
	videoPath            = ".video"
	imagePath            = ".image"
	userIconSizeLimit    = 50 * 1024
	imageSizeLimit       = 5 * 1024 * 1024
	videoSizeLimit       = 10 * 1024 * 1024
	uploadedZipFileLimit = 100 * 1024 * 1024
)

const (
	CUSTOMER     = "CUSTOMER"
	AGENT        = "AGENT"
	ADMIN        = "ADMIN"
	ORDERED      = "ORDERED"
	CANCELLED    = "CANCELLED"
	DOWNPAYMENT  = "DOWNPAYMENT"
	PAID         = "PAID"
	DELIVERED    = "DELIVERED"
	RECEIVED     = "RECEIVED"
	MCANCELLED   = "M-CANCELLED"
	MDOWNPAYMENT = "M-DOWNPAYMENT"
	MPAID        = "M-PAID"
	MDELIVERED   = "M-DELIVERED"
	MRECEIVED    = "M-RECEIVED"
)

//TODO level comes from configs db
const (
	LEVEL1 = "1"
	LEVEL2 = "2"
	LEVEL3 = "3"
	LEVEL4 = "4"
	LEVEL5 = "5"
	LEVEL6 = "6"
)
