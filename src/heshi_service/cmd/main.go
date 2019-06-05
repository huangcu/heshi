package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"log"
	"util"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"

	"github.com/jinzhu/configor"
)

var config = struct {
	Db struct {
		Database string
		User     string `default:"root" env:"heshi_mysql_user"`
		Password string `env:"heshi_mysql_password"`
		Port     uint   `default:"3306"`
	}
}{}

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatalf("fail to open db. err: %s", errors.GetMessage(err))
	}

	v4, _ := uuid.NewV4()
	id := v4.String()
	password := util.Encrypt("admin")
	q := fmt.Sprintf(`INSERT INTO users (id, username, password, email, user_type, invitation_code) 
	VALUES ('%s', 'admin','%s','admin@admin.com', 'admin', 'ignore')`, id, password)
	if _, err := db.Exec(q); err != nil {
		log.Fatalf("fail to create supre admin. err: %s", errors.GetMessage(err))
	}
	q = fmt.Sprintf(`INSERT INTO admins (user_id, level) VALUES ('%s',10)`, id)
	if _, err := db.Exec(q); err != nil {
		log.Fatalf("fail to create supre admin. err: %s", errors.GetMessage(err))
	}
}

// var dbCfgTeemplate = `
// db:
// 	database: heshi_dev
// 	user: root
// 	password: admin
// 	port: 3306
// `
func openDB() (*sql.DB, error) {
	// df := "database." + os.Getenv("stage") + ".yml"
	if err := configor.Load(&config, "config.yml"); err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?parseTime=true",
		config.Db.User, config.Db.Password, config.Db.Database))
	if err != nil {
		return nil, err
	}
	return db, nil
}
