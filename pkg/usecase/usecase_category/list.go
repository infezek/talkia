package usecase_category

import (
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type List struct {
	RepoCategory repository.RepositoryCategory
	Cfg          *config.Config
}

func NewList(RepoCategory repository.RepositoryCategory, cfg *config.Config) *List {
	return &List{
		RepoCategory: RepoCategory,
		Cfg:          cfg,
	}
}

func (u *List) Execute(pagintion entity.Pagination) (entity.Pagination, error) {
	categories, totalData, err := u.RepoCategory.List(pagintion)
	if err != nil {
		return entity.Pagination{}, err
	}
	pagination := entity.NewPagination(pagintion.PerPage, pagintion.Page)
	pagination.SetData(categories)
	pagination.SetTotal(totalData)
	return *pagination, nil
}
