package repository_user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	repositoryDB "github.com/infezek/app-chat/db/repository"
	"github.com/infezek/app-chat/pkg/domain/entity"
)

type RepositoryUsers struct {
	db   *sql.DB
	repo *repositoryDB.Queries
}

func New(db *sql.DB) *RepositoryUsers {
	repo := repositoryDB.New(db)
	return &RepositoryUsers{
		db:   db,
		repo: repo,
	}
}

func (ru *RepositoryUsers) Create(user entity.User) error {
	return ru.repo.CreateUser(context.Background(), repositoryDB.CreateUserParams{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Platform: repositoryDB.UsersPlatform(user.Platform),
		Location: user.Location,
		Language: repositoryDB.UsersLanguage(user.Language),
	})
}
func (ru *RepositoryUsers) List(pagination entity.Pagination) ([]entity.User, int64, error) {
	repos, err := ru.repo.List(context.Background(), repositoryDB.ListParams{
		Limit:  pagination.PerPage,
		Offset: pagination.Offset,
	})
	if err != nil {
		return []entity.User{}, 0, err
	}
	users := []entity.User{}
	for _, repo := range repos {
		id, err := uuid.Parse(repo.ID)
		if err != nil {
			return []entity.User{}, 0, err
		}
		categories, err := findCategoriesByUserID(ru.repo, id)
		if err != nil {
			return []entity.User{}, 0, err
		}
		chats, err := findChatsByUserID(ru.repo, id)
		if err != nil {
			return []entity.User{}, 0, err
		}
		bots, err := findBotsByUserID(ru.repo, id)
		if err != nil {
			return []entity.User{}, 0, err
		}

		users = append(users, entity.User{
			ID:         id,
			Username:   repo.Username,
			Email:      repo.Email,
			Password:   repo.Password,
			Platform:   entity.Platform(repo.Platform),
			Language:   entity.Language(repo.Language),
			Location:   repo.Location,
			Gender:     &entity.GenderFeminine,
			AvatarURL:  repo.AvatarUrl.String,
			Categories: categories,
			Chats:      chats,
			Bots:       bots,
			CreatedAt:  repo.CreatedAt,
		})
	}
	total, err := ru.repo.ListCount(context.Background())
	if err != nil {
		return []entity.User{}, 0, err
	}
	return users, total, nil
}

func (ru *RepositoryUsers) Update(user entity.User) error {
	return ru.repo.UpdateUser(context.Background(), repositoryDB.UpdateUserParams{
		ID:       user.ID.String(),
		Username: user.Username,
		Language: repositoryDB.UsersLanguage(user.Language),
	})
}

func (ru *RepositoryUsers) Desactive(userID uuid.UUID) error {
	return ru.repo.DeleteUser(context.Background(), userID.String())
}

func (ru *RepositoryUsers) FindByID(userID uuid.UUID) (entity.User, error) {
	repo, err := ru.repo.GetUserByID(context.Background(), userID.String())
	if err != nil {
		return entity.User{}, err
	}
	id, err := uuid.Parse(repo.ID)
	if err != nil {
		return entity.User{}, err
	}
	categories, err := findCategoriesByUserID(ru.repo, userID)
	if err != nil {
		return entity.User{}, err
	}
	chats, err := findChatsByUserID(ru.repo, userID)
	if err != nil {
		return entity.User{}, err
	}
	bots, err := findBotsByUserID(ru.repo, userID)
	if err != nil {
		return entity.User{}, err
	}
	var gender *entity.Gender
	gender = new(entity.Gender)
	if repo.Gender.Valid {
		*gender = entity.Gender(repo.Gender.UsersGender)
	}
	return entity.User{
		ID:         id,
		Username:   repo.Username,
		Email:      repo.Email,
		Password:   repo.Password,
		Platform:   entity.Platform(repo.Platform),
		Language:   entity.Language(repo.Language),
		Gender:     gender,
		Location:   repo.Location,
		Chats:      chats,
		Bots:       bots,
		CreatedAt:  repo.CreatedAt,
		Categories: categories,
	}, nil
}

func (ru *RepositoryUsers) FindUserByEmail(email string) (*entity.User, error) {
	userRepo, err := ru.repo.FindUserByEmail(context.Background(), email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return &entity.User{}, err
	}
	userID, err := uuid.Parse(userRepo.ID)
	if err != nil {
		return &entity.User{}, err
	}
	chats, err := findChatsByUserID(ru.repo, userID)
	if err != nil {
		return &entity.User{}, err
	}
	categories, err := findCategoriesByUserID(ru.repo, userID)
	if err != nil {
		return &entity.User{}, err
	}
	var gender *entity.Gender
	if userRepo.Gender.Valid {
		*gender = entity.Gender(string(userRepo.Gender.UsersGender))
	}
	return &entity.User{
		ID:         userID,
		Username:   userRepo.Username,
		Email:      userRepo.Email,
		AvatarURL:  userRepo.AvatarUrl.String,
		Password:   userRepo.Password,
		CreatedAt:  userRepo.CreatedAt,
		Platform:   entity.Platform(userRepo.Platform),
		Gender:     gender,
		Language:   entity.Language(userRepo.Language),
		Location:   userRepo.Location,
		Chats:      chats,
		Bots:       []entity.Bot{},
		Categories: categories,
	}, nil
}

func (ru *RepositoryUsers) CreateChat(chat entity.Chat) error {
	return ru.repo.CreateChat(context.Background(), repositoryDB.CreateChatParams{
		UserID: chat.UserID.String(),
		BotID:  chat.BotID.String(),
	})
}

func (ru *RepositoryUsers) DeleteCategoriesUserID(userID uuid.UUID) error {
	err := ru.repo.DeleteCategoriesUserID(context.Background(), userID.String())
	if err != nil {
		return err
	}
	return nil
}

func (ru *RepositoryUsers) AddCategoriesByUserID(user entity.User) error {
	tx, err := ru.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	repo := repositoryDB.New(tx)
	for _, category := range user.Categories {
		id, err := uuid.NewV7()
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("error rolling back transaction: %v", rbErr)
			}
			return err
		}
		err = repo.AddCategoriesByUserID(context.Background(), repositoryDB.AddCategoriesByUserIDParams{
			ID:         id.String(),
			UserID:     user.ID.String(),
			CategoryID: category.ID.String(),
		})
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("error rolling back transaction: %v", rbErr)
			}
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ru *RepositoryUsers) PublishBot(botID uuid.UUID, userID uuid.UUID) error {
	return ru.repo.PublishBot(context.Background(), repositoryDB.PublishBotParams{
		ID:     botID.String(),
		UserID: userID.String(),
	})
}

func (ru *RepositoryUsers) LikeToBot(userID, botID uuid.UUID) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	fmt.Println("id", userID.String())
	return ru.repo.LikeToBot(context.Background(), repositoryDB.LikeToBotParams{
		ID:     id.String(),
		UserID: userID.String(),
		BotID:  botID.String(),
	})
}

func (ru *RepositoryUsers) FindLikeBotByUserID(userID uuid.UUID) ([]entity.Bot, error) {
	repos, err := ru.repo.FindLikeBotByUserID(context.Background(), userID.String())
	if err != nil {
		return []entity.Bot{}, err
	}
	var bots []entity.Bot
	for _, repo := range repos {
		id, err := uuid.Parse(repo.ID)
		if err != nil {
			return []entity.Bot{}, err
		}
		bots = append(bots, entity.Bot{
			ID:            id,
			UserID:        userID,
			Name:          repo.Name,
			Description:   repo.Description,
			AvatarURL:     repo.AvatarUrl,
			BackgroundURL: repo.BackgroundUrl,
			Personality:   repo.Personality,
			Location:      repo.Location,
			Published:     repo.Published,
			CategoryID:    uuid.MustParse(repo.CategoryID),
			Active:        repo.Active,
			CreatedAt:     repo.CreatedAt,
			UpdatedAt:     repo.UpdatedAt,
		})
	}
	return bots, nil
}

func (ru *RepositoryUsers) FindLikeBotByBotID(botID uuid.UUID) ([]entity.User, error) {
	repos, err := ru.repo.FindLikeBotByBotID(context.Background(), botID.String())
	if err != nil {
		return []entity.User{}, err
	}
	var users []entity.User
	for _, repo := range repos {
		id, err := uuid.Parse(repo.ID)
		if err != nil {
			return []entity.User{}, err
		}
		var gender *entity.Gender
		if repo.Gender.Valid {
			*gender = entity.Gender(string(repo.Gender.UsersGender))
		}
		users = append(users, entity.User{
			ID:       id,
			Username: repo.Username,
			Email:    repo.Email,
			Password: repo.Password,
			Platform: entity.Platform(repo.Platform),
			Gender:   gender,
			Location: repo.Location,
		})
	}
	return users, nil
}

func (ru *RepositoryUsers) FindLikeBotByUserAndBot(userID uuid.UUID, botID uuid.UUID) (ok bool, err error) {
	_, err = ru.repo.FindLikeBotByUserAndBot(context.Background(), repositoryDB.FindLikeBotByUserAndBotParams{
		UserID: userID.String(),
		BotID:  botID.String(),
	})
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err != sql.ErrNoRows {
		return true, nil
	}
	return false, nil
}

func (ru *RepositoryUsers) UpdateAvatarURL(user entity.User) error {
	err := ru.repo.UpdateAvatarURL(context.Background(), repositoryDB.UpdateAvatarURLParams{
		AvatarUrl: sql.NullString{
			String: user.AvatarURL,
			Valid:  true,
		},
		ID: user.ID.String(),
	})
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func (ru *RepositoryUsers) ListMyBots(pagination entity.Pagination, userID uuid.UUID) ([]entity.Bot, int64, error) {
	repos, err := ru.repo.ListBotsByUserID(context.Background(), repositoryDB.ListBotsByUserIDParams{
		UserID: userID.String(),
		Limit:  pagination.PerPage,
		Offset: pagination.Offset,
	})
	if err != nil {
		return []entity.Bot{}, 0, err
	}
	bots := []entity.Bot{}
	for _, repo := range repos {
		id, err := uuid.Parse(repo.ID)
		if err != nil {
			return []entity.Bot{}, 0, err
		}
		categoryUUID, err := uuid.Parse(repo.CategoryID)
		if err != nil {
			return []entity.Bot{}, 0, err
		}
		bots = append(bots, entity.Bot{
			ID:            id,
			UserID:        userID,
			Name:          repo.Name,
			CategoryID:    categoryUUID,
			Description:   repo.Description,
			AvatarURL:     repo.AvatarUrl,
			BackgroundURL: repo.BackgroundUrl,
			Personality:   repo.Personality,
			Location:      repo.Location,
			Published:     repo.Published,
			Active:        repo.Active,
			Like:          repo.Likes,
			CreatedAt:     repo.CreatedAt,
			UpdatedAt:     repo.UpdatedAt,
		})
	}
	total, err := ru.repo.ListBotsByUserIDCount(context.Background(), userID.String())
	if err != nil {
		return []entity.Bot{}, 0, err
	}
	return bots, total, nil
}

func (ru *RepositoryUsers) ListCategoriesByUserID(userID uuid.UUID) ([]entity.Category, error) {
	repos, err := ru.repo.ListCategoriesByUserID(context.Background(), userID.String())
	if err != nil {
		return []entity.Category{}, err
	}
	categories := []entity.Category{}
	for _, repo := range repos {
		id, err := uuid.Parse(repo.ID)
		if err != nil {
			return []entity.Category{}, err
		}
		categories = append(categories, entity.Category{
			ID:        id,
			Name:      repo.Name,
			Active:    repo.Active,
			CreatedAt: repo.CreatedAt,
			UpdateAt:  repo.UpdatedAt,
		})
	}
	return categories, nil
}
