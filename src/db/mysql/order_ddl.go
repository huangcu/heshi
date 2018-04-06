package mysql

//same transaction_id, multiple items, multiple order records
//if only one, id is same value with transaction_id
//status: ORDERED/CANCELLED/SOLD/DOWNPAYMENT - > cancel order without downpayment in 1 day
//agent return_point: based on payment time's (agent level TODO)???
const orderDdl = `
CREATE TABLE IF NOT EXISTS orders
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	item_id INT NOT NULL,
	item_price FLOAT NOT NULL,
	item_category VARCHAR(20) NOT NULL,
	item_quantity INT NOT NULL DEFAULT 1,
	transaction_id VARCHAR(225) NOT NULL,
	buyer_id VARCHAR(225) NOT NULL,
	downpayment FLOAT,
	chosen_by VARCHAR(28),
	sold_price_usd FLOAT,
	sold_price_cny FLOAT,
	sold_price_eur FLOAT,
	return_point FLOAT,
	status VARCHAR(20) NOT NULL DEFAULT 'ORDERED',
	extra_info VARCHAR(225),
	special_notice VARCHAR(225),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`

// 购物车
const shoppingCart = `
CREATE TABLE IF NOT EXISTS shopping_cart
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	item_id INT NOT NULL,
	item_price FLOAT NOT NULL,
	item_category VARCHAR(20) NOT NULL,
	item_quantity INT NOT NULL DEFAULT 1,
	user_id VARCHAR(225) NOT NULL,
	extra_info VARCHAR(225),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`

// 我的收藏 TODO to be removed
const interestedItemDdl = `
CREATE TABLE IF NOT EXISTS interested_items
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	user_id VARCHAR(225) NOT NULL,
	item_type VARCHAR(28) NOT NULL,
	item_id VARCHAR(58) NOT NULL,
	item_accessory int(11),
	confirmed_for_check VARCHAR(8) NOT NULL DEFAULT 'No',
	available VARCHAR(25) NOT NULL DEFAULT 'TOBECHECKED',
	special_notice VARCHAR(225),
	event_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`
