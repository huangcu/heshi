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
const diamondDdl = `
CREATE TABLE IF NOT EXISTS diamonds
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	stock_ref VARCHAR(225) NOT NULL,
	shape VARCHAR(225),
	carat FLOAT NOT NULL,
	color VARCHAR(225),
	clarity VARCHAR(225),
	grading_lab INT NOT NULL,
	certificate_number VARCHAR(225),
	cut_grade VARCHAR(225),
	polish VARCHAR(58),
	symmetry VARCHAR(58),
	fluorescence_intensity VARCHAR(58),
	country VARCHAR(58) NOT NULL,
	supplier VARCHAR(15) NOT NULL,
	price_no_added_value DECIMAL(12,2) NOT NULL,
	price_retail DECIMAL(12,2) NOT NULL,
	certificate_link TEXT,
	clarity_number VARCHAR(5) NOT NULL,
	cut_number VARCHAR(5) NOT NULL,
	featured VARCHAR(5) NOT NULL,
	recommand_words TEXT,
	extra_words VARCHAR(255),
	status VARCHAR(58) NOT NULL,
	ordered_by INT,
	picked_up VARCHAR(8) NOT NULL,
	sold VARCHAR(5) NOT NULL,
	sold_price FLOAT,
	profitable varchar(5) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
	`

const jewelryDdl = `
CREATE TABLE IF NOT EXISTS jewelrys
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	category INT NOT NULL,
	unit_number VARCHAR(28) NOT NULL,
	dia_shape VARCHAR(18),
	material VARCHAR(15),
	metal_weight DECIMAL(12,2),
	need_diamond VARCHAR(5) NOT NULL,
	name VARCHAR(58) NOT NULL,
	name_suffix INT,
	dia_size_min FLOAT NOT NULL,
	dia_size_max FLOAT NOT NULL,
	small_dias VARCHAR(8),
	small_dia_num INT,
	small_dia_carat DECIMAL(5,3),
	mounting_type VARCHAR(225),
	main_dia_num INT NOT NULL,
	main_dia_size VARCHAR(225),
	video_link VARCHAR(225),
	text TEXT,
	online VARCHAR(8) NOT NULL,
	verified VARCHAR(12) NOT NULL,
	in_stock VARCHAR(5) NOT NULL,
	featured VARCHAR(28) NOT NULL,
	price DECIMAL(12,2),
	stock_quantity TINYINT(4) NOT NULL,
	profitable varchar(5) NOT NULL,
	clarity_number VARCHAR(5) NOT NULL,
	totally_scanned INT NOT NULL,
	free_acc VARCHAR(8) NOT NULL,
	last_scan_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	offline_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
	`

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
const smallDiamondDdl = `
CREATE TABLE IF NOT EXISTS small_diamonds
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	size_from DECIMAL(4,3) NOT NULL,
	size_to DECIMAL(4,3) NOT NULL,
	price DECIMAL(12,2) NOT NULL,
	quantity TINYINT(4) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`
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
