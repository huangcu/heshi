package mysql

//Direction: TO user or FROM user
const wechatMessageDdl = `
CREATE TABLE IF NOT EXISTS wechat_messages
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	user VARCHAR(225),
	msg_type VARCHAR(25) NOT NULL DEFAULT 'TEXT',
	content VARCHAR(225),
	kf_account VARCHAR(225),
	direction VARCHAR(8),
	pic_url VARCHAR(225),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`

const messageDdl = `
CREATE TABLE IF NOT EXISTS messages (
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	user_id VARCHAR(225) NOT NULL,
	msg_type VARCHAR(25) NOT NULL DEFAULT 'TEXT',
	content TEXT NOT NULL,
	serve_id VARCHAR(225),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`
