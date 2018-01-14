package mysql

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

const interestedItemDdl = `
CREATE TABLE IF NOT EXISTS interested_items
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	user_id VARCHAR(225) NOT NULL,
	item_type VARCHAR(28) NOT NULL,
	item_id VARCHAR(58) NOT NULL,
	item_accessory int(11),
	confirmed_for_check VARCHAR(8) NOT NULL DEFAULT 'No',
	available VARCHAR(8) NOT NULL DEFAULT 'TOBECHECKED',
	special_notice VARCHAR(225),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`
