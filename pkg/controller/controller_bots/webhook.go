package controller_bots

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/infezek/app-chat/db/migrate"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/database"
	"github.com/infezek/app-chat/pkg/domain/adapter"
	"github.com/infezek/app-chat/pkg/gateway"
	"github.com/infezek/app-chat/pkg/repository/repository_chat"
	"github.com/infezek/app-chat/pkg/repository/repository_user"
	"github.com/infezek/app-chat/pkg/usecase/chat"
	"github.com/infezek/app-chat/pkg/utils/middleware"
	"github.com/infezek/app-chat/pkg/utils/util_url"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
)

type Params struct {
	RepositoryUser interface{}
	GatewayChat    interface{}
}

var (
	Users map[*websocket.Conn]string = make(map[*websocket.Conn]string)
)

func WebSocket(app *fiber.App, adapterToken adapter.AdapterToken) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	cfg, err := config.New("chat")
	if err != nil {
		panic(err)
	}
	migrate.Run(cfg)
	db, err := database.New(cfg)
	if err != nil {
		panic(err)
	}
	repoUser := repository_user.New(db)
	repoChat := repository_chat.New(db)
	gateway := gateway.NewGatewayBot(cfg)
	usecaseSendMessage := chat.NewBotUseCase(repoUser, repoChat, gateway, cfg)
	jwt := middleware.NewAuthMiddleware(cfg.OpenIAToken)
	app.Get(util_url.New("/ws/chats/:chatID"), jwt, websocket.New(func(c *websocket.Conn) {
		var (
			mt  int
			msg []byte
		)

		for {
			chatID := c.Params("chatID")
			log := logrus.WithFields(logrus.Fields{"chatID": chatID})
			userID, ok := Users[c]
			if !ok {
				bearer := c.Headers("Authorization")
				paramsUser, err := adapterToken.DecodeToken(bearer)
				if err != nil {
					log.Error("error parsing token")
					continue
				}
				fmt.Println(paramsUser.UserID)
				Users[c] = paramsUser.UserID
			}
			log = log.WithFields(logrus.Fields{"userID": userID})
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Errorf("error reading message, %s", err)
				continue
			}
			message, err := usecaseSendMessage.SendMessage(chatID, Users[c], string(msg))
			if err != nil {
				log.Errorf("error sending message, %s", err)
				resp, _ := json.Marshal(ResponseError{Error: err.Error()})
				if err = c.WriteMessage(mt, resp); err != nil {
					log.Errorf("error writing message: %s", err)
				}
				continue
			}
			b, err := json.Marshal(Response{Message: message})
			if err != nil {
				log.Errorf("error marshaling message, %s", err)
				continue
			}
			if err = c.WriteMessage(mt, b); err != nil {
				log.Errorf("error writing message: %s", err)
			}
		}
	}))
}

type Response struct {
	Message string `json:"message"`
}

type ResponseError struct {
	Error string `json:"error"`
}
