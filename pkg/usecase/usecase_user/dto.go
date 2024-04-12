package usecase_user

import (
	"time"

	"github.com/infezek/app-chat/pkg/domain/entity"
)

type UpdateUserDtoInput struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	Language string `json:"language"`
}

type CreateUserDtoInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Platform string `json:"platform"`
	Gender   string `json:"gender"`
	Location string `json:"location"`
	Language string `json:"language"`
}

type CreateUserDtoOutput struct {
	ID         string            `json:"id"`
	Username   string            `json:"username"`
	Email      string            `json:"email"`
	Platform   string            `json:"platform"`
	Location   string            `json:"location"`
	Chats      []entity.Chat     `json:"chats"`
	Bots       []entity.Bot      `json:"bots"`
	Categories []entity.Category `json:"categories"`
	CreatedAt  time.Time         `json:"created_at"`
}
