package controller_community

import (
	"github.com/gofiber/fiber/v2"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/adapter"
	"github.com/infezek/app-chat/pkg/domain/repository"
	"github.com/infezek/app-chat/pkg/usecase/usecase_community"
	"github.com/infezek/app-chat/pkg/utils/middleware"
	"github.com/infezek/app-chat/pkg/utils/util_url"
)

func Http(app *fiber.App, cfg *config.Config, repoCommunity repository.RepositoryCommunity, repoUser repository.RepositoryUser, adapterToken adapter.AdapterToken) {
	jwt := middleware.NewAuthMiddleware(cfg.OpenIAToken)
	app.Get(util_url.New("/community"), jwt, func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		paramsUser, err := adapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		usecase := usecase_community.NewList(cfg, repoCommunity, repoUser)
		output, err := usecase.Execute(paramsUser.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(output)
	})
}
