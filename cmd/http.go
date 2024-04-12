package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/infezek/app-chat/db/migrate"
	"github.com/infezek/app-chat/pkg/adapter"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/controller/controller_bots"
	"github.com/infezek/app-chat/pkg/controller/controller_categories"
	"github.com/infezek/app-chat/pkg/controller/controller_chats"
	"github.com/infezek/app-chat/pkg/controller/controller_community"
	"github.com/infezek/app-chat/pkg/controller/controller_users"
	"github.com/infezek/app-chat/pkg/database"
	"github.com/infezek/app-chat/pkg/process_error"
	"github.com/infezek/app-chat/pkg/repository/repository_bot"
	"github.com/infezek/app-chat/pkg/repository/repository_category"
	"github.com/infezek/app-chat/pkg/repository/repository_chat"
	"github.com/infezek/app-chat/pkg/repository/repository_community"
	"github.com/infezek/app-chat/pkg/repository/repository_user"
	"github.com/infezek/app-chat/pkg/utils/middleware"
	"github.com/infezek/app-chat/pkg/utils/util_url"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func http(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Start HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Infof("Starting HTTP server on port %s", os.Getenv("PORT"))
			app := fiber.New(fiber.Config{
				ErrorHandler: process_error.Handler,
			})
			configNewRelic := Config{
				License: cfg.NewRelicLicense,
				AppName: cfg.NewRelicAppName,
				Enabled: cfg.NewRelicEnabled,
			}
			nr, err := newrelic.NewApplication(
				newrelic.ConfigAppName(configNewRelic.AppName),
				newrelic.ConfigLicense(configNewRelic.License),
			)
			if err != nil {
				log.Fatalf("Erro ao inicializar o New Relic: %v", err)
			}
			app.Use(fiberNewRelic(configNewRelic))
			app.Use(logger.New(logger.Config{
				Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
			}))

			app.Get(util_url.New("/metrics"), monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
			migrate.Run(cfg)
			db, err := database.New(cfg)
			if err != nil {
				panic(err)
			}

			repoUser := repository_user.New(db)
			repoCategory := repository_category.New(db)
			repoBot := repository_bot.New(db)
			repoChat := repository_chat.New(db)
			repoCommunity := repository_community.New(db)
			adapterImage := adapter.NewImage(cfg)
			adapterToken := adapter.NewToken(cfg)

			l := logrus.NewEntry(logrus.New())
			adapterLogger, err := adapter.NewLogger(l, nr)
			if err != nil {
				fmt.Println("error ao criar logger", err.Error())
			}
			middleware := middleware.NewMiddlewareHandler(repoUser, adapterToken)
			controller_categories.Http(app, repoCategory, cfg)
			controller_users.Http(app, cfg, repoUser, adapterToken, adapterImage, adapterLogger, middleware)
			controller_bots.Http(app, repoBot, repoCategory, adapterToken, adapterImage, cfg, middleware)
			controller_chats.Http(app, repoChat, repoUser, repoBot, adapterToken, cfg)
			controller_bots.WebSocket(app, adapterToken)
			controller_community.Http(app, cfg, repoCommunity, repoUser, adapterToken)

			app.Get(util_url.New("/"), func(c *fiber.Ctx) error {
				return c.JSON(map[string]string{"message": "ok"})
			})

			port := os.Getenv("PORT")
			if port == "" {
				port = "3000"
			}
			if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
				panic(err)
			}
		},
	}
}
