package usecase_user

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type UseCase struct {
	RepoUser repository.RepositoryUser
	Cfg      *config.Config
}

func CreateUserUseCase(repoUser repository.RepositoryUser, cfg *config.Config) *UseCase {
	return &UseCase{
		RepoUser: repoUser,
		Cfg:      cfg,
	}
}

func (u *UseCase) Execute(params CreateUserDtoInput) (CreateDto, error) {
	user, err := u.RepoUser.FindUserByEmail(params.Email)
	if user == nil {
		user = entity.NewUser(
			uuid.Nil,
			params.Username,
			params.Email,
			params.Password,
			entity.Language(params.Language),
			entity.Platform(params.Platform),
			nil,
			[]entity.Chat{},
			[]entity.Bot{},
			[]entity.Category{},
			time.Now())
		if err := u.RepoUser.Create(*user); err != nil {
			return CreateDto{}, domain_error.BadRequest(err.Error())
		}
		return CreateDto{
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
		}, nil

	}
	if err != nil && sql.ErrNoRows != err {
		return CreateDto{}, domain_error.BadRequest(err.Error())
	}
	return CreateDto{
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
	}, nil
}

type CreateDto struct {
	Username  string
	AvatarURL string
}
