package usecase_category

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type CategoryUpdateUseCase struct {
	RepoCategory repository.RepositoryCategory
	Cfg          *config.Config
}

func NewUpdate(RepoCategory repository.RepositoryCategory, cfg *config.Config) *CategoryUpdateUseCase {
	return &CategoryUpdateUseCase{
		RepoCategory: RepoCategory,
		Cfg:          cfg,
	}
}

func (u *CategoryUpdateUseCase) Execute(params UpdateDtoInput) (UpdateDtoOutput, error) {
	id, err := uuid.Parse(params.ID)
	if err != nil {
		return UpdateDtoOutput{}, err
	}
	category, err := u.RepoCategory.FindByID(id)
	if category == nil {
		return UpdateDtoOutput{}, fmt.Errorf("category not found")
	}
	uuid, err := uuid.Parse(params.ID)
	if err != nil {
		return UpdateDtoOutput{}, err
	}
	category = entity.NewCategory(
		uuid,
		params.Name,
		params.Active,
		nil)
	if err := u.RepoCategory.Update(*category); err != nil {
		return UpdateDtoOutput{}, err
	}
	return UpdateDtoOutput{
		ID:        category.ID.String(),
		Name:      category.Name,
		Active:    category.Active,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdateAt,
	}, nil
}
