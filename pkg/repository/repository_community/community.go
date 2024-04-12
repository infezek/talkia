package repository_community

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	repositoryDB "github.com/infezek/app-chat/db/repository"
	"github.com/infezek/app-chat/pkg/domain/entity"
)

type RepositoryCommunity struct {
	db   *sql.DB
	repo *repositoryDB.Queries
}

func New(db *sql.DB) *RepositoryCommunity {
	repo := repositoryDB.New(db)
	return &RepositoryCommunity{db, repo}
}

func (r *RepositoryCommunity) ListTreands(startDate, endDate time.Time) ([]entity.Bot, error) {
	resp, err := r.repo.ListTrends(context.Background(), repositoryDB.ListTrendsParams{
		CreatedAt:   startDate,
		CreatedAt_2: endDate,
	})
	if err != nil && err != sql.ErrNoRows {
		return []entity.Bot{}, err
	}
	bots := make([]entity.Bot, 0)
	for _, bot := range resp {
		botID, err := uuid.Parse(bot.ID)
		if err != nil {
			return []entity.Bot{}, err
		}
		userID, err := uuid.Parse(bot.UserID)
		if err != nil {
			return []entity.Bot{}, err
		}
		categoryID, err := uuid.Parse(bot.CategoryID)
		if err != nil {
			return []entity.Bot{}, err
		}
		bots = append(bots, entity.Bot{
			ID:            botID,
			UserID:        userID,
			Name:          bot.Name,
			CategoryID:    categoryID,
			Description:   bot.Description,
			AvatarURL:     bot.AvatarUrl,
			Personality:   bot.Personality,
			BackgroundURL: bot.BackgroundUrl,
			Published:     bot.Published,
			Active:        bot.Active,
			Like:          bot.Likes,
			CreatedAt:     bot.CreatedAt,
			UpdatedAt:     bot.UpdatedAt,
			Location:      bot.Location,
		})
	}
	return bots, nil
}

func (r *RepositoryCommunity) ListCategoriesTreands(startDate, endDate time.Time) ([]entity.Category, error) {
	categoriesRepo, err := r.repo.ListCategoriesTrends(context.Background(), repositoryDB.ListCategoriesTrendsParams{
		CreatedAt:   startDate,
		CreatedAt_2: endDate,
	})
	if err != nil && err != sql.ErrNoRows {
		return []entity.Category{}, err
	}
	categories := []entity.Category{}
	for _, category := range categoriesRepo {
		botsRepo, err := r.repo.ListBotsOfCategoryTrends(context.Background(), repositoryDB.ListBotsOfCategoryTrendsParams{
			CreatedAt:   startDate,
			CreatedAt_2: endDate,
			CategoryID:  category.ID,
		})
		if err != nil && err != sql.ErrNoRows {
			return []entity.Category{}, err
		}
		categoryID, err := uuid.Parse(category.ID)
		if err != nil {
			return []entity.Category{}, err
		}
		bots := []*entity.Bot{}
		for _, bot := range botsRepo {
			botID, err := uuid.Parse(bot.ID)
			if err != nil {
				return []entity.Category{}, err
			}
			userID, err := uuid.Parse(bot.UserID)
			if err != nil {
				return []entity.Category{}, err
			}
			categoryID, err := uuid.Parse(bot.CategoryID)
			if err != nil {
				return []entity.Category{}, err
			}
			bots = append(bots, &entity.Bot{
				ID:            botID,
				UserID:        userID,
				Name:          bot.Name,
				CategoryID:    categoryID,
				Description:   bot.Description,
				AvatarURL:     bot.AvatarUrl,
				Personality:   bot.Personality,
				BackgroundURL: bot.BackgroundUrl,
				Published:     bot.Published,
				Active:        bot.Active,
				Like:          bot.Likes,
				CreatedAt:     bot.CreatedAt,
				UpdatedAt:     bot.UpdatedAt,
				Location:      bot.Location,
			})
		}
		categories = append(categories, entity.Category{
			ID:        categoryID,
			Name:      category.Name,
			Bots:      bots,
			Active:    category.Active,
			CreatedAt: category.CreatedAt,
		})
	}
	return categories, nil
}

func (r *RepositoryCommunity) ListTreandsByLocation(location string) ([]entity.Bot, error) {
	return []entity.Bot{}, fmt.Errorf("not implemented")
}
