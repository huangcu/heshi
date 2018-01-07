package mysql

const userDdl = `
CREATE TABLE IF NOT EXISTS users
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	username VARCHAR(225) UNIQUE,
	cellphone VARCHAR(225) UNIQUE,
	email VARCHAR(225) UNIQUE,
	password VARCHAR(225) NOT NULL,
	user_type VARCHAR(25) NOT NULL,
	real_name VARCHAR(225),
	wechat_openid VARCHAR(225) UNIQUE,
	wechat_id VARCHAR(225) UNIQUE,
	wechat_name VARCHAR(225) UNIQUE,
	wechat_qr VARCHAR(225) UNIQUE,
	address VARCHAR(225),
	additional_info TEXT,
	recommended_by VARCHAR(225), 
	invitation_code VARCHAR(225) NOT NULL UNIQUE,
	discount INT NOT NULL DEFAULT 98,
	point INT NOT NULL DEFAULT 0,
	total_purchase_amount FLOAT NOT NULL DEFAULT 0,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	icon VARCHAR(255) DEFAULT "beyourdiamond.jpg"
) ENGINE=INNODB;
`
const adminDdl = `
CREATE TABLE IF NOT EXISTS admins
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	user_id VARCHAR(225) NOT NULL UNIQUE,
	level INT NOT NULL,
	wechat_kefu VARCHAR(225),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE=INNODB;
	`

//user VS agent
const agentDdl = `
CREATE TABLE IF NOT EXISTS agents
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	user_id VARCHAR(225) NOT NULL UNIQUE,
	level INT NOT NULL,
	discount FLOAT DEFAULT 0 NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE=INNODB;
`

//invitation code in user table VS seperate table????
const invitationCodeDdl = `
CREATE TABLE IF NOT EXISTS invitation_codes
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	user_id VARCHAR(225) NOT NULL,
	invitation_code VARCHAR(225) NOT NULL UNIQUE,
	discount INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`

// OpenId   string `json:"openid"`   // 用户的唯一标识
// Nickname string `json:"nickname"` // 用户昵称
// Sex      int    `json:"sex"`      // 用户的性别, 值为1时是男性, 值为2时是女性, 值为0时是未知
// City     string `json:"city"`     // 普通用户个人资料填写的城市
// Province string `json:"province"` // 用户个人资料填写的省份
// Country  string `json:"country"`  // 国家, 如中国为CN
// 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），
// 用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
// HeadImageURL string `json:"headimgurl,omitempty"`
// Privilege []string `json:"privilege,omitempty"` // 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
// UnionId   string   `json:"unionid,omitempty"`   // 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
const wechatUserDdl = `
CREATE TABLE IF NOT EXISTS wechat_users
(
	openid VARCHAR(225) PRIMARY KEY NOT NULL,
	nickname VARCHAR(225),
	sex TINYINT(4),
	city VARCHAR(225),
	province VARCHAR(225),
	contry VARCHAR(225),
	head_image_url VARCHAR(225),
	privilege VARCHAR(225),
	unionid VARCHAR(225) UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;`
