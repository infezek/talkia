package controller_chats

import (
	"github.com/gofiber/fiber/v2"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/adapter"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
	"github.com/infezek/app-chat/pkg/usecase/usecase_chat"
	"github.com/infezek/app-chat/pkg/utils/middleware"
	"github.com/infezek/app-chat/pkg/utils/util_pagination"
	"github.com/infezek/app-chat/pkg/utils/util_url"

	"github.com/go-playground/validator/v10"
)

func Http(app *fiber.App, repoChat repository.RepositoryChat, repoUser repository.RepositoryUser, repoBot repository.RepositoryBot, adapterToken adapter.AdapterToken, cfg *config.Config) {
	jwt := middleware.NewAuthMiddleware(cfg.OpenIAToken)
	app.Post(util_url.New("/chats/bot/:botID"), jwt, func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		paramsUser, err := adapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		botID := c.Params("botID")
		usecase := usecase_chat.NewCreate(repoChat, repoBot, repoUser, cfg)
		output, err := usecase.Execute(usecase_chat.CreateDtoInput{
			UserID: paramsUser.UserID,
			BotID:  botID,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(output)
	})

	app.Get(util_url.New("/chats"), jwt, func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		paramsUser, err := adapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		pg := util_pagination.NewWithParams(c, util_pagination.Pagination{})
		usecase := usecase_chat.NewList(repoChat, cfg)
		output, err := usecase.Execute(entity.Pagination{
			PerPage: pg.Limit,
			Offset:  pg.Offset,
			Page:    pg.Page,
		}, paramsUser.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(output)
	})

	app.Get(util_url.New("/chats/:chatID"), jwt, func(c *fiber.Ctx) error {
		chatID := c.Params("chatID")
		bearer := c.Get("Authorization")
		paramsUser, err := adapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		usecase := usecase_chat.NewChatGet(repoChat, cfg)
		output, err := usecase.Execute(chatID, paramsUser.UserID)
		if err != nil {
			return err
		}
		return c.JSON(output)
	})
}

func validateParams(params interface{}) error {
	validate := validator.New()
	err := validate.Struct(params)
	if err != nil {
		return err
	}
	return nil
}
