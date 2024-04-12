package usecase_bot

import (
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/adapter"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/repository"
	"github.com/sirupsen/logrus"
)

type UploadImage struct {
	RepoBot      repository.RepositoryBot
	AdapterImage adapter.AdapterImagem
	Cfg          *config.Config
}

func NewUploadImage(repoBot repository.RepositoryBot, adapter adapter.AdapterImagem, cfg *config.Config) *UploadImage {
	return &UploadImage{
		RepoBot:      repoBot,
		AdapterImage: adapter,
		Cfg:          cfg,
	}
}

func (b *UploadImage) Execute(params UploadImageDtoInput) error {
	logrus.Infof("[UploadImage] started")
	defer logrus.Infof("[UploadImage] finished")
	botID, err := uuid.Parse(params.BotID)
	if err != nil {
		logrus.Infof("[UploadImage] 1 %s", err.Error())
		return err
	}
	bot, err := b.RepoBot.FindByID(botID)
	if err != nil {
		logrus.Infof("[UploadImage] 2 %s", err.Error())
		return domain_error.NotFound("not found bot")
	}
	userUUID, err := uuid.Parse(params.UserID)
	if err != nil {
		logrus.Infof("[UploadImage] 3 %s", err.Error())
		return err
	}
	if bot.UserID != userUUID {
		logrus.Infof("[UploadImage] 4 %s", userUUID.String())
		return fmt.Errorf("user not authorized")
	}
	wg := new(sync.WaitGroup)
	ch := make(chan error, 2)
	ImageUUID := map[TypeImage]string{}
	logrus.Infof("[UploadImage] 5 %d", len(params.Files))
	for _, file := range params.Files {
		if file.File != nil {
			wg.Add(1)
			id, err := uuid.NewV7()
			if err != nil {
				return err
			}
			extensao := strings.Split(file.Name, ".")
			nameFile := fmt.Sprintf("%s.%s", id.String(), extensao[len(extensao)-1])
			ImageUUID[file.Type] = nameFile
			go b.AdapterImage.Upload(file.File, id.String(), file.Name, wg, ch)
		}
	}
	wg.Wait()
	close(ch)
	var errors []string
	for err := range ch {
		if err != nil {
			logrus.Infof("[UploadImage] 6 %s", err.Error())
			errors = append(errors, err.Error())
		}
	}
	if len(errors) > 0 {
		logrus.Infof("[UploadImage] 7 %s", strings.Join(errors, ","))
		return fmt.Errorf("%s", strings.Join(errors, ","))
	}
	for _, file := range params.Files {
		if file.Type == TypeImageAvatar {
			id, ok := ImageUUID[file.Type]
			if !ok {
				logrus.Infof("[UploadImage] 8")
				return fmt.Errorf("avatar image not found")
			}
			bot.UpdateAvatarURL(id)
		} else if file.Type == TypeImageBackground {
			id, ok := ImageUUID[file.Type]
			if !ok {
				logrus.Infof("[UploadImage] 9")
				return fmt.Errorf("background image not found")
			}
			bot.UpdateBackgroundURL(id)
		}
	}
	if err := b.RepoBot.Update(bot); err != nil {
		logrus.Infof("[UploadImage] 10 %s", err.Error())
		return err
	}
	return nil
}
