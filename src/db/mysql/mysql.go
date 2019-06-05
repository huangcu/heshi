package mysql

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"log"
	"os"
	"util"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/configor"
	uuid "github.com/satori/go.uuid"
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
	tableName := []string{"users", "agents", "admins"}
	tableName = append(tableName, "diamonds", "jewelrys", "small_diamonds", "gems")
	tableName = append(tableName, "promotions", "orders", "shopping_cart", "currency_exchange_rates")
	tableName = append(tableName, "invitation_codes", "wechat_users", "discounts", "level_rate_rules")
	tableName = append(tableName, "interested_items", "user_using_records", "user_active_record")
	tableName = append(tableName, "suppliers", "price_settings_universal", "action_logs")
	tableName = append(tableName, "messages", "wechat_messages", "historys", "product_stock_handle_records")
	tableDdl := []string{userDdl, agentDdl, adminDdl}
	tableDdl = append(tableDdl, diamondDdl, jewelryDdl, smallDiamondDdl, gemDdl)
	tableDdl = append(tableDdl, promotionDdl, orderDdl, shoppingCart, currencyExchangeRateDdl)
	tableDdl = append(tableDdl, invitationCodeDdl, wechatUserDdl, discountDdl, levelRateDdl)
	tableDdl = append(tableDdl, interestedItemDdl, userUsingRecordDdl, userActiveRecordDdl)
	tableDdl = append(tableDdl, supplierDdl, priceSettingUniversalDdl, actionLogDdl)
	tableDdl = append(tableDdl, messageDdl, wechatMessageDdl, historyTrackDdl, productStockHandleHistoryDdl)
	if len(tableName) != len(tableDdl) {
		return nil, errors.New("db DDL number is not a match to table number")
	}

	for i := 0; i < len(tableName); i++ {
		table := tableName[i]
		ddl := tableDdl[i]
		if _, err := db.Exec(ddl); err != nil {
			return nil, fmt.Errorf("fail to create table %s with %s; err: %s", table, ddl, errors.GetMessage(err))
		}
	}
	db.Close()

	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?parseTime=true",
		Config.Db.User, Config.Db.Password, Config.Db.Database))
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(0)

	if os.Getenv("STAGE") == "dev" {
		if err := createDefaultDevUser(db); err != nil {
			log.Fatal(err.Error())
		}
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

func createDefaultDevUser(db *sql.DB) error {
	if exist, err := isUserExistByEmail(db, "agent@dev.com"); (err == nil) && !exist {
		v4, _ := uuid.NewV4()
		id := v4.String()
		qAgent := fmt.Sprintf(`INSERT INTO users (id, username, password, email, user_type, invitation_code) 
			VALUES ('%s', 'agent','%s','agent@dev.com', 'AGENT', 'ignore_agent')`, id, util.Encrypt("agent"))
		if _, err := db.Exec(qAgent); err != nil {
			return err
		}
		if _, err := db.Exec(fmt.Sprintf(`INSERT INTO agents (user_id, level) VALUES ('%s',10)`, id)); err != nil {
			return err
		}
	}

	if exist, err := isUserExistByEmail(db, "admin@dev.com"); (err == nil) && !exist {
		v4, _ := uuid.NewV4()
		id := v4.String()
		qAdmin := fmt.Sprintf(`INSERT INTO users (id, username, password, email, user_type, invitation_code) 
			VALUES ('%s', 'admin','%s','admin@dev.com', 'ADMIN', 'ignore_admin')`, id, util.Encrypt("admin"))
		if _, err := db.Exec(qAdmin); err != nil {
			return err
		}
		if _, err := db.Exec(fmt.Sprintf(`INSERT INTO admins (user_id, level) VALUES ('%s',10)`, id)); err != nil {
			return err
		}
	}

	if exist, err := isUserExistByEmail(db, "customer@dev.com"); (err == nil) && !exist {
		v4, _ := uuid.NewV4()
		id := v4.String()
		qCustomer := fmt.Sprintf(`INSERT INTO users (id, username, password, email, user_type, invitation_code) 
			VALUES ('%s', 'customer','%s','customer@dev.com', 'CUSTOMER', 'ignore_customer')`, id, util.Encrypt("customer"))
		if _, err := db.Exec(qCustomer); err != nil {
			return err
		}
	}
	return nil
}

func isUserExistByEmail(db *sql.DB, email string) (bool, error) {
	var count int
	if err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM users WHERE email='%s'", email)).Scan(&count); err != nil {
		return false, err
	}
	return count == 1, nil
}
