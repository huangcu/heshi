package mysql

const diamondDdl = `
CREATE TABLE IF NOT EXISTS diamonds
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	diamond_id VARCHAR(225) NOT NULL unique,
	stock_ref VARCHAR(225) NOT NULL unique,
	shape VARCHAR(225) NOT NULL,
	carat FLOAT NOT NULL,
	color VARCHAR(225),
	clarity VARCHAR(225),
	grading_lab VARCHAR(58) NOT NULL,
	certificate_number VARCHAR(225) NOT NULL,
	cut_grade VARCHAR(225) NOT NULL,
	polish VARCHAR(58) NOT NULL,
	symmetry VARCHAR(58) NOT NULL,
	fluorescence_intensity VARCHAR(58) NOT NULL,
	country VARCHAR(58) NOT NULL,
	supplier VARCHAR(15) NOT NULL,
	price_no_added_value DECIMAL(12,2) NOT NULL,
	price_retail DECIMAL(12,2) NOT NULL,
	featured VARCHAR(5) NOT NULL DEFAULT 'NO',
	recommand_words TEXT,
	extra_words VARCHAR(255),
	status VARCHAR(58) NOT NULL DEFAULT 'AVAIABLE',
	ordered_by VARCHAR(225),
	picked_up VARCHAR(8),
	sold_price FLOAT,
	profitable varchar(5) NOT NULL DEFAULT 'YES',
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
// <td width="88">镶碎钻</td> small_dias
// <td width="88">小钻数量</td> small_dia_num
// <td width="88">小钻总重</td> small_dia_carat
// <td width="88">镶嵌方式</td>  mounting_type
// <td width="88">价格</td> price
// unit_number 盒子号(库存归档)
// name_suffix INT,
const jewelryDdl = `
CREATE TABLE IF NOT EXISTS jewelrys
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	stock_id VARCHAR(225) NOT NULL unique,
	name VARCHAR(58) NOT NULL,
	need_diamond VARCHAR(5) NOT NULL,
	category VARCHAR(12) NOT NULL,
	mounting_type VARCHAR(225) NOT NULL,
	material VARCHAR(15) NOT NULL,
	metal_weight DECIMAL(12,2) NOT NULL,
	dia_shape VARCHAR(18) NOT NULL,
	price DECIMAL(12,2) NOT NULL,
	unit_number VARCHAR(28),
	dia_size_min FLOAT,
	dia_size_max FLOAT,
	main_dia_num INT,
	main_dia_size FLOAT,
	small_dias VARCHAR(8),
	small_dia_num INT,
	small_dia_carat DECIMAL(5,3),
	video_link VARCHAR(225),
	text TEXT,
	online VARCHAR(8) NOT NULL DEFAULT 'NO',
	verified VARCHAR(12) NOT NULL DEFAULT 'NO',
	in_stock VARCHAR(5) NOT NULL DEFAULT 'YES',
	featured VARCHAR(28) NOT NULL DEFAULT 'NO',
	stock_quantity TINYINT(4) NOT NULL DEFAULT 1,
	profitable varchar(5) NOT NULL DEFAULT 'YES',
	totally_scanned INT NOT NULL DEFAULT 0,
	free_acc VARCHAR(8) NOT NULL DEFAULT 'NO',
	offline_at TIMESTAMP,
	last_scan_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
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
	verified VARCHAR(12) NOT NULL DEFAULT 'NO',
	in_stock VARCHAR(5) NOT NULL DEFAULT 'YES',
	featured VARCHAR(28) NOT NULL DEFAULT 'NO',
	stock_quantity TINYINT(4) NOT NULL DEFAULT 1,
	profitable varchar(5) NOT NULL DEFAULT 'NO',
	totally_scanned INT NOT NULL DEFAULT 0,
	free_acc VARCHAR(8) NOT NULL DEFAULT 'NO',
	last_scan_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	offline_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
	`
