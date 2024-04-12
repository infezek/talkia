package controller_users

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/adapter"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
	"github.com/infezek/app-chat/pkg/usecase/usecase_bot"
	"github.com/infezek/app-chat/pkg/usecase/usecase_user"
	"github.com/infezek/app-chat/pkg/utils/middleware"
	"github.com/infezek/app-chat/pkg/utils/util_pagination"
	"github.com/infezek/app-chat/pkg/utils/util_url"

	"github.com/go-playground/validator/v10"
)

func Http(app *fiber.App,
	cfg *config.Config,
	repoUser repository.RepositoryUser,
	adapterToken adapter.AdapterToken,
	adapterImage adapter.AdapterImagem,
	logger adapter.AdapterLogger,
	middlerwareHandler middleware.MiddlewareHandlerInterface,
) {
	app.Post(util_url.New("/users"), func(c *fiber.Ctx) error {
		var params CreateUserParams
		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validateHandler(params); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		usecase := usecase_user.CreateUserUseCase(repoUser, cfg)
		out, err := usecase.Execute(usecase_user.CreateUserDtoInput{
			Username: params.Username,
			Email:    params.Email,
			Password: params.Password,
			Platform: params.Platform,
			Location: params.Location,
			Language: params.Language,
		})
		if err != nil {
			return err
		}
		usecaseLogin := usecase_user.NewLoginUseCase(cfg, repoUser, adapterToken)
		token, err := usecaseLogin.Execute(params.Email)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"token":      token,
			"username":   out.Username,
			"avatar_url": out.AvatarURL,
			"bucket_url": cfg.BucketImagesURL,
		})
	})
	jwt := middleware.NewAuthMiddleware(cfg.OpenIAToken)
	app.Post(util_url.New("/users/categories"), jwt, middlerwareHandler.FindUser(), func(c *fiber.Ctx) error {
		logger.Info("users/categories")
		bearer := c.Get("Authorization")
		paramsUser, err := adapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		logger.AddParam("user_id", paramsUser.UserID)
		logger.AddParam("email", paramsUser.Email)
		logger.Info("params user")
		var param []string
		if err := c.BodyParser(&param); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		usecase := usecase_user.NewAddCategory(cfg, repoUser)
		err = usecase.Execute(usecase_user.AddCategoryDtoInput{
			UserID:     paramsUser.UserID,
			Categories: param,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusCreated).Send(nil)
	})

	app.Post(util_url.New("/like/bot/:botID"), jwt, middlerwareHandler.FindUser(), func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		paramsUser, err := adapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		botID := c.Params("botID")
		usecase := usecase_bot.NewLike(cfg, repoUser)
		if err := usecase.Execute(usecase_bot.LikeInput{
			UserID: paramsUser.UserID,
			BotID:  botID,
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusCreated).Send(nil)
	})

	app.Put(util_url.New("/users"), jwt, middlerwareHandler.FindUser(), func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		paramsUser, err := adapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		var params UpdateUserParams
		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validateHandler(params); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		usecase := usecase_user.UpdateUseCase(cfg, repoUser)
		if err := usecase.Execute(usecase_user.UpdateUserDtoInput{
			UserID:   paramsUser.UserID,
			Username: params.Username,
			Language: params.Language,
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusCreated).Send(nil)
	})

	app.Post(util_url.New("/login"), func(c *fiber.Ctx) error {
		var params LoginParams
		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validateHandler(params); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		usecase := usecase_user.NewLoginUseCase(cfg, repoUser, adapterToken)
		token, err := usecase.Execute(params.Email)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"token": token})
	})

	app.Put(util_url.New("/users/image"), jwt, middlerwareHandler.FindUser(), func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		paramsUser, err := adapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		c.MultipartForm()
		usecase := usecase_user.NewUploadImage(cfg, repoUser, adapterImage)
		file, err := uploadImage(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := usecase.Execute(usecase_user.UploadImageDtoInput{
			UserID: paramsUser.UserID,
			Avatar: file,
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusOK).JSON(map[string]string{"message": "Image uploaded successfully!"})
	})

	app.Get(util_url.New("/users/bots"), jwt, middlerwareHandler.FindUser(), func(c *fiber.Ctx) error {
		pagination := util_pagination.NewWithParams(c, util_pagination.Pagination{})
		bearer := c.Get("Authorization")
		paramsUser, err := adapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		usecase := usecase_user.NewMyBots(repoUser, cfg)
		output, err := usecase.Execute(entity.Pagination{
			PerPage: pagination.Limit,
			Offset:  pagination.Offset,
			Page:    pagination.Page,
		}, paramsUser.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(output)
	})

}

type LoginParams struct {
	Email string `json:"email" validate:"required,email"`
}

type CreateUserParams struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Platform string `json:"platform" validate:"oneof=android ios"`
	Location string `json:"location" validate:"required"`
	Language string `json:"language" validate:"oneof=english portuguese spanish"`
}

type UpdateUserParams struct {
	Username string `json:"username" validate:"required"`
	Language string `json:"language" validate:"oneof=english portuguese spanish"`
}

func validateHandler(params interface{}) error {
	validate := validator.New()
	err := validate.Struct(params)
	if err != nil {
		return err
	}
	return nil
}
