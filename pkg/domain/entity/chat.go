package entity

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID              uuid.UUID    `json:"id"`
	UserID          uuid.UUID    `json:"-"`
	BotID           uuid.UUID    `json:"-"`
	Bot             Bot          `json:"bot"`
	Messages        []Message    `json:"messages,omitempty"`
	UserPreferences []Preference `json:"preferences,omitempty"`
}

func NewChat(id uuid.UUID, userID uuid.UUID, botID uuid.UUID, bot Bot, messages []Message, preferences []Preference) Chat {
	if id == uuid.Nil {
		id, _ = uuid.NewV7()
	}

	return Chat{
		ID:              id,
		UserID:          userID,
		BotID:           botID,
		Bot:             bot,
		Messages:        messages,
		UserPreferences: preferences,
	}
}

var (
	MessageSystem Who = "system"
	MessageUser   Who = "user"
)

type Who string

func (w Who) String() string {
	return string(w)
}

type Message struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	UserID    uuid.UUID
	Who       Who
	Message   string
	CreatedAt time.Time
}

func NewMessage(id uuid.UUID, chatID uuid.UUID, userID uuid.UUID, who Who, message string, createdAt time.Time) Message {
	if id == uuid.Nil {
		id, _ = uuid.NewV7()
	}
	return Message{
		ID:        id,
		ChatID:    chatID,
		UserID:    userID,
		Who:       who,
		Message:   message,
		CreatedAt: createdAt,
	}
}

func (m *Message) AddChat(chat Chat) {
	m.ChatID = chat.ID
	m.UserID = chat.BotID
}

type Preference struct {
	ID              uuid.UUID
	ChatID          uuid.UUID
	UserID          uuid.UUID
	PreferenceKey   PreferenceKey
	PreferenceValue string
}

type PreferenceKey string

var (
	PreferenceUser PreferenceKey = "User"
)

func (p Preference) User() string {
	if p.PreferenceKey == PreferenceUser {
		return string(p.PreferenceKey)
	}
	return ""
}
