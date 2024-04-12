package usecase_category

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type Create struct {
	RepoCategory repository.RepositoryCategory
	Cfg          *config.Config
}

func NewCreate(RepoCategory repository.RepositoryCategory, cfg *config.Config) *Create {
	return &Create{
		RepoCategory: RepoCategory,
		Cfg:          cfg,
	}
}

func (u *Create) Execute(params CreateDtoInput) (CreateDtoOutput, error) {
	category, err := u.RepoCategory.FindByName(params.Name)
	if category != nil {
		return CreateDtoOutput{}, fmt.Errorf("category already exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return CreateDtoOutput{}, err
	}
	category = entity.NewCategory(
		uuid.Nil,
		params.Name,
		params.Active,
		nil)

	if err := u.RepoCategory.Create(*category); err != nil {
		return CreateDtoOutput{}, err
	}
	return CreateDtoOutput{
		ID:        category.ID.String(),
		Name:      category.Name,
		Active:    category.Active,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdateAt,
	}, nil
}
