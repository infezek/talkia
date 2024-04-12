package domain_error

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NotFound(message string) error {
	return fiber.NewError(http.StatusNotFound, message)
}

func Unprocessable(message []string) error {
	return fiber.NewError(http.StatusUnprocessableEntity, message...)
}

func Unauthorized(message ...string) error {
	msg := "unauthorized user"
	if len(message) == 1 {
		msg = strings.Join(message, ", ")
	}
	return fiber.NewError(http.StatusNotFound, msg)
}

func BadRequest(message string) error {
	return fiber.NewError(http.StatusBadRequest, message)
}

func Conflict(message string) error {
	return fiber.NewError(http.StatusConflict, message)
}
