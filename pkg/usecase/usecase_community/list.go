package usecase_community

import (
	"time"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type ListComunity struct {
	repoCommunity repository.RepositoryCommunity
	repoUser      repository.RepositoryUser
	Cfg           *config.Config
}

func NewList(cfg *config.Config, repoCommunity repository.RepositoryCommunity, repoUser repository.RepositoryUser) *ListComunity {
	return &ListComunity{
		repoCommunity: repoCommunity,
		repoUser:      repoUser,
		Cfg:           cfg,
	}
}

func (l *ListComunity) Execute(userID string) (entity.Community, error) {
	trendsDateEnd := time.Now().AddDate(0, 0, 3)
	trendsDateStart := trendsDateEnd.AddDate(0, 0, -7)
	trends, err := l.repoCommunity.ListTreands(trendsDateStart, trendsDateEnd)
	if err != nil {
		return entity.Community{}, domain_error.BadRequest(err.Error())
	}
	categories, err := l.repoCommunity.ListCategoriesTreands(trendsDateStart, trendsDateEnd)
	if err != nil {
		return entity.Community{}, err
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return entity.Community{}, domain_error.BadRequest("invalid user id")
	}
	categoriesByUser, err := l.repoUser.ListCategoriesByUserID(userUUID)
	if err != nil {
		return entity.Community{}, err
	}

	var categoryOfUser []entity.Category
	var categoryNotUser []entity.Category
	for _, category := range categories {
		userNotCategory := true
		if len(category.Bots) == 0 {
			continue
		}
		for _, categoryByUser := range categoriesByUser {
			if category.ID.String() == categoryByUser.ID.String() {
				categoryOfUser = append(categoryOfUser, category)
				userNotCategory = false
				break
			}
		}
		if userNotCategory {
			categoryNotUser = append(categoryNotUser, category)
		}
	}
	result := append(categoryOfUser, categoryNotUser...)
	return entity.Community{
		Trends:     trends,
		Categories: result,
	}, nil

}
