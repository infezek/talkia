package repository_category

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	repositoryDB "github.com/infezek/app-chat/db/repository"
	"github.com/infezek/app-chat/pkg/domain/entity"
)

type RepositorCategories struct {
	db   *sql.DB
	repo *repositoryDB.Queries
}

func New(db *sql.DB) *RepositorCategories {
	repo := repositoryDB.New(db)
	return &RepositorCategories{db, repo}
}

func (ru *RepositorCategories) Create(category entity.Category) error {
	return ru.repo.CreateCategory(context.Background(), repositoryDB.CreateCategoryParams{
		Name:      category.Name,
		ID:        category.ID.String(),
		Active:    category.Active,
		CreatedAt: category.CreatedAt,
	})
}
func (ru *RepositorCategories) FindByID(id uuid.UUID) (*entity.Category, error) {
	categoryRepo, err := ru.repo.FindCategoryByID(context.Background(), id.String())
	if err != nil {
		return &entity.Category{}, err
	}
	categoryID, err := uuid.Parse(categoryRepo.ID)
	if err != nil {
		return &entity.Category{}, err
	}
	return &entity.Category{
		ID:        categoryID,
		Name:      categoryRepo.Name,
		Active:    categoryRepo.Active,
		CreatedAt: categoryRepo.CreatedAt,
	}, nil
}
func (ru *RepositorCategories) FindByName(name string) (*entity.Category, error) {
	categoryRepo, err := ru.repo.FindCategoryByName(context.Background(), name)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return &entity.Category{}, err
	}
	categoryID, err := uuid.Parse(categoryRepo.ID)
	if err != nil {
		return &entity.Category{}, err
	}
	return &entity.Category{
		ID:        categoryID,
		Name:      categoryRepo.Name,
		Active:    categoryRepo.Active,
		CreatedAt: categoryRepo.CreatedAt,
	}, nil
}
func (ru *RepositorCategories) Update(category entity.Category) error {
	return ru.repo.UpdateCategory(context.Background(), repositoryDB.UpdateCategoryParams{
		ID:     category.ID.String(),
		Name:   category.Name,
		Active: category.Active,
	})
}

func (ru *RepositorCategories) Desactive(id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func (rc *RepositorCategories) List(pagintion entity.Pagination) (categories []entity.Category, total int64, err error) {
	repoCategories, err := rc.repo.ListCategories(context.Background(), repositoryDB.ListCategoriesParams{
		Limit:  pagintion.PerPage,
		Offset: pagintion.Offset,
	})
	if err != nil && err != sql.ErrNoRows {
		return categories, total, err
	}
	if err != nil && err == sql.ErrNoRows {
		return []entity.Category{}, total, err
	}
	for _, category := range repoCategories {
		id, err := uuid.Parse(category.ID)
		if err != nil {
			return categories, total, nil

		}
		categories = append(categories, entity.Category{
			ID:        id,
			Name:      category.Name,
			Active:    category.Active,
			CreatedAt: category.CreatedAt,
		})
	}
	total, err = rc.repo.ListCategoriesCount(context.Background())
	if len(repoCategories) == 0 {
		return []entity.Category{}, total, nil
	}
	return
}
