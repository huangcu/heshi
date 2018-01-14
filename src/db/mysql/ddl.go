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

//  exchange fluctuations
// deposit, balance
const discountDdl = `
CREATE TABLE IF NOT EXISTS discounts (
id VARCHAR(225) PRIMARY KEY NOT NULL,
discount_code VARCHAR(225) NOT NULL,
discount INT NOT NULL DEFAULT 98,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;`

const configDdl = `
CREATE TABLE IF NOT EXISTS configs (
id VARCHAR(225) PRIMARY KEY NOT NULL,
rate float NOT NULL,
created_by VARCHAR(225) NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`
