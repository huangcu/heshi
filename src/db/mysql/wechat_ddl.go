package mysql

const messageDdl = `
CREATE TABLE IF NOT EXISTS messages
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	from_user VARCHAR(225),
	msg_type VARCHAR(25) NOT NULL,
	context VARCHAR(225),
	kf_account VARCHAR(225),
	media_id VARCHAR(225),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`
