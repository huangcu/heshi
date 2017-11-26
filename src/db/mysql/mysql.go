package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/configor"
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
	id VARCHAR(255) PRIMARY KEY NOT NULL,
	level INTEGER NOT NULL,
	wechat_kefu VARCHAR(225)
);
`
const agentDdl = `
CREATE TABLE IF NOT EXISTS agents
(
	id VARCHAR(255) PRIMARY KEY NOT NULL,
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

// var dbCfgTeemplate = `
// db:
// 	database: heshi_dev
// 	user: root
// 	password: admin
// 	port: 3306
// `

func OpenDB() (*sql.DB, error) {
	// df := "database." + os.Getenv("stage") + ".yml"
	if err := configor.Load(&Config, "config.yml"); err != nil {
		return nil, err
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/", Config.Db.User, Config.Db.Password))
	if err != nil {
		return nil, err
	}

	createDatabase(db, Config.Db.Database)
	// if err := createDatabase(db, Config.Db.Database); err != nil {
	// 	return nil, err
	// }

	for table, ddl := range map[string]string{
		"users":  userDdl,
		"agents": agentDdl,
		"admins": adminDdl,
	} {
		if _, err := db.Exec(ddl); err != nil {
			return nil, fmt.Errorf("fail to create table %s with %s", table, ddl)
		}
	}

	db.Close()

	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?parseTime=true",
		Config.Db.User, Config.Db.Password, Config.Db.Database))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createDatabase(db *sql.DB, name string) {
	q := `CREATE DATABASE IF NOT EXISTS %s 
	DEFAULT CHARACTER SET utf8 
	DEFAULT COLLATE utf8_general_ci`
	ddl := fmt.Sprintf(q, name)
	if _, err := db.Exec(ddl); err != nil {
		panic(err)
	}

	if _, err := db.Exec(fmt.Sprintf("USE %s", name)); err != nil {
		panic(err)
	}
}
