package repository_chat

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	repositoryDB "github.com/infezek/app-chat/db/repository"
	"github.com/infezek/app-chat/pkg/domain/entity"
)

type RepositoryChat struct {
	db   *sql.DB
	repo *repositoryDB.Queries
}

func New(db *sql.DB) *RepositoryChat {
	repo := repositoryDB.New(db)
	return &RepositoryChat{db, repo}
}

func (rc *RepositoryChat) Create(chat entity.Chat) error {
	return rc.repo.CreateChat(context.Background(), repositoryDB.CreateChatParams{
		ID:     chat.ID.String(),
		UserID: chat.UserID.String(),
		BotID:  chat.Bot.ID.String(),
	})
}

func (rc *RepositoryChat) SaveMessage(chat entity.Chat, message entity.Message) error {
	if message.Who == entity.MessageUser {
		return rc.repo.CreateMessage(context.Background(), repositoryDB.CreateMessageParams{
			ID:      message.ID.String(),
			ChatID:  chat.ID.String(),
			UserID:  chat.UserID.String(),
			Who:     repositoryDB.MessagesWhoUser,
			Message: message.Message,
		})
	}
	return rc.repo.CreateMessage(context.Background(), repositoryDB.CreateMessageParams{
		ID:      message.ID.String(),
		ChatID:  chat.ID.String(),
		UserID:  chat.UserID.String(),
		Who:     repositoryDB.MessagesWhoSystem,
		Message: message.Message,
	})
}

func (rc *RepositoryChat) SavePreferences(chat entity.Chat, message entity.Message) error {
	return rc.repo.CreatePreference(context.Background(), repositoryDB.CreatePreferenceParams{
		ChatID:          chat.ID.String(),
		UserID:          chat.UserID.String(),
		PreferenceKey:   string(entity.PreferenceUser),
		PreferenceValue: message.Message,
	})
}

func (rc *RepositoryChat) Update(chat entity.Chat) error {
	return nil
}

func (rc *RepositoryChat) Desactive(id uuid.UUID) error {
	return nil
}

func (rc *RepositoryChat) FindByID(id uuid.UUID) (entity.Chat, error) {
	chat, err := rc.repo.GetChatByID(context.Background(), id.String())
	if err != nil {
		return entity.Chat{}, err
	}
	chatID, err := uuid.Parse(chat.ID)
	if err != nil {
		return entity.Chat{}, err
	}
	userID, err := uuid.Parse(chat.UserID)
	if err != nil {
		return entity.Chat{}, err
	}
	messages, err := rc.repo.GetMessagesByChatID(context.Background(), id.String())
	var msg []entity.Message
	for _, m := range messages {
		msgID, err := uuid.Parse(m.ID)
		if err != nil {
			return entity.Chat{}, err
		}
		msg = append(msg, entity.Message{
			ID:      msgID,
			Who:     entity.Who(m.Who),
			Message: m.Message,
			ChatID:  chatID,
			UserID:  userID,
		})
	}
	bot, err := rc.repo.GetBotByID(context.Background(), chat.BotID)
	userPreferences, err := rc.repo.GetPreferencesByChatID(context.Background(), chat.ID)
	if err != nil {
		return entity.Chat{}, err
	}
	var preferences []entity.Preference
	for _, p := range userPreferences {
		preferenceID, err := uuid.Parse(p.ID)
		if err != nil {
			return entity.Chat{}, err
		}
		preferences = append(preferences, entity.Preference{
			ID:              preferenceID,
			PreferenceKey:   entity.PreferenceKey(p.PreferenceKey),
			PreferenceValue: p.PreferenceValue,
			ChatID:          chatID,
			UserID:          userID,
		})
	}
	botID, err := uuid.Parse(bot.ID)
	if err != nil {
		return entity.Chat{}, err
	}
	categoryID, err := uuid.Parse(bot.CategoryID)
	if err != nil {
		return entity.Chat{}, err
	}
	return entity.Chat{
		ID:              chatID,
		UserID:          userID,
		UserPreferences: preferences,
		Bot: entity.Bot{
			ID:            botID,
			Name:          bot.Name,
			Personality:   bot.Personality,
			Description:   bot.Description,
			AvatarURL:     bot.AvatarUrl,
			CategoryID:    categoryID,
			UserID:        userID,
			BackgroundURL: bot.BackgroundUrl,
			Location:      bot.Location,
			Published:     bot.Published,
			Active:        bot.Active,
			Like:          bot.Likes,
			CreatedAt:     bot.CreatedAt,
			UpdatedAt:     bot.UpdatedAt,
		},
		Messages: msg,
	}, nil
}

func (rc *RepositoryChat) FindByUserID(pagintion entity.Pagination, userID uuid.UUID) (chats []entity.Chat, total int64, err error) {
	repoChats, err := rc.repo.ListChatsByUserID(context.Background(), repositoryDB.ListChatsByUserIDParams{
		UserID: userID.String(),
		Limit:  pagintion.PerPage,
		Offset: pagintion.Offset,
	})
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if err != nil && err == sql.ErrNoRows {
		return
	}
	for _, chat := range repoChats {
		chatID, err := uuid.Parse(chat.ID)
		if err != nil {
			return chats, total, err
		}
		botID, err := uuid.Parse(chat.BotID)
		if err != nil {
			return chats, total, err
		}
		bot, err := rc.repo.GetBotByID(context.Background(), chat.BotID)
		if err != nil {
			return chats, total, err
		}
		botCategory, err := uuid.Parse(bot.CategoryID)
		if err != nil {
			return chats, total, err
		}
		chats = append(chats, entity.Chat{
			ID:     chatID,
			UserID: userID,
			BotID:  botID,
			Bot: entity.Bot{
				ID:            botID,
				Name:          bot.Name,
				UserID:        userID,
				Personality:   bot.Personality,
				CategoryID:    botCategory,
				Description:   bot.Description,
				AvatarURL:     bot.AvatarUrl,
				BackgroundURL: bot.BackgroundUrl,
				Location:      bot.Location,
				Published:     bot.Published,
				Active:        bot.Active,
				Like:          bot.Likes,
				CreatedAt:     bot.CreatedAt,
				UpdatedAt:     bot.UpdatedAt,
			},
			Messages: []entity.Message{},
		})
	}
	total, err = rc.repo.ListChatsByUserIDCount(context.Background(), userID.String())
	if len(repoChats) == 0 {
		return []entity.Chat{}, total, nil
	}
	return
}
