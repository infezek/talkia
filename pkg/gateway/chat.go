package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	histories := []string{}

	var mm []Message

	for _, preference := range chat.UserPreferences {
		if preference.PreferenceKey != "User" {
			continue
		}
		histories = append(histories, preference.PreferenceValue)
		mm = append(mm, Message{
			Role:    "user",
			Content: preference.PreferenceValue,
		})
	}
	historico := ""
	if len(histories) > 0 {
		//historico = fmt.Sprintf("historico [%s] ", strings.Join(histories, ", "))
	}
	language := "english"
	if user.Language.String() != "" {
		language = user.Language.String()
	}
	mm = append(mm, Message{
		Role: "user",
		Content: fmt.Sprintf(
			`%s Responda em %s, responda apenas assunto relacioando %s, (seu nome é  %s), simule uma conversa comigo, responda direto sem nenhum parametro antes como ("Resposta: ", "Darth Vader: ", "Eu: ", "Você: ", "Personagem: ") responda com texto resumido e der alguns detalhes na sua resposta, Responda a pergunta: %s`,
			historico,
			language,
			chat.Bot.Description,
			chat.Bot.Name,
			message),
	})
	reqBody := RequestBody{
		Model: "gpt-4-turbo-preview",
		ResponseFormat: map[string]string{
			"type": "text",
		},
		Messages: mm,
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
	return
}
