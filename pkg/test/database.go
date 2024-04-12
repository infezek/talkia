package test

import (
	"database/sql"
	"fmt"

	"github.com/infezek/app-chat/pkg/config"
)

type DatabaseTest struct {
	DB  *sql.DB
	Cfg *config.Config
}

func ConfigTest() *config.Config {
	return &config.Config{
		DBHost:     "172.17.0.1",
		DBPort:     "3306",
		DBUser:     "root",
		DBPassword: "123",
		DBName:     "app_chat_test",
	}
}

func NewDatabaseTest(db *sql.DB, cfg *config.Config) *DatabaseTest {
	return &DatabaseTest{
		DB:  db,
		Cfg: cfg,
	}
}

func (d *DatabaseTest) Delete() {
	result, err := d.DB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", d.Cfg.DBName))
	if err != nil {
		fmt.Println("err:", err.Error())
	}
	fmt.Println(result)
}

func (d *DatabaseTest) Create() {
	result, err := d.DB.Exec(fmt.Sprintf("DROP DATABASE IF NOT EXISTS  %s", d.Cfg.DBName))
	if err != nil {
		fmt.Println("err:", err.Error())
	}
	fmt.Println(result)
}
