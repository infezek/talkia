package usecase_chat

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type List struct {
	RepoChat repository.RepositoryChat
	Cfg      *config.Config
}

func NewList(repoChat repository.RepositoryChat, cfg *config.Config) *List {
	return &List{
		RepoChat: repoChat,
		Cfg:      cfg,
	}
}

func (u *List) Execute(pagination entity.Pagination, userID string) (entity.Pagination, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return entity.Pagination{}, domain_error.BadRequest(fmt.Sprintf("invalid user id: %s", userID))
	}
	chats, total, err := u.RepoChat.FindByUserID(pagination, userUUID)
	if err != nil {
		return entity.Pagination{}, domain_error.BadRequest(err.Error())
	}
	pg := entity.NewPagination(pagination.PerPage, pagination.Page)
	pg.SetData(chats)
	pg.SetTotal(total)
	return *pg, nil
}
