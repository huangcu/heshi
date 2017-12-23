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

// FOREIGN KEY (diamond_id) REFERENCES diamonds (id),
const supplierDdl = `
CREATE TABLE IF NOT EXISTS suppliers
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	name VARCHAR(255) NOT NULL,
	prefix VARCHAR(8) NOT NULL,
	connected VARCHAR(5) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`

const appointmentDdl = `
CREATE TABLE IF NOT EXISTS appointments
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	item_id INT NOT NULL,
	item_price FLOAT NOT NULL,
	item_category INT NOT NULL,
	buyer_id TINYINT(4) NOT NULL,
	chosen_by VARCHAR(28) NOT NULL,
	extra_info VARCHAR(255),
	ordered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`
const orderDdl = `
CREATE TABLE IF NOT EXISTS orders
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	appointment_id VARCHAR(225) NOT NULL,
	FOREIGN KEY (appointment_id) REFERENCES appointments (id)
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
) ENGINE=INNODB`
