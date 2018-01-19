package mysql

// 	+----+---------------+-----------+-----------+
// | id | supplier_name | id_prefix | connected |
// +----+---------------+-----------+-----------+
// |  1 | KGK           | 1U        | NA        |
// |  2 | DIAM          | DIAM      | NA        |
// |  3 | SUN           | SUN       | NA        |
// |  4 | BEYOU-HESHI   | HS        | NA        |
// |  5 | PG            | PG        | NA        |
// |  8 | HST           | HST       | YES       |
// |  7 | pk            |           | NA        |
// |  9 | CN            | CN        | YES       |

// FOREIGN KEY (diamond_id) REFERENCES diamonds (id),
const supplierDdl = `
CREATE TABLE IF NOT EXISTS suppliers
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	name VARCHAR(255) NOT NULL,
	prefix VARCHAR(8) NOT NULL,
	connected VARCHAR(5) NOT NULL DEFAULT 'NO',
	status VARCHAR(8) NOT NULL DEFAULT 'active',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`

// +----------------+--------------+------+-----+---------+----------------+
// | id             | int(11)      | NO   | PRI | NULL    | auto_increment |
// | carat_from     | float        | NO   |     | NULL    |                |
// | carat_to       | float        | NO   |     | NULL    |                |
// | color          | varchar(88)  | NO   |     | NULL    |                |
// | clarity        | varchar(88)  | NO   |     | NULL    |                |
// | cut            | varchar(255) | NO   |     | NULL    |                |
// | symmetry       | varchar(255) | NO   |     | NULL    |                |
// | polish         | varchar(255) | NO   |     | NULL    |                |
// | fluo           | varchar(255) | NO   |     | NULL    |                |
// | certificate    | varchar(58)  | NO   |     | NULL    |                |
// | shape          | varchar(28)  | NO   |     | NULL    |                |
// | the_para_value | float        | NO   |     | NULL    |                |
// | priority       | int(11)      | NO   |     | NULL    |                |
const priceSettingDdl = `
CREATE TABLE IF NOT EXISTS price_settings
(
	id VARCHAR(225) PRIMARY KEY NOT NULL,
	carat_from FLOAT NOT NULL,
	carat_to FLOAT NOT NULL,
	color VARCHAR(88) NOT NULL,
	clarity VARCHAR(88) NOT NULL,
	cut VARCHAR(225) NOT NULL,
	symmetry VARCHAR(225) NOT NULL,
	polish VARCHAR(225) NOT NULL,
	fluo VARCHAR(225) NOT NULL,
	certificate VARCHAR(225) NOT NULL,
	the_para_value FLOAT NOT NULL,
	priority INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`
