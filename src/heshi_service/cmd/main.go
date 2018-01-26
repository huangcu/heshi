package main

import (
	"database/sql"
	"fmt"
	"log"
	"util"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/configor"
	uuid "github.com/satori/go.uuid"
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
		log.Fatalf("fail to open db. err: %s", err.Error())
	}

	id := uuid.NewV4().String()
	password := util.Encrypt("admin")
	q := `INSERT INTO users (id, username, password, email, user_type, invitation_code) 
	VALUES (?, 'admin',?,'admin@admin.com', 'admin', 'ignore')`
	if _, err := db.Exec(q, id, password); err != nil {
		log.Fatalf("fail to create supre admin. err: %s", err.Error())
	}
	q = `INSERT INTO admins (user_id, level) VALUES (?,10)`
	if _, err := db.Exec(q, id); err != nil {
		log.Fatalf("fail to create supre admin. err: %s", err.Error())
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
