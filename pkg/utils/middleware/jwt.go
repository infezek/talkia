package middleware

import (
	"database/sql"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/infezek/app-chat/pkg/domain/adapter"
	"github.com/infezek/app-chat/pkg/domain/domain_error"
	"github.com/infezek/app-chat/pkg/domain/repository"
)

func NewAuthMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
	})
}

type MiddlewareHandler struct {
	RepoUser     repository.RepositoryUser
	AdapterToken adapter.AdapterToken
}

func NewMiddlewareHandler(repoUser repository.RepositoryUser, adapterToken adapter.AdapterToken) *MiddlewareHandler {
	return &MiddlewareHandler{
		RepoUser:     repoUser,
		AdapterToken: adapterToken,
	}
}

func (m *MiddlewareHandler) FindUser() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		paramsUser, err := m.AdapterToken.DecodeToken(bearer)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		userUUID, err := uuid.Parse(paramsUser.UserID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		user, err := m.RepoUser.FindByID(userUUID)
		if err != nil && err != sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		if err == sql.ErrNoRows {
			return domain_error.NotFound("User not found")
		}
		if user.ID == uuid.Nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Next()
	}
}

type MiddlewareHandlerInterface interface {
	FindUser() func(c *fiber.Ctx) error
}
