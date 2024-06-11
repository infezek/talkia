package usecase_user

import (
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/adapter"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type UploadImage struct {
	RepoUser     repository.RepositoryUser
	AdapterImage adapter.AdapterImagem
	Cfg          *config.Config
}

func NewUploadImage(cfg *config.Config, repoUser repository.RepositoryUser, adapter adapter.AdapterImagem) *UploadImage {
	return &UploadImage{
		RepoUser:     repoUser,
		AdapterImage: adapter,
		Cfg:          cfg,
	}
}

func (b *UploadImage) Execute(params UploadImageDtoInput) (UploadImageOutput, error) {
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		return UploadImageOutput{}, domain_error.BadRequest("invalid user id")
	}
	user, err := b.RepoUser.FindByID(userID)
	if err != nil {
		return UploadImageOutput{}, domain_error.NotFound("user not found")
	}
	id, err := uuid.NewV7()
	if err != nil {
		return UploadImageOutput{}, domain_error.BadRequest("error to generate id")
	}
	extensao := strings.Split(params.Avatar.Name, ".")
	nameFile := fmt.Sprintf("%s.%s", id.String(), extensao[len(extensao)-1])
	wg := sync.WaitGroup{}
	ch := make(chan error, 1)
	wg.Add(1)
	go b.AdapterImage.Upload(params.Avatar.File, id.String(), params.Avatar.Name, &wg, ch)
	wg.Wait()
	close(ch)
	var errors []string
	for err := range ch {
		if err != nil {
			errors = append(errors, err.Error())
		}
	}
	if len(errors) > 0 {
		return UploadImageOutput{}, domain_error.BadRequest(fmt.Sprintf("%s", strings.Join(errors, ",")))
	}
	user.UpdateAvatar(nameFile)
	if err := b.RepoUser.UpdateAvatarURL(user); err != nil {
		return UploadImageOutput{}, domain_error.BadRequest("error to update avatar")
	}
	return UploadImageOutput{
		AvatarURL: nameFile,
	}, nil
}
