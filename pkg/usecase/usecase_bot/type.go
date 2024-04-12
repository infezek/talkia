package usecase_bot

import (
	"time"
)

type UploadImageDtoInput struct {
	BotID  string
	UserID string
	Files  []File
}

type File struct {
	File []byte
	Name string
	Type TypeImage
}

type TypeImage string

var (
	TypeImageAvatar     TypeImage = "avatar"
	TypeImageBackground TypeImage = "background"
)

type CreateDtoInput struct {
	UserID      string
	Name        string
	CategoryID  string
	Personality string
	Description string
	Location    string
}

type UpdateDtoInput struct {
	ID            string
	UserID        string
	Name          string
	Personality   string
	Description   string
	AvatarURL     string
	BackgroundURL string
	Language      string
	Location      string
}

type CreateDtoOutput struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CategoryID  string    `json:"category_id"`
	Personality string    `json:"personality"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateDtoOutput CreateDtoOutput

type LikeInput struct {
	UserID string `json:"user_id"`
	BotID  string `json:"bot_id"`
}
