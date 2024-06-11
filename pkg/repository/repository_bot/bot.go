package repository_bot

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	repositoryDB "github.com/infezek/app-chat/db/repository"
	"github.com/infezek/app-chat/pkg/domain/entity"
)

type RepositoryBots struct {
	db   *sql.DB
	repo *repositoryDB.Queries
}

func New(db *sql.DB) *RepositoryBots {
	repo := repositoryDB.New(db)
	return &RepositoryBots{db, repo}
}

func (r *RepositoryBots) Create(bot entity.Bot) error {
	return r.repo.CreateBot(context.Background(), repositoryDB.CreateBotParams{
		ID:            bot.ID.String(),
		UserID:        bot.UserID.String(),
		CategoryID:    bot.CategoryID.String(),
		Name:          bot.Name,
		Personality:   bot.Personality,
		Description:   bot.Description,
		BackgroundUrl: bot.BackgroundURL,
		AvatarUrl:     bot.AvatarURL,
		Location:      bot.Location,
	})
}
func (r *RepositoryBots) Update(bot entity.Bot) error {
	return r.repo.UpdateBot(context.Background(), repositoryDB.UpdateBotParams{
		ID:            bot.ID.String(),
		Name:          bot.Name,
		Personality:   bot.Personality,
		Description:   bot.Description,
		BackgroundUrl: bot.BackgroundURL,
		AvatarUrl:     bot.AvatarURL,
		Location:      bot.Location,
		Active:        bot.Active,
		CategoryID:    bot.CategoryID.String(),
	})
}
func (r *RepositoryBots) Desactive(id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}
func (r *RepositoryBots) FindByName(pagination entity.Pagination, name string) ([]entity.Bot, error) {
	botRepo, err := r.repo.FindBotsByName(context.Background(), repositoryDB.FindBotsByNameParams{
		Name:   fmt.Sprintf("%%%s%%", name),
		Limit:  pagination.PerPage,
		Offset: pagination.Offset,
	})
	if err != nil {
		return []entity.Bot{}, err
	}
	bots := []entity.Bot{}
	for _, bot := range botRepo {
		botUUID, err := uuid.Parse(bot.ID)
		if err != nil {
			return []entity.Bot{}, err
		}
		userUUID, err := uuid.Parse(bot.UserID)
		if err != nil {
			return []entity.Bot{}, err
		}
		bots = append(bots, entity.Bot{
			ID:            botUUID,
			UserID:        userUUID,
			Name:          bot.Name,
			Personality:   bot.Personality,
			Description:   bot.Description,
			AvatarURL:     bot.AvatarUrl,
			BackgroundURL: bot.BackgroundUrl,
			CreatedAt:     bot.CreatedAt,
			UpdatedAt:     bot.UpdatedAt,
			CategoryID:    uuid.MustParse(bot.CategoryID),
			Location:      bot.Location,
			Published:     bot.Published,
			Active:        bot.Active,
			Like:          bot.Likes,
		})
	}

	return bots, nil

}
func (r *RepositoryBots) FindByID(id uuid.UUID) (entity.Bot, error) {

	bot, err := r.repo.GetBotByID(context.Background(), id.String())
	if err != nil {
		return entity.Bot{}, err
	}

	return entity.Bot{
		ID:            uuid.MustParse(bot.ID),
		UserID:        uuid.MustParse(bot.UserID),
		Name:          bot.Name,
		CategoryID:    uuid.MustParse(bot.CategoryID),
		Location:      bot.Location,
		Published:     bot.Published,
		Active:        bot.Active,
		Personality:   bot.Personality,
		Description:   bot.Description,
		AvatarURL:     bot.AvatarUrl,
		BackgroundURL: bot.BackgroundUrl,
		Like:          bot.Likes,
		CreatedAt:     bot.CreatedAt,
		UpdatedAt:     bot.UpdatedAt,
	}, nil
}
