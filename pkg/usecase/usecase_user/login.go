package usecase_user

import (
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/adapter"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type Login struct {
	AdapterToken adapter.AdapterToken
	RepoUser     repository.RepositoryUser
	Cfg          *config.Config
}

func NewLoginUseCase(cfg *config.Config, repoUser repository.RepositoryUser, token adapter.AdapterToken) *Login {
	return &Login{
		AdapterToken: token,
		RepoUser:     repoUser,
		Cfg:          cfg,
	}
}

func (l *Login) Execute(email string) (string, error) {
	user, err := l.RepoUser.FindUserByEmail(email)
	if err != nil {
		return "", domain_error.BadRequest(err.Error())
	}
	if user == nil {
		return "", domain_error.NotFound("user not found")
	}
	token, err := l.AdapterToken.CreateToken(user.ID.String(), user.Email)
	if err != nil {
		return "", domain_error.BadRequest(err.Error())
	}
	return token, nil
}
