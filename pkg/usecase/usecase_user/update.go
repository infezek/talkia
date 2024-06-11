package usecase_user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/entity"
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

func (u *Update) Execute(params UpdateUserDtoInput) (UpdateDto, error) {
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		return UpdateDto{}, domain_error.BadRequest(fmt.Sprintf("not parse id %s", err.Error()))
	}
	user, err := u.RepoUser.FindByID(userID)
	if err != nil {
		return UpdateDto{}, domain_error.BadRequest(err.Error())
	}
	user.Update(params.Username, params.Language)
	err = u.RepoUser.Update(user)
	if err != nil {
		return UpdateDto{}, domain_error.BadRequest(err.Error())
	}
	return UpdateDto{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Platform:  user.Platform,
		Gender:    user.Gender,
		Location:  user.Location,
		Language:  user.Language,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}, nil
}

type UpdateDto struct {
	ID        uuid.UUID       `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Platform  entity.Platform `json:"platform"`
	Gender    *entity.Gender  `json:"gender"`
	AvatarURL string          `json:"avatar_url"`
	Location  string          `json:"location"`
	Language  entity.Language `json:"language"`
	CreatedAt time.Time       `json:"created_at"`
}
