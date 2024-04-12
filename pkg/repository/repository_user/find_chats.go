package repository_user

import (
	"context"

	"github.com/google/uuid"
	repositoryDB "github.com/infezek/app-chat/db/repository"
	"github.com/infezek/app-chat/pkg/domain/entity"
)

func findChatsByUserID(repo *repositoryDB.Queries, userID uuid.UUID) ([]entity.Chat, error) {
	chatsRepo, err := repo.FindChatsByUserID(context.Background(), userID.String())
	if err != nil {
		return []entity.Chat{}, err
	}
	chats := []entity.Chat{}
	for _, chat := range chatsRepo {
		chatID, err := uuid.Parse(chat.ID)
		if err != nil {
			return []entity.Chat{}, err
		}
		botID, err := uuid.Parse(chat.BotID)
		if err != nil {
			return []entity.Chat{}, err
		}
		chats = append(chats, entity.Chat{
			ID:              chatID,
			BotID:           botID,
			Messages:        []entity.Message{},
			UserPreferences: []entity.Preference{},
			UserID:          userID,
			Bot:             entity.Bot{},
		})
	}
	return chats, nil
}
