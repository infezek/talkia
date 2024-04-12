package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/database"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/repository/repository_category"
	"github.com/infezek/app-chat/pkg/repository/repository_user"
	"github.com/spf13/cobra"
)

func Seed() *cobra.Command {
	return &cobra.Command{
		Use:   "seed",
		Short: "Seed data",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.New("chat")
			if err != nil {
				panic(err)
			}
			db, err := database.New(cfg)
			if err != nil {
				panic(err)
			}
			repoCategory := repository_category.New(db)
			for _, category := range categories {
				err := repoCategory.Create(category)
				if err != nil {
					panic(err)
				}
			}

			repoUser := repository_user.New(db)
			for _, user := range users {
				err := repoUser.Create(user)
				if err != nil {
					panic(err)
				}
			}
		},
	}
}

var (
	err         error = nil
	userID_1, _       = uuid.Parse("018ea5d2-7a82-7367-98a8-7aeb13495684")
	userID_2, _       = uuid.Parse("018ea5d8-f5bb-71bf-94b0-c50922044c74")
	userID_3, _       = uuid.Parse("018ea5d9-09dc-7d03-8739-a0373e7a4cd6")

	users = []entity.User{
		{
			ID:        userID_1,
			Username:  "Ezequiel Lopes",
			Password:  "password1",
			Email:     "ezequiel@gmail.com",
			Platform:  "android",
			Gender:    &entity.GenderMasculine,
			Location:  "São Paulo",
			Language:  entity.LanguagePortuguese,
			CreatedAt: time.Now(),
		},
		{
			ID:        userID_2,
			Username:  "Maria Eduarda",
			Password:  "password2",
			Email:     "maria.eduarda@gmail.com",
			Platform:  "ios",
			Gender:    nil,
			Location:  "Rio de Janeiro",
			Language:  entity.LanguageSpanish,
			CreatedAt: time.Now(),
		},
		{
			ID:        userID_3,
			Username:  "João Pedro",
			Password:  "password3",
			Email:     "joao.pedro@gmail.com",
			Platform:  "android",
			Gender:    &entity.GenderMasculine,
			Location:  "Europa",
			Language:  entity.LanguageEnglish,
			CreatedAt: time.Now(),
		},
	}

	categoryID_1, _ = uuid.Parse("018ea5da-0796-7219-bd90-883cae348caf")
	categoryID_2, _ = uuid.Parse("018ea5da-1e95-783c-90ba-eabc282093ae")
	categoryID_3, _ = uuid.Parse("018ea5da-3304-7c02-a4b4-74e6e0796965")
	categoryID_4, _ = uuid.Parse("018ea5da-4b3d-7f3b-8b9b-7b3e7e7a4b3d")
	categoryID_5, _ = uuid.Parse("018ea5db-25fe-736c-85b9-2f72c40b942a")

	categories = []entity.Category{
		{
			ID:        categoryID_1,
			Name:      "Anime",
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:        categoryID_2,
			Name:      "Harry Potter",
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:        categoryID_3,
			Name:      "The Lord of the Rings",
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:        categoryID_4,
			Name:      "Desenho Animado",
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:        categoryID_5,
			Name:      "Star Wars",
			Active:    true,
			CreatedAt: time.Now(),
		},
	}
)
