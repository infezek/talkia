package usecase_bot

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type Create struct {
	RepoBot      repository.RepositoryBot
	RepoCategory repository.RepositoryCategory
	Cfg          *config.Config
}

func NewCreate(repoBot repository.RepositoryBot, repoCategory repository.RepositoryCategory, cfg *config.Config) *Create {
	return &Create{
		RepoBot:      repoBot,
		RepoCategory: repoCategory,
		Cfg:          cfg,
	}
}

func (u *Create) Execute(params CreateDtoInput) (CreateDtoOutput, error) {
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		return CreateDtoOutput{}, fmt.Errorf("invalid user id")
	}
	categoryID, err := uuid.Parse(params.CategoryID)
	if err != nil {
		return CreateDtoOutput{}, fmt.Errorf("invalid category id")
	}
	category, err := u.RepoCategory.FindByID(categoryID)
	if err != nil {
		return CreateDtoOutput{}, fmt.Errorf("category not found")
	}
	bot := entity.NewBot(
		uuid.Nil,
		userID,
		category.ID,
		params.Name,
		params.Personality,
		params.Description,
		"",
		"",
		params.Location,
		time.Now(),
		time.Now(),
	)
	if err := u.RepoBot.Create(*bot); err != nil {
		return CreateDtoOutput{}, err
	}
	return CreateDtoOutput{
		ID:          bot.ID.String(),
		Name:        bot.Name,
		CategoryID:  bot.CategoryID.String(),
		Personality: bot.Personality,
		Description: bot.Description,
		Location:    bot.Location,
		CreatedAt:   bot.CreatedAt,
		UpdatedAt:   bot.UpdatedAt,
	}, nil
}
