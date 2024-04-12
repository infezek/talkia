package usecase_bot

import (
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type ListByName struct {
	RepoBot repository.RepositoryBot
	Cfg     *config.Config
}

func NewListByName(cfg *config.Config, repoBot repository.RepositoryBot) *ListByName {
	return &ListByName{
		RepoBot: repoBot,
		Cfg:     cfg,
	}
}

func (u *ListByName) Execute(pagintion entity.Pagination, name string) (*entity.Pagination, error) {
	bots, err := u.RepoBot.FindByName(pagintion, name)
	if err != nil {
		return nil, domain_error.BadRequest(err.Error())
	}
	pagination := entity.NewPagination(pagintion.PerPage, pagintion.Page)
	pagination.SetData(bots)
	return pagination, nil
}
