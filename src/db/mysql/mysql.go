package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/configor"
)

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
	tableName := []string{"users", "agents", "admins", "diamonds", "jewelrys", "small_diamonds", "promotions", "suppliers", "appointments", "orders", "currency_exchange_rates"}
	tableName = append(tableName, "invitation_codes", "wechat_users", "discounts", "configs")
	tableDdl := []string{userDdl, agentDdl, adminDdl, diamondDdl, jewelryDdl, smallDiamondDdl, promotionDdl, supplierDdl, appointmentDdl, orderDdl, currencyExchangeRateDdl}
	tableDdl = append(tableDdl, invitationCodeDdl, wechatUserDdl, discountDdl, configDdl)
	if len(tableName) != len(tableDdl) {
		return nil, errors.New("db DDL number is not a match to table number")
	}

	for i := 0; i < len(tableName); i++ {
		table := tableName[i]
		ddl := tableDdl[i]
		if _, err := db.Exec(ddl); err != nil {
			return nil, fmt.Errorf("fail to create table %s with %s; err: %s", table, ddl, err.Error())
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
