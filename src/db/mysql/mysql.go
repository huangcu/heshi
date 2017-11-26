package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// user_type: 0 admin, 1 agent, 2 customer
const userDdl = `
CREATE TABLE IF NOT EXISTS users
(
	id VARCHAR(255) PRIMARY KEY NOT NULL,
	username VARCHAR(225),
	cellphone VARCHAR(225),
	email VARCHAR(225),
	password VARCHAR(225),
	user_type INTEGER NOT NULL,
	real_name VARCHAR(225),
	wechat_id VARCHAR(225),
	wechat_name VARCHAR(225),
	address VARCHAR(225),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	icon VARCHAR(255) DEFAULT "beyourdiamond.jpg"
);
`
const adminDdl = `
CREATE TABLE IF NOT EXISTS admins
(
	id VARCHAR(255) PRIMARY KEY NOT NULL
	level INTEGER NOT NULL,
	wechat_kefu VARCHAR(225)
);
`
const agentDdl = `
CREATE TABLE IF NOT EXISTS agents
(
	id VARCHAR(255) PRIMARY KEY NOT NULL
	level INTEGER NOT NULL,
	discount float DEFAULT 0 NOT NULL
);
`

// const customer_ddl= `
// CREATE TABLE IF NOT EXISTS customers
// (
// 	id VARCHAR(255) PRIMARY KEY NOT NULL
// 	level INTEGER NOT NULL,
// );
// `

var Config = struct {
	Db struct {
		Database string
		User     string `default:"root" env:"heshi_mysql_user"`
		Password string `env:"heshi_mysql_password"`
		Port     uint   `default:"3306"`
	}
}{}

var dbCfgTeemplate = `
db: 
	database: heshi_dev
	user: root
	password: admin
	port: 3306
`

func OpenDB() (*sql.DB, error) {
	// df := "database." + os.Getenv("stage") + ".yml"
	return nil, nil
}
