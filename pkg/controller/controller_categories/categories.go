package controller_categories

import (
	"github.com/gofiber/fiber/v2"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
	"github.com/infezek/app-chat/pkg/usecase/usecase_category"
	"github.com/infezek/app-chat/pkg/utils/middleware"
	"github.com/infezek/app-chat/pkg/utils/util_pagination"
	"github.com/infezek/app-chat/pkg/utils/util_url"

	"github.com/go-playground/validator/v10"
)

func Http(app *fiber.App, repoCategory repository.RepositoryCategory, cfg *config.Config) {
	jwt := middleware.NewAuthMiddleware(cfg.OpenIAToken)
	app.Get(util_url.New("/categories"), jwt, func(c *fiber.Ctx) error {
		pagination := util_pagination.NewWithParams(c, util_pagination.Pagination{})
		usecase := usecase_category.NewList(repoCategory, cfg)
		output, err := usecase.Execute(entity.Pagination{
			PerPage: pagination.Limit,
			Offset:  pagination.Offset,
			Page:    pagination.Page,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(output)
	})

	app.Post(util_url.New("/categories"), jwt, func(c *fiber.Ctx) error {
		var params CategoriesParams
		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validateCategoriesParams(params); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		usecase := usecase_category.NewCreate(repoCategory, cfg)
		output, err := usecase.Execute(usecase_category.CreateDtoInput{
			Name:   params.Name,
			Active: params.Active,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(output)
	})

	app.Put(util_url.New("/categories"), jwt, func(c *fiber.Ctx) error {
		var params CategoriesUpdateParams
		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validateCategoriesUpdateParams(params); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		usecase := usecase_category.NewUpdate(repoCategory, cfg)
		output, err := usecase.Execute(usecase_category.UpdateDtoInput{
			ID:     params.ID,
			Name:   params.Name,
			Active: params.Active,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(output)
	})
}

type CategoriesUpdateParams struct {
	ID string `json:"id" validate:"required"`
	CategoriesParams
}

func validateCategoriesUpdateParams(params CategoriesUpdateParams) error {
	validate := validator.New()
	err := validate.Struct(params)
	if err != nil {
		return err
	}
	return nil
}

type CategoriesParams struct {
	Name   string `json:"name" validate:"required"`
	Active bool   `json:"active" validate:"boolean"`
}

func validateCategoriesParams(params CategoriesParams) error {
	validate := validator.New()
	err := validate.Struct(params)
	if err != nil {
		return err
	}
	return nil
}
