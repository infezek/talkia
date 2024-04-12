package repository_user

import (
	"context"

	"github.com/google/uuid"
	repositoryDB "github.com/infezek/app-chat/db/repository"
	"github.com/infezek/app-chat/pkg/domain/entity"
)

func findBotsByUserID(repo *repositoryDB.Queries, userID uuid.UUID) ([]entity.Bot, error) {
	botsRepo, err := repo.FindBotsByUserID(context.Background(), userID.String())
	if err != nil {
		return []entity.Bot{}, err
	}
	bots := []entity.Bot{}
	for _, bot := range botsRepo {
		botID, err := uuid.Parse(bot.ID)
		if err != nil {
			return []entity.Bot{}, err
		}
		bots = append(bots, entity.Bot{
			ID:            botID,
			UserID:        userID,
			Name:          bot.Name,
			CategoryID:    uuid.MustParse(bot.CategoryID),
			Personality:   bot.Personality,
			Published:     bot.Published,
			Active:        bot.Active,
			Description:   bot.Description,
			AvatarURL:     bot.AvatarUrl,
			BackgroundURL: bot.BackgroundUrl,
			Location:      bot.Location,
			Like:          bot.Likes,
			CreatedAt:     bot.CreatedAt,
			UpdatedAt:     bot.UpdatedAt,
		})
	}
	return bots, nil
}
