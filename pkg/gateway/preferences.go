package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/domain/gateway"
)

func (g *GatewayBot) ProcessPreferences(chat entity.Chat, message string, wg *sync.WaitGroup, ch chan gateway.GatewayResponse) {
	defer wg.Done()
	client := &http.Client{}
	histories := []string{}
	for _, preference := range chat.UserPreferences {
		if preference.PreferenceKey != "User" {
			continue
		}
		histories = append(histories, preference.PreferenceValue)
	}
	historico := ""
	hh := ""
	if len(histories) > 0 {
		historico = fmt.Sprintf("historico [%s] ", strings.Join(histories, ", "))
		hh = "exceto do histórico resuma."
	}
	mss := fmt.Sprintf(
		`%sSe a mensagem conter algo importante sobre o usuario como nome, idade, preferencia, gosto etc resuma a mensage tipo ('nome do usuario é', 'usuario gosta de', 'usuario tem ', 'usuario prefere', ) ,%s Caso contrário, responda com '####'. Mensagem do usuario: %s`,
		historico,
		hh,
		message)
	reqBody := RequestBody{
		Model: "gpt-4-turbo-preview",
		ResponseFormat: map[string]string{
			"type": "text",
		},
		Messages: []Message{
			{Role: "user",
				Content: mss},
		},
	}
	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		ch <- *gateway.NewProcessPreferencesGateway().WithError(err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, g.Cfg.OpenIAURL, bytes.NewBuffer(jsonReqBody))
	if err != nil {
		ch <- *gateway.NewProcessPreferencesGateway().WithError(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.Cfg.OpenIAToken))
	resp, err := client.Do(req)
	if err != nil {
		ch <- *gateway.NewProcessPreferencesGateway().WithError(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- *gateway.NewProcessPreferencesGateway().WithError(err)
		return
	}
	var respBody RequestResponse
	if err := json.Unmarshal(body, &respBody); err != nil {
		ch <- *gateway.NewProcessPreferencesGateway().WithError(err)
		return
	}
	ch <- *gateway.NewProcessPreferencesGateway().WithMessage(entity.Message{
		ChatID:  chat.ID,
		UserID:  chat.UserID,
		Who:     entity.MessageSystem,
		Message: respBody.Choices[0].Message.Content,
	})
	return
}
