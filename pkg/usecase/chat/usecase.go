package chat

import (
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/gateway"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

type UseCase struct {
	RepoUser repository.RepositoryUser
	RepoChat repository.RepositoryChat
	Gateway  gateway.GatewayBot
	Cfg      *config.Config
}

func NewBotUseCase(repoUser repository.RepositoryUser, repoChat repository.RepositoryChat, gateway gateway.GatewayBot, cfg *config.Config) *UseCase {
	return &UseCase{
		RepoUser: repoUser,
		RepoChat: repoChat,
		Gateway:  gateway,
		Cfg:      cfg,
	}
}
