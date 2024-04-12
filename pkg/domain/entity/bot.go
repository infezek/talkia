package entity

import (
	"time"

	"github.com/google/uuid"
)

type Bot struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Name          string    `json:"name"`
	CategoryID    uuid.UUID `json:"category_id"`
	Personality   string    `json:"personality"`
	Description   string    `json:"description"`
	AvatarURL     string    `json:"avatar_url"`
	BackgroundURL string    `json:"background_url"`
	Location      string    `json:"location"`
	Published     bool      `json:"published"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Like          int64     `json:"like"`
}

type LikeBot struct {
	ID     uuid.UUID
	UserID uuid.UUID
	BotID  uuid.UUID
}

func NewBot(id, userID, categoryID uuid.UUID, name, personality, description, avatarURL, backgroundURL, location string, createdAt, updatedAt time.Time) *Bot {
	var err error
	if id == uuid.Nil {
		id, err = uuid.NewV7()
		if err != nil {
			return nil
		}
	}
	return &Bot{
		ID:            id,
		UserID:        userID,
		CategoryID:    categoryID,
		Name:          name,
		Personality:   personality,
		Description:   description,
		AvatarURL:     avatarURL,
		BackgroundURL: backgroundURL,
		Location:      location,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (b *Bot) Update(name, personality, description, avatarURL, BackgroundURL string, location string) {
	b.Name = name
	b.Personality = personality
	b.Description = description
	b.AvatarURL = avatarURL
	b.Location = location
	b.UpdatedAt = time.Now()
}

func (b *Bot) UpdateAvatarURL(avatarURL string) {
	b.AvatarURL = avatarURL
	b.UpdatedAt = time.Now()
}

func (b *Bot) UpdateBackgroundURL(backgroundURL string) {
	b.BackgroundURL = backgroundURL
	b.UpdatedAt = time.Now()
}
