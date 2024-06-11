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
		userToken, err := usecaseLogin.Execute(params.Email)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"id":         out.ID,
			"token":      userToken.Token,
			"username":   out.Username,
			"avatar_url": out.AvatarURL,
			"bucket_url": cfg.BucketImagesURL,
			"email":      out.Email,
			"platform":   out.Platform,
			"gender":     out.Gender,
			"location":   out.Location,
			"language":   out.Language,
			"created_at": out.CreatedAt,
		})
	})
	jwt := middleware.NewAuthMiddleware(cfg.OpenIAToken)
	app.Post(util_url.New("/users/categories"), jwt, middlerwareHandler.FindUser(), func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		paramsUser, err := adapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

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
		user, err := usecase.Execute(usecase_user.UpdateUserDtoInput{
			UserID:   paramsUser.UserID,
			Username: params.Username,
			Language: params.Language,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusOK).JSON(user)
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
		userToken, err := usecase.Execute(params.Email)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"id":         userToken.ID,
			"token":      userToken.Token,
			"username":   userToken.Username,
			"avatar_url": userToken.AvatarURL,
			"bucket_url": cfg.BucketImagesURL,
			"email":      userToken.Email,
			"platform":   userToken.Platform,
			"gender":     userToken.Gender,
			"location":   userToken.Location,
			"language":   userToken.Language,
			"created_at": userToken.CreatedAt,
		})
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
		out, err := usecase.Execute(usecase_user.UploadImageDtoInput{
			UserID: paramsUser.UserID,
			Avatar: file,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusOK).JSON(map[string]string{"avatar_url": out.AvatarURL})
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
