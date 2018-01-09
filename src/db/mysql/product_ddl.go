package mysql

const diamondDdl = `
CREATE TABLE IF NOT EXISTS diamonds
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	diamond_id VARCHAR(225) NOT NULL unique,
	stock_ref VARCHAR(225) NOT NULL unique,
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
	clarity_number VARCHAR(5) NOT NULL,
	cut_number VARCHAR(5) NOT NULL,
	certificate_link TEXT,
	featured VARCHAR(5),
	recommand_words TEXT,
	extra_words VARCHAR(255),
	status VARCHAR(58) NOT NULL DEFAULT "AVAIABLE",
	ordered_by INT,
	picked_up VARCHAR(8),
	sold VARCHAR(5),
	sold_price FLOAT,
	profitable varchar(5),
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
