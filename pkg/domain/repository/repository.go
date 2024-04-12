package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/domain/entity"
)

type RepositoryUser interface {
	Create(user entity.User) error
	Update(user entity.User) error
	Desactive(id uuid.UUID) error
	CreateChat(chat entity.Chat) error
	ListMyBots(pagination entity.Pagination, userID uuid.UUID) ([]entity.Bot, int64, error)
	FindByID(id uuid.UUID) (entity.User, error)
	FindUserByEmail(email string) (*entity.User, error)
	DeleteCategoriesUserID(userID uuid.UUID) error
	AddCategoriesByUserID(user entity.User) error
	PublishBot(botID uuid.UUID, userID uuid.UUID) error
	LikeToBot(userID, botID uuid.UUID) error
	FindLikeBotByUserAndBot(userID uuid.UUID, botID uuid.UUID) (bool, error)
	FindLikeBotByUserID(userID uuid.UUID) ([]entity.Bot, error)
	FindLikeBotByBotID(botID uuid.UUID) ([]entity.User, error)
	UpdateAvatarURL(user entity.User) error
	ListCategoriesByUserID(userID uuid.UUID) ([]entity.Category, error)
}

type RepositoryBot interface {
	Create(bot entity.Bot) error
	Update(bot entity.Bot) error
	Desactive(id uuid.UUID) error
	FindByName(pagination entity.Pagination, name string) ([]entity.Bot, error)
	FindByID(id uuid.UUID) (entity.Bot, error)
}

type RepositoryChat interface {
	Create(chat entity.Chat) error
	SaveMessage(chat entity.Chat, message entity.Message) error
	SavePreferences(chat entity.Chat, message entity.Message) error
	Update(chat entity.Chat) error
	Desactive(chatID uuid.UUID) error
	FindByID(chatID uuid.UUID) (entity.Chat, error)
	FindByUserID(pagination entity.Pagination, userID uuid.UUID) (chats []entity.Chat, total int64, err error)
}

type RepositoryCategory interface {
	Create(category entity.Category) error
	FindByID(id uuid.UUID) (*entity.Category, error)
	FindByName(name string) (*entity.Category, error)
	Update(category entity.Category) error
	Desactive(id uuid.UUID) error
	List(pagination entity.Pagination) ([]entity.Category, int64, error)
}

type RepositoryCommunity interface {
	ListTreands(startDate, endDate time.Time) ([]entity.Bot, error)
	ListCategoriesTreands(startDate, endDate time.Time) ([]entity.Category, error)
	ListTreandsByLocation(location string) ([]entity.Bot, error)
}
