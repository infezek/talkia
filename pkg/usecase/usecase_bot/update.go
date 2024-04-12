package usecase_bot

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type Update struct {
	RepoBot repository.RepositoryBot
	Cfg     *config.Config
}

func NewUpdate(repoBot repository.RepositoryBot, cfg *config.Config) *Create {
	return &Create{
		RepoBot: repoBot,
		Cfg:     cfg,
	}
}

func (u *Update) Execute(params UpdateDtoInput) (UpdateDtoOutput, error) {
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		return UpdateDtoOutput{}, fmt.Errorf("invalid user id")
	}
	botID, err := uuid.Parse(params.ID)
	if err != nil {
		return UpdateDtoOutput{}, fmt.Errorf("invalid bot id")
	}
	bot, err := u.RepoBot.FindByID(botID)
	if err != nil {
		return UpdateDtoOutput{}, err
	}
	if bot.UserID.String() != userID.String() {
		return UpdateDtoOutput{}, fmt.Errorf("unauthorized")
	}
	if bot.Published {
		return UpdateDtoOutput{}, fmt.Errorf("cannot update published bot")
	}
	bot.Update(
		params.Name,
		params.Personality,
		params.Description,
		params.AvatarURL,
		params.BackgroundURL,
		params.Location,
	)
	if err := u.RepoBot.Update(bot); err != nil {
		return UpdateDtoOutput{}, err
	}
	return UpdateDtoOutput{
		ID:          bot.ID.String(),
		Name:        bot.Name,
		Personality: bot.Personality,
		Description: bot.Description,
		Location:    bot.Location,
		CreatedAt:   bot.CreatedAt,
		UpdatedAt:   bot.UpdatedAt,
	}, nil
}
