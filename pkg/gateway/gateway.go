package gateway

import "github.com/infezek/app-chat/pkg/config"

type GatewayBot struct {
	Cfg *config.Config
}

func NewGatewayBot(cfg *config.Config) *GatewayBot {
	return &GatewayBot{Cfg: cfg}
}
