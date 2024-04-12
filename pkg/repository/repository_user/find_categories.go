package repository_user

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	repositoryDB "github.com/infezek/app-chat/db/repository"
	"github.com/infezek/app-chat/pkg/domain/entity"
)

func findCategoriesByUserID(repo *repositoryDB.Queries, userID uuid.UUID) ([]entity.Category, error) {
	categoriesRepo, err := repo.FindCategoriesByUserID(context.Background(), userID.String())
	if sql.ErrNoRows == err {
		return []entity.Category{}, nil
	}
	if err != nil && sql.ErrNoRows != err {
		return []entity.Category{}, err
	}
	categories := []entity.Category{}
	for _, category := range categoriesRepo {
		categoryID, err := uuid.Parse(category.ID)
		if err != nil {
			return []entity.Category{}, err
		}
		categories = append(categories, entity.Category{
			ID:        categoryID,
			Name:      string(category.Name),
			Active:    category.Active,
			CreatedAt: category.CreatedAt,
		})
	}
	return categories, nil
}
