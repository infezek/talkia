package controller_bots

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/adapter"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/gateway"
	"github.com/infezek/app-chat/pkg/domain/repository"
	"github.com/infezek/app-chat/pkg/usecase/chat"
	"github.com/infezek/app-chat/pkg/usecase/usecase_bot"
	"github.com/infezek/app-chat/pkg/utils/middleware"
	"github.com/infezek/app-chat/pkg/utils/util_pagination"
	"github.com/infezek/app-chat/pkg/utils/util_url"
)

func Http(
	app *fiber.App,
	repoBot repository.RepositoryBot,
	repoCategory repository.RepositoryCategory,
	repoChat repository.RepositoryChat,
	repoUser repository.RepositoryUser,
	gateway gateway.GatewayBot,
	adapterToken adapter.AdapterToken,
	adapterImage adapter.AdapterImagem,
	cfg *config.Config,
	middlerwareHandler middleware.MiddlewareHandlerInterface,
) {
	jwt := middleware.NewAuthMiddleware(cfg.OpenIAToken)
	app.Post(util_url.New("/bots"), jwt, middlerwareHandler.FindUser(), func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		paramsUser, err := adapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		var params Create
		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validateParams(params); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		usecase := usecase_bot.NewCreate(repoBot, repoCategory, cfg)
		output, err := usecase.Execute(usecase_bot.CreateDtoInput{
			UserID:      paramsUser.UserID,
			CategoryID:  params.CategoryID,
			Name:        params.Name,
			Personality: params.Personality,
			Description: params.Description,
			Location:    params.Location,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(output)
	})

	app.Put(util_url.New("/bots/:id"), jwt, middlerwareHandler.FindUser(), func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		paramsUser, err := adapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		var params Update
		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validateParams(params); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		usecase := usecase_bot.NewUpdate(repoBot, cfg)
		output, err := usecase.Execute(usecase_bot.UpdateDtoInput{
			UserID:      paramsUser.UserID,
			ID:          c.Params("id"),
			CategoryID:  params.CategoryID,
			Name:        params.Name,
			Personality: params.Personality,
			Description: params.Description,
			Location:    params.Location,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(output)
	})

	app.Put(util_url.New("/bots-image"), jwt, middlerwareHandler.FindUser(), func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		paramsUser, err := adapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		c.MultipartForm()
		usecase := usecase_bot.NewUploadImage(repoBot, adapterImage, cfg)
		files, err := uploadImage(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		out, err := usecase.Execute(usecase_bot.UploadImageDtoInput{
			BotID:  c.FormValue("bot_id"),
			UserID: paramsUser.UserID,
			Files:  files,
		})
		if err != nil {
			return err
		}
		return c.Status(http.StatusOK).JSON(map[string]string{"avatar": out.Avatar, "background": out.Background})
	})

	app.Get(util_url.New("/bots"), jwt, func(c *fiber.Ctx) error {
		pagination := util_pagination.NewWithParams(c, util_pagination.Pagination{})
		botName := c.Query("bot_name", "")
		usecase := usecase_bot.NewListByName(cfg, repoBot)
		output, err := usecase.Execute(entity.Pagination{
			PerPage: pagination.Limit,
			Offset:  pagination.Offset,
			Page:    pagination.Page,
		}, botName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(output)
	})

	app.Post(util_url.New("/bots/conversation"), jwt, func(c *fiber.Ctx) error {
		var params Conversation
		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validateParams(params); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		conversation := chat.NewBotUseCase(repoUser, repoChat, gateway, cfg)
		msg, err := conversation.ConversationTest(chat.ConversationInput{
			Name:        params.Name,
			Description: params.Description,
			Text:        params.Message,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": msg})
	})
}
