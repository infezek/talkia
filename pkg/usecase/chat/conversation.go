package chat

import (
	"sync"

	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/gateway"
)

type ConversationInput struct {
	Name        string
	Description string
	Text        string
}

func (u *UseCase) ConversationTest(input ConversationInput) (string, error) {
	user := entity.User{
		Language: "portuguese",
	}
	chat := entity.Chat{
		Bot: entity.Bot{
			Name:        input.Name,
			Description: input.Description,
		},
	}
	wg := sync.WaitGroup{}
	messageCh := make(chan gateway.GatewayResponse, 1)
	wg.Add(1)
	go u.Gateway.SendMessage(chat, user, input.Text, &wg, messageCh)
	wg.Wait()
	close(messageCh)
	for response := range messageCh {
		return response.Message.Message, nil
	}
	return "", nil
}
