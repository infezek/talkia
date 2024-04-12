package test

import (
	"database/sql"

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

func (d *DatabaseTest) Trucate() {
	d.DB.Exec("DELETE FROM user_like_bot")
	d.DB.Exec("DELETE FROM users_categories")
	d.DB.Exec("DELETE FROM preferences")
	d.DB.Exec("DELETE FROM messages")
	d.DB.Exec("DELETE FROM chats")
	d.DB.Exec("DELETE FROM messages")
	d.DB.Exec("DELETE FROM categories")
	d.DB.Exec("DELETE FROM bots")
	d.DB.Exec("DELETE FROM users")
}
