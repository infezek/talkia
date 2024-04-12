package usecase_chat

import (
	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type Chat struct {
	RepoChat repository.RepositoryChat
	Cfg      *config.Config
}

func NewChatGet(repoChat repository.RepositoryChat, cfg *config.Config) *Chat {
	return &Chat{
		RepoChat: repoChat,
		Cfg:      cfg,
	}
}

func (u *Chat) Execute(chatID string, userID string) (entity.Chat, error) {
	chatUUID, err := uuid.Parse(chatID)
	if err != nil {
		return entity.Chat{}, domain_error.BadRequest("invalid chat id")
	}
	chat, err := u.RepoChat.FindByID(chatUUID)
	if err != nil {
		return entity.Chat{}, domain_error.NotFound("chat already exists")
	}
	if chat.UserID.String() != userID {
		return entity.Chat{}, domain_error.BadRequest("chat not of user")
	}
	return chat, nil
}
