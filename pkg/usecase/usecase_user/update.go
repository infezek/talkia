package usecase_user

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type Update struct {
	RepoUser repository.RepositoryUser
	Cfg      *config.Config
}

func UpdateUseCase(cfg *config.Config, repoUser repository.RepositoryUser) *Update {
	return &Update{
		RepoUser: repoUser,
		Cfg:      cfg,
	}
}

func (u *Update) Execute(params UpdateUserDtoInput) error {
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		return domain_error.BadRequest(fmt.Sprintf("not parse id %s", err.Error()))
	}
	user, err := u.RepoUser.FindByID(userID)
	user.Update(params.Username, params.Language)
	err = u.RepoUser.Update(user)
	if err != nil {
		return domain_error.BadRequest(err.Error())
	}
	return nil
}
