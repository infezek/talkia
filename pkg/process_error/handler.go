package process_error

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func Handler(c *fiber.Ctx, err error) error {
	var e *fiber.Error
	code := fiber.StatusInternalServerError
	if errors.As(err, &e) {
		code = e.Code
	}
	return c.Status(code).JSON(map[string]string{"error": e.Message})
}
