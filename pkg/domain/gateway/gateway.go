package gateway

import (
	"sync"

	"github.com/infezek/app-chat/pkg/domain/entity"
)

type GatewayResponse struct {
	Type    GatewayType
	Message entity.Message
	Error   error
}

func (g *GatewayResponse) WithError(err error) *GatewayResponse {
	g.Error = err
	return g
}

func (g *GatewayResponse) WithMessage(message entity.Message) *GatewayResponse {
	g.Message = message
	return g
}

type GatewayType string

var (
	TypeSendMessage        GatewayType = "SendMessage"
	TypeProcessPreferences GatewayType = "ProcessPreferences"
)

func NewSendMessageGateway() *GatewayResponse {
	return &GatewayResponse{
		Type: TypeSendMessage,
	}
}

func NewProcessPreferencesGateway() *GatewayResponse {
	return &GatewayResponse{
		Type: TypeProcessPreferences,
	}
}

type GatewayBot interface {
	SendMessage(entityChat entity.Chat, user entity.User, message string, wg *sync.WaitGroup, ch chan GatewayResponse)
	ProcessPreferences(entityChat entity.Chat, message string, wg *sync.WaitGroup, ch chan GatewayResponse)
}
