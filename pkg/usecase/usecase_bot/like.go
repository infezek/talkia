package usecase_bot

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type Like struct {
	RepoUser repository.RepositoryUser
	Cfg      *config.Config
}

func NewLike(cfg *config.Config, repoUser repository.RepositoryUser) *Like {
	return &Like{
		RepoUser: repoUser,
		Cfg:      cfg,
	}
}

func (u *Like) Execute(params LikeInput) error {
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		return fmt.Errorf("invalid user id")
	}
	botID, err := uuid.Parse(params.BotID)
	if err != nil {
		return fmt.Errorf("invalid user id")
	}
	ok, err := u.RepoUser.FindLikeBotByUserAndBot(userID, botID)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("user already liked bot")
	}
	if err = u.RepoUser.LikeToBot(userID, botID); err != nil {
		return err
	}
	return nil
}
