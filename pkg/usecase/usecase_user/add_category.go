package usecase_user

import (
	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type AddCategoryUseCase struct {
	RepoUser repository.RepositoryUser
	Cfg      *config.Config
}

func NewAddCategory(cfg *config.Config, repoUser repository.RepositoryUser) *AddCategoryUseCase {
	return &AddCategoryUseCase{
		RepoUser: repoUser,
		Cfg:      cfg,
	}
}

type AddCategoryDtoInput struct {
	UserID     string
	Categories []string
}

func (u *AddCategoryUseCase) Execute(params AddCategoryDtoInput) error {
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		return domain_error.BadRequest(err.Error())
	}
	user, err := u.RepoUser.FindByID(userID)
	if err != nil {
		return domain_error.BadRequest(err.Error())
	}
	var categoriesFormatted []uuid.UUID
	for _, category := range params.Categories {
		id, err := uuid.Parse(category)
		if err != nil {
			return domain_error.BadRequest("category id invalid")
		}
		categoriesFormatted = append(categoriesFormatted, id)
	}
	if err := u.RepoUser.DeleteCategoriesUserID(userID); err != nil {
		return domain_error.BadRequest(err.Error())
	}
	user.RemoveCategories()
	user.AddCategories(categoriesFormatted)
	if err := u.RepoUser.AddCategoriesByUserID(user); err != nil {
		return domain_error.BadRequest(err.Error())
	}
	return nil
}
