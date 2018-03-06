package mysql

//Direction: TO user or FROM user
const messageDdl = `
CREATE TABLE IF NOT EXISTS messages
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	user VARCHAR(225),
	msg_type VARCHAR(25) NOT NULL DEFAULT 'TEXT',
	content VARCHAR(225),
	kf_account VARCHAR(225),
	direction VARCHAR(8),
	media_id VARCHAR(225),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`
