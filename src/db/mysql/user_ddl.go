package mysql

const userDdl = `
CREATE TABLE IF NOT EXISTS users
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	username VARCHAR(225) UNIQUE,
	cellphone VARCHAR(225) UNIQUE,
	email VARCHAR(225) UNIQUE,
	password VARCHAR(225),
	user_type VARCHAR(25) NOT NULL,
	real_name VARCHAR(225),
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
	invitation_code VARCHAR(20) NOT NULL UNIQUE,
	FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE=INNODB;
`

//invitation code in user table VS seperate table????
const invitationCodeDdl = `
CREAT TABLE IF NOT EXISTS invitation_codes
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	user_id VARCHAR(225) NOT NULL,
	invitation_code VARCHAR(225) NOT NULL UNIQUE,
	discount NOT NULL
) ENGINE=INNODB;
`
