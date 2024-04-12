package migrate

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql" // Importa o driver do MySQL
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/infezek/app-chat/pkg/config"
)

func Run(cfg *config.Config) {
	file := "file://db/migrations"
	if cfg.Environment == "debug" {
		file = "file://../db/migrations"
	}
	m, err := migrate.New(
		file,
		fmt.Sprintf(
			"mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		))
	if err != nil {
		log.Fatal("ss", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("s*", err)
	}
}
