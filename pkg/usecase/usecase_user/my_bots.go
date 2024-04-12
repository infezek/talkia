package usecase_user

import (
	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type MyBots struct {
	RepoUser repository.RepositoryUser
	Cfg      *config.Config
}

func NewMyBots(repo repository.RepositoryUser, cfg *config.Config) *MyBots {
	return &MyBots{
		RepoUser: repo,
		Cfg:      cfg,
	}
}

func (u *MyBots) Execute(pagintion entity.Pagination, userID string) (entity.Pagination, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return entity.Pagination{}, err
	}
	categories, totalData, err := u.RepoUser.ListMyBots(pagintion, userUUID)
	if err != nil {
		return entity.Pagination{}, err
	}
	pagination := entity.NewPagination(pagintion.PerPage, pagintion.Page)
	pagination.SetData(categories)
	pagination.SetTotal(totalData)
	return *pagination, nil
}
