package mysql

const promotionDdl = `
CREATE TABLE IF NOT EXISTS promotions
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	diamond_id VARCHAR(225),
	jewelry_id VARCHAR(225),
	type VARCHAR(5),
	price DECIMAL(12,2),
	begin TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	end TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	FOREIGN KEY (jewelry_id) REFERENCES jewelrys (id)
	) ENGINE=INNODB;
	`

//note: where rate comes from
const currencyExchangeRateDdl = `
CREATE TABLE IF NOT EXISTS currency_exchange_rates
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	note VARCHAR(225) NOT NULL,
	base VARCHAR(4) NOT NULL DEFAULT 'USD',
	usd FLOAT NOT NULL,
	cny FLOAT NOT NULL,
	eur FLOAT NOT NULL,
	cad FLOAT NOT NULL,
	aud FLOAT NOT NULL,
	chf FLOAT NOT NULL,
	rub FLOAT NOT NULL,
	nzd FLOAT NOT NULL,
	usd_fluc FLOAT,
	cny_fluc FLOAT,
	eur_fluc FLOAT,
	cad_fluc FLOAT,
	aud_fluc FLOAT,
	chf_fluc FLOAT,
	rub_fluc FLOAT,
	nzd_fluc FLOAT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	) ENGINE=INNODB;
`

const discountDdl = `
CREATE TABLE IF NOT EXISTS discounts (
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	discount_code VARCHAR(225) NOT NULL,
	discount INT NOT NULL DEFAULT 98,
	created_by VARCHAR(225) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	) ENGINE=INNODB;`

//TYPE: customer(level), agent(level), rate
const levelRateDdl = `
CREATE TABLE IF NOT EXISTS level_rate_rules (
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	exchange_rate_float FLOAT,
	level VARCHAR(20),
	discount INT,
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
item_type VARCHAR(225) NOT NULL,
device VARCHAR(25) NOT NULL,
action_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`

const historyTrackDdl = `
CREATE TABLE IF NOT EXISTS historys (
id VARCHAR(225) PRIMARY KEY NOT NULL,
user_id VARCHAR(225) NOT NULL,
item_id VARCHAR(225) NOT NULL,
table_name VARCHAR(20) NOT NULL,
field_name VARCHAR(10) NOT NULL,
field_value VARCHAR(25) NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`
