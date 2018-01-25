package mysql

const diamondDdl = `
CREATE TABLE IF NOT EXISTS diamonds
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	diamond_id VARCHAR(225) NOT NULL unique,
	stock_ref VARCHAR(225) NOT NULL unique,
	shape VARCHAR(225) NOT NULL DEFAULT '-',
	carat FLOAT NOT NULL,
	color VARCHAR(225) DEFAULT '-',
	clarity VARCHAR(225) DEFAULT '-',
	grading_lab VARCHAR(58) NOT NULL DEFAULT '-',
	certificate_number VARCHAR(225) NOT NULL DEFAULT '-',
	cut_grade VARCHAR(225) NOT NULL DEFAULT '-',
	polish VARCHAR(58) NOT NULL DEFAULT '-',
	symmetry VARCHAR(58) NOT NULL DEFAULT '-',
	fluorescence_intensity VARCHAR(58) NOT NULL DEFAULT '-',
	country VARCHAR(58) NOT NULL,
	supplier VARCHAR(15) NOT NULL,
	price_no_added_value DECIMAL(12,2) NOT NULL,
	price_retail DECIMAL(12,2) NOT NULL,
	featured VARCHAR(5),
	recommand_words TEXT,
	extra_words VARCHAR(255),
	status VARCHAR(58) NOT NULL DEFAULT "AVAIABLE",
	ordered_by VARCHAR(225),
	picked_up VARCHAR(8),
	sold_price FLOAT,
	profitable varchar(5),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
	`

// category
// <option value="JR">素金戒指</option>
// <option value="JE">素金耳环／耳钉</option>
// <option value="JP">素金吊坠／项链</option>
// <option value="ZR">镶碎钻戒指</option>
// <option value="ZE">镶碎钻耳环／耳钉</option>
// <option value="ZP">镶碎钻吊坠／项链</option>
// <option value="CR">成品戒指</option>
// <option value="CE">成品耳环／耳钉</option>
// <option value="CP">成品吊坠／项链</option>

// <td width="88">唯一商品号(StockID)</td>
// <td width="88">货号(Name)</td>
// <td width="88">材料</td> material
// <td width="88">金重</td>     metal_weight
// <td width="88">是否空托</td>  need_diamond
// <td width="88">最小钻石尺寸</td> dia_size_min
// <td width="88">最大钻石尺寸</td> dia_size_max
// <td width="88">镶碎钻</td>
// <td width="88">小钻数量</td>
// <td width="88">小钻总重</td>
// <td width="88">镶嵌方式</td>  mounting_type
// <td width="88">价格</td> price
const jewelryDdl = `
CREATE TABLE IF NOT EXISTS jewelrys
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	category INT NOT NULL,
	name VARCHAR(58) NOT NULL,
	name_suffix INT,
	stock_id VARCHAR(225) NOT NULL unique,
	metal_weight DECIMAL(12,2),
	material VARCHAR(15),
	need_diamond VARCHAR(5) NOT NULL,
	dia_size_min FLOAT,
	dia_size_max FLOAT,
	mounting_type VARCHAR(225),
	price DECIMAL(12,2),
	unit_number VARCHAR(28),
	dia_shape VARCHAR(18),
	small_dias VARCHAR(8),
	small_dia_num INT,
	small_dia_carat DECIMAL(5,3),
	main_dia_num INT,
	main_dia_size FLOAT,
	video_link VARCHAR(225),
	text TEXT,
	online VARCHAR(8) NOT NULL,
	verified VARCHAR(12) NOT NULL,
	in_stock VARCHAR(5) NOT NULL,
	featured VARCHAR(28) NOT NULL,
	stock_quantity TINYINT(4) NOT NULL,
	profitable varchar(5) NOT NULL,
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

const gemDdl = `
CREATE TABLE IF NOT EXISTS gems
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	name VARCHAR(58) NOT NULL,
	stock_id VARCHAR(225) NOT NULL unique,
	size FLOAT NOT NULL,
	material VARCHAR(15) NOT NULL,
	price DECIMAL(12,2) NOT NULL,
	shape VARCHAR(18) NOT NULL,
	text TEXT NOT NULL,
	certificate VARCHAR(225) NOT NULL unique,
	online VARCHAR(8) NOT NULL DEFAULT 'NO',
	verified VARCHAR(12) NOT NULL DEFAULT 'YES',
	in_stock VARCHAR(5) NOT NULL DEFAULT 'YES',
	featured VARCHAR(28) NOT NULL DEFAULT 'NO',
	stock_quantity TINYINT(4) NOT NULL DEFAULT 1,
	profitable varchar(5) NOT NULL DEFAULT 'NOT',
	totally_scanned INT NOT NULL DEFAULT 0,
	free_acc VARCHAR(8) NOT NULL DEFAULT 'NOT',
	last_scan_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	offline_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
	`
