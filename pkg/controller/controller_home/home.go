package controller_home

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/infezek/app-chat/pkg/utils/util_url"
)

func Http(app *fiber.App) {
	app.Get(util_url.New("/"), func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{"message": "ok"})
	})
	app.Get(util_url.New("/metrics"), monitor.New(monitor.Config{Title: "talkIA"}))
}
