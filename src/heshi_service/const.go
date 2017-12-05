package main

const (
	VEMSG_SHOULD_BE_JSON         = "should be JSON"
	VEMSG_SHOULD_NOT_BE_EMPTY    = "should not be empty"
	VEMSG_ALREADY_EXIST          = "already exists."
	VEMSG_SERVER_ERROR           = "something is wrong, please try later"
	VEMSG_ERROR_RECOMMAND_CODE   = "您的推荐码不正确，请核实"
	VEMSG_ALREADY_RECOMMANDED    = "您以前已经输入过一次推荐码，不需要再用其他推荐码了"
	VEMSG_NONEED_RECOMMANDED     = "您的用户级别已经很高，不需要再被别人推荐了"
	VEMSG_CANNOT_RECOMMAND       = "被您推荐的人不能再推荐您"
	VEMSG_RECOMMAND_TOOTHER      = "请将该页发送给你的朋友"
	VEMSG_ACCOUNT_ALREADY_EXISTS = "该帐户已存在"
	VEMSG_PASSWORD_WARNING       = "密码请使用英文字母或数字组合"
	VEMSG_LOGIN_ERROR_EMAIL      = "邮箱或密码错误，请重试"
	VEMSG_LOGIN_ERROR_CELLPHONE  = "电话号码或密码错误，请重试"
	VEMSG_LOGIN_ERROR_USERNAME   = "用户名或密码错误，请重试"

	VEMSG_USER_CELLPHONE_EMAIL_EMPTY = "you must input cellphone or email;"
	VEMSG_USER_PASSWORD_EMPTY        = "password can not be empty;"
	VEMSG_USER_USERTYPE_EMPTY        = "user type can not be empty;"
	VEMSG_USER_USERNAME_ERROR1       = "user name length is 6 to 40;"
	VEMSG_USER_USERNAME_ERROR2       = "user name should only contain characher and number;"
	VEMSG_USER_USERTYPE_NOT_VALID    = "user_type value is not valid;"
	VEMSG_USER_EMAIL_NOT_VALID       = "email input is not a valid email address;"
	VEMSG_USER_CELLPHONE_NOT_VALID   = "cellphone input is not a valid cellphone number;"
)

var (
	VALID_USERTYPE = []string{"admin", "customer", "agent"}
)
