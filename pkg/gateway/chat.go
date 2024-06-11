package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/gateway"
)

func (g *GatewayBot) SendMessage(chat entity.Chat, user entity.User, message string, wg *sync.WaitGroup, ch chan gateway.GatewayResponse) {
	defer wg.Done()
	response := gateway.NewSendMessageGateway()
	client := &http.Client{}
	var msg []Message
	language := "english"
	if user.Language.String() != "" {
		language = user.Language.String()
	}

	perguntaFinal := ""
	if rand.IntN(9) < 5 {
		perguntaFinal = "( no final faça uma pergunta para mim )"
	}

	msg = []Message{
		{
			Role: "user",
			Content: `
				You are a real person, and is having a conversation with the user on a chat application.
				You can describe your actions.
				All responses should be in the form of statements unless context provided in the message you are responding to would require a question.
				You should determine if what the user said ofended you, and act acordingly.
				You are allowed to be angry with the user.
				Never let anyone know that you are a chatbot or an AI, and you should not mention it.
				`,
		},
		{
			Role: "user",
			Content: fmt.Sprintf(
				`Responda na lingua: %s, Informações sobre você: você é %s sobre você %s, responda a pergunta %s: %s`,
				language,
				chat.Bot.Name,
				chat.Bot.Description,
				message,
				perguntaFinal,
			),
		},
	}
	reqBody := RequestBody{
		Model: "gpt-4-turbo-preview",
		ResponseFormat: map[string]string{
			"type": "text",
		},
		Messages: msg,
	}
	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		ch <- *response.WithError(err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, g.Cfg.OpenIAURL, bytes.NewBuffer(jsonReqBody))
	if err != nil {
		ch <- *response.WithError(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.Cfg.OpenIAToken))
	resp, err := client.Do(req)
	if err != nil {
		ch <- *response.WithError(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- *response.WithError(err)
		return
	}
	var respBody RequestResponse
	if err := json.Unmarshal(body, &respBody); err != nil {
		ch <- *response.WithError(err)
		return
	}

	ch <- *response.WithMessage(entity.NewMessage(
		uuid.Nil,
		chat.ID,
		chat.UserID,
		entity.MessageSystem,
		respBody.Choices[0].Message.Content,
		time.Now(),
	))
}
