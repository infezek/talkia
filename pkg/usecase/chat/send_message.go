package chat

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/gateway"
	"github.com/sirupsen/logrus"
)

func (u *UseCase) SendMessage(chatID, userID string, text string) (string, error) {
	chatUUID, err := uuid.Parse(chatID)
	if err != nil {
		logrus.Error("error parse chat id")
		return "", err
	}
	chat, err := u.RepoChat.FindByID(chatUUID)
	if err != nil {
		return "", domain_error.NotFound("chat not found" + err.Error())
	}
	if chat.UserID != uuid.MustParse(userID) {
		return "", fmt.Errorf("user not allowed to send message")
	}
	user, err := u.RepoUser.FindByID(chat.UserID)
	if err != nil {
		return "", err
	}
	messageUser := entity.NewMessage(uuid.Nil, chat.ID, chat.UserID, entity.MessageUser, text, time.Now())
	err = u.RepoChat.SaveMessage(chat, messageUser)
	if err != nil {
		return "", err
	}
	wg := sync.WaitGroup{}
	messageCh := make(chan gateway.GatewayResponse, 2)
	wg.Add(2)
	go u.Gateway.SendMessage(chat, user, text, &wg, messageCh)
	go u.Gateway.ProcessPreferences(chat, text, &wg, messageCh)
	wg.Wait()
	close(messageCh)
	var messageSendMessage entity.Message
	var messageProcessPreferences entity.Message
	for response := range messageCh {
		if response.Error != nil {
			return "", response.Error
		}
		if response.Type == gateway.TypeSendMessage {
			messageSendMessage = response.Message
		}
		if response.Type == gateway.TypeProcessPreferences {
			messageProcessPreferences = response.Message
		}
	}
	if err := u.RepoChat.SaveMessage(chat, messageSendMessage); err != nil {
		return "", err
	}
	if err := u.savePreferences(chat, messageProcessPreferences); err != nil {
		return "", err
	}
	return messageSendMessage.Message, nil
}

func (u *UseCase) savePreferences(chat entity.Chat, messageProcessPreferences entity.Message) error {
	if len(strings.Split(messageProcessPreferences.Message, "##")) > 0 {
		return nil
	}
	if err := u.RepoChat.SavePreferences(chat, messageProcessPreferences); err != nil {
		return err
	}
	return nil
}
