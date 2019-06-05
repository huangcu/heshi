package mysql

// PROM_TYPE: DISCOUNT, FREE_ACCESSORY, SPECIAL_OFFER
// STATUS: ACTIVE, INACTIVE
const promotionDdl = `
CREATE TABLE IF NOT EXISTS promotions
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	prom_type VARCHAR(20),
	prom_discount TINYINT(4),
	prom_price DECIMAL(12,2),
	status VARCHAR(10) DEFAULT 'INACTIVE' NOT NULL,
	begin_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	end_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	created_by VARCHAR(225) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	) ENGINE=INNODB;
	`

//note: where rate comes from
const currencyExchangeRateDdl = `
CREATE TABLE IF NOT EXISTS currency_exchange_rates
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	note VARCHAR(225) NOT NULL,
	base VARCHAR(4) NOT NULL DEFAULT 'USD',
	usd DECIMAL(12,6) NOT NULL,
	cny DECIMAL(12,6) NOT NULL,
	eur DECIMAL(12,6) NOT NULL,
	cad DECIMAL(12,6) NOT NULL,
	aud DECIMAL(12,6) NOT NULL,
	chf DECIMAL(12,6) NOT NULL,
	rub DECIMAL(12,6) NOT NULL,
	nzd DECIMAL(12,6) NOT NULL,
	usd_fluc DECIMAL(6,3),
	cny_fluc DECIMAL(6,3),
	eur_fluc DECIMAL(6,3),
	cad_fluc DECIMAL(6,3),
	aud_fluc DECIMAL(6,3),
	chf_fluc DECIMAL(6,3),
	rub_fluc DECIMAL(6,3),
	nzd_fluc DECIMAL(6,3),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	) ENGINE=INNODB;
`

const discountDdl = `
CREATE TABLE IF NOT EXISTS discounts (
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	discount_code VARCHAR(225) NOT NULL,
	discount TINYINT(4) NOT NULL DEFAULT 98,
	created_by VARCHAR(225) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	) ENGINE=INNODB;`

//TYPE: customer(level), agent(level), rate
const levelRateDdl = `
CREATE TABLE IF NOT EXISTS level_rate_rules (
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	exchange_rate_float DECIMAL(6,3),
	level TINYINT(4),
	discount TINYINT(4),
	amount INT,
	pieces INT,
	return_point_percent INT,
	rule_type VARCHAR(20) NOT NULL,
	created_by VARCHAR(225) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`

//to track use activity - pages
const userActiveRecordDdl = `
CREATE TABLE IF NOT EXISTS user_active_records (
id VARCHAR(225) PRIMARY KEY NOT NULL,
user_id VARCHAR(225) NOT NULL,
page VARCHAR(225) NOT NULL,
remote_addr VARCHAR(50),
device VARCHAR(25) NOT NULL,
action_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`

//userusingrecord - track browser items (which diamond, which jewerly etc ) using PC MOBILE etc
//item_id (diamond stock_ref, jewelry stock_id, item_type: diamond or jewelry)
const userUsingRecordDdl = `
CREATE TABLE IF NOT EXISTS user_using_records (
id VARCHAR(225) PRIMARY KEY NOT NULL,
user_id VARCHAR(225) NOT NULL,
item_id VARCHAR(225) NOT NULL,
item_type VARCHAR(20) NOT NULL,
device VARCHAR(25) NOT NULL,
remote_addr VARCHAR(50),
action_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`

// id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
// id VARCHAR(225) PRIMARY KEY NOT NULL,
const historyTrackDdl = `
CREATE TABLE IF NOT EXISTS historys (
id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
user_id VARCHAR(225) NOT NULL,
item_id VARCHAR(225) NOT NULL,
table_name VARCHAR(30) NOT NULL,
field_name VARCHAR(50) NOT NULL,
new_value VARCHAR(225) NOT NULL,
old_value VARCHAR(225) NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`
