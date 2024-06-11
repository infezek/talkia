package usecase_user

import (
	"time"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/adapter"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/entity"
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

func (l *Login) Execute(email string) (LoginDto, error) {
	user, err := l.RepoUser.FindUserByEmail(email)
	if err != nil {
		return LoginDto{}, domain_error.BadRequest(err.Error())
	}
	if user == nil {
		return LoginDto{}, domain_error.NotFound("user not found")
	}
	token, err := l.AdapterToken.CreateToken(user.ID.String(), user.Email)
	if err != nil {
		return LoginDto{}, domain_error.BadRequest(err.Error())
	}
	return LoginDto{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Platform:  user.Platform,
		Gender:    user.Gender,
		AvatarURL: user.AvatarURL,
		Location:  user.Location,
		Language:  user.Language,
		CreatedAt: user.CreatedAt,
		Token:     token,
	}, nil
}

type LoginDto struct {
	ID        uuid.UUID       `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Platform  entity.Platform `json:"platform"`
	Gender    *entity.Gender  `json:"gender"`
	AvatarURL string          `json:"avatar_url"`
	Location  string          `json:"location"`
	Language  entity.Language `json:"language"`
	CreatedAt time.Time       `json:"created_at"`
	Token     string          `json:"token"`
}
