package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID  `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Password   string     `json:"-"`
	Platform   Platform   `json:"platform"`
	Gender     *Gender    `json:"gender"`
	AvatarURL  string     `json:"avatar_url"`
	Location   string     `json:"location"`
	Language   Language   `json:"language"`
	Categories []Category `json:"categories"`
	Chats      []Chat     `json:"chats"`
	Bots       []Bot      `json:"bots"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (u *User) CodeLanguage() string {
	if u.Language == "portuguese" {
		return "responder en español"
	}
	if u.Language == "english" {
		return "answer in English"
	}
	if u.Language == "spanish" {
		return "responder em português"
	}
	return "responda no idioma da pergunta"
}

type Language string

func (l Language) String() string {
	return string(l)
}

var (
	LanguageEnglish    Language = "english"
	LanguagePortuguese Language = "portuguese"
	LanguageSpanish    Language = "spanish"
)

type Platform string
type Gender string

var (
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
)

var (
	GenderMasculine Gender = "masculine"
	GenderFeminine  Gender = "feminine"
	GenderNeutral   Gender = "neutral"
)

func NewUser(id uuid.UUID, username string, email string, password string, language Language, platform Platform, gender *Gender, chats []Chat, bots []Bot, categories []Category, createdAt time.Time) *User {
	if id == uuid.Nil {
		id, _ = uuid.NewV7()
	}
	return &User{
		ID:         id,
		Username:   username,
		Email:      email,
		Password:   password,
		Language:   language,
		Platform:   platform,
		Categories: categories,
		Gender:     gender,
		Chats:      chats,
		Bots:       bots,
		CreatedAt:  createdAt,
	}
}

func (u *User) Update(username, language string) {
	u.Username = username
	u.Language = Language(language)
}

func (u *User) FindChatByID(chatID uuid.UUID) (Chat, error) {
	for _, c := range u.Chats {
		if c.ID == chatID {
			return c, nil
		}
	}
	return Chat{}, errors.New("chat not found")
}

func (u *User) CreateChat(chat Chat) {
	u.Chats = append(u.Chats, chat)
}

func (u *User) AddBot(bot Bot) error {
	for _, b := range u.Bots {
		if b.ID == bot.ID {
			return errors.New("bot already exists")
		}
	}
	u.Bots = append(u.Bots, bot)
	return nil
}

func (u *User) RemoveBot(bot Bot) error {
	for i, b := range u.Bots {
		if b.ID == bot.ID {
			u.Bots = append(u.Bots[:i], u.Bots[i+1:]...)
			return nil
		}
	}
	return errors.New("bot not found")
}

func (u *User) CreateBot(bot Bot) {
	u.Bots = append(u.Bots, bot)
}

func (u *User) AddCategory(categoryID uuid.UUID) {
	for _, c := range u.Categories {
		if c.ID == categoryID {
			return
		}
	}
	u.Categories = append(u.Categories, Category{ID: categoryID})
}

func (u *User) RemoveCategories() {
	u.Categories = []Category{}
}

func (u *User) AddCategories(categories []uuid.UUID) {
	for _, c := range categories {
		u.Categories = append(u.Categories, Category{ID: c})
	}
}

func (u *User) VerifyIfTheConversationHasAlreadyStarted(botID uuid.UUID) bool {
	for _, c := range u.Chats {
		if c.BotID == botID {
			return true
		}
	}
	return false
}

func (u *User) UpdateAvatar(avatarURL string) {
	u.AvatarURL = avatarURL
}
