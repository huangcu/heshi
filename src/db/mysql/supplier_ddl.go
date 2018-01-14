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
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
) ENGINE=INNODB;
`
