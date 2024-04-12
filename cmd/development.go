package main

import (
	"fmt"
	"time"

	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/database"
	"github.com/spf13/cobra"
)

type AA struct {
	name string
	tt   time.Time
}

func development(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:  "dev",
		Long: "Development commands",
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println(cfg)
			db, err := database.New(cfg)
			if err != nil {
				panic(err)
			}
			res, err := db.Exec("DROP DATABASE IF EXISTS app_chat")
			if err != nil {
				panic(err)
			}
			resRows, err := res.RowsAffected()
			fmt.Println("[1]", resRows)
			res, err = db.Exec("CREATE DATABASE IF NOT EXISTS app_chat")
			if err != nil {
				panic(err)
			}
			resRows, err = res.RowsAffected()
			fmt.Println("[2]", resRows)

		},
	}
}
