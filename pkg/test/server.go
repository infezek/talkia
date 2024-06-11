package test

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/infezek/app-chat/pkg/adapter"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/controller/controller_bots"
	"github.com/infezek/app-chat/pkg/controller/controller_categories"
	"github.com/infezek/app-chat/pkg/controller/controller_chats"
	"github.com/infezek/app-chat/pkg/controller/controller_community"
	"github.com/infezek/app-chat/pkg/controller/controller_home"
	"github.com/infezek/app-chat/pkg/controller/controller_users"
	"github.com/infezek/app-chat/pkg/database"
	"github.com/infezek/app-chat/pkg/domain/gateway"
	"github.com/infezek/app-chat/pkg/domain/repository"
	"github.com/infezek/app-chat/pkg/process_error"
	"github.com/infezek/app-chat/pkg/repository/repository_bot"
	"github.com/infezek/app-chat/pkg/repository/repository_category"
	"github.com/infezek/app-chat/pkg/repository/repository_chat"
	"github.com/infezek/app-chat/pkg/repository/repository_community"
	"github.com/infezek/app-chat/pkg/repository/repository_user"
	"github.com/infezek/app-chat/pkg/utils/middleware"
)

func Implementations(cfg *config.Config) Params {
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

	return Params{
		DB:            db,
		RepoUser:      repoUser,
		RepoCategory:  repoCategory,
		RepoBot:       repoBot,
		RepoChat:      repoChat,
		RepoCommunity: repoCommunity,
		AdapterImage:  adapterImage,
		AdapterToken:  adapterToken,
		Cfg:           cfg,
	}

}

type Params struct {
	RepoUser      repository.RepositoryUser
	RepoCategory  repository.RepositoryCategory
	RepoBot       repository.RepositoryBot
	RepoChat      repository.RepositoryChat
	RepoCommunity repository.RepositoryCommunity
	GatewayBot    gateway.GatewayBot
	AdapterImage  *adapter.AdapterImagem
	AdapterToken  *adapter.AdapterToken
	Cfg           *config.Config
	DB            *sql.DB
}

func Server(params Params) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: process_error.Handler,
	})

	middleware := middleware.NewMiddlewareHandler(params.RepoUser, params.AdapterToken)
	controller_categories.Http(app, params.RepoCategory, params.Cfg)
	controller_users.Http(app, params.Cfg, params.RepoUser, params.AdapterToken, params.AdapterImage, middleware)
	controller_bots.Http(app, params.RepoBot, params.RepoCategory, params.RepoChat, params.RepoUser, params.GatewayBot, params.AdapterToken, params.AdapterImage, params.Cfg, middleware)
	controller_chats.Http(app, params.RepoChat, params.RepoUser, params.RepoBot, params.AdapterToken, params.Cfg)
	controller_bots.WebSocket(app, params.AdapterToken)
	controller_community.Http(app, params.Cfg, params.RepoCommunity, params.RepoUser, params.AdapterToken)
	controller_home.Http(app)

	return app
}
