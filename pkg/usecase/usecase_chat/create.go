package usecase_chat

import (
	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type Create struct {
	RepoChat repository.RepositoryChat
	RepoBot  repository.RepositoryBot
	RepoUser repository.RepositoryUser
	Cfg      *config.Config
}

func NewCreate(repoChat repository.RepositoryChat, repoBot repository.RepositoryBot, repoUser repository.RepositoryUser, cfg *config.Config) *Create {
	return &Create{
		RepoChat: repoChat,
		RepoBot:  repoBot,
		RepoUser: repoUser,
		Cfg:      cfg,
	}
}

func (u *Create) Execute(params CreateDtoInput) (CreateDtoOutput, error) {
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		return CreateDtoOutput{}, domain_error.BadRequest("invalid user id")
	}
	botID, err := uuid.Parse(params.BotID)
	if err != nil {
		return CreateDtoOutput{}, domain_error.BadRequest("invalid bot id")
	}
	user, err := u.RepoUser.FindByID(userID)
	if err != nil {
		return CreateDtoOutput{}, domain_error.BadRequest("user not found")
	}
	if user.VerifyIfTheConversationHasAlreadyStarted(botID) {
		return CreateDtoOutput{}, domain_error.BadRequest("conversation already exists")
	}
	bot, err := u.RepoBot.FindByID(botID)
	if err != nil {
		return CreateDtoOutput{}, domain_error.NotFound("bot not found")
	}
	chat := entity.NewChat(
		uuid.Nil,
		userID,
		botID,
		bot,
		[]entity.Message{},
		[]entity.Preference{},
	)
	if err := u.RepoChat.Create(chat); err != nil {
		return CreateDtoOutput{}, domain_error.BadRequest(err.Error())
	}
	return CreateDtoOutput{
		ID: chat.ID.String(),
	}, nil
}
