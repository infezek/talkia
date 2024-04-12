package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Config struct {
	License                string
	AppName                string
	Enabled                bool
	TransportType          string
	Application            *newrelic.Application
	ErrorStatusCodeHandler func(c *fiber.Ctx, err error) int
}

var ConfigDefault = Config{
	Application:            nil,
	License:                "",
	AppName:                "fiber-api",
	Enabled:                false,
	ErrorStatusCodeHandler: DefaultErrorStatusCodeHandler,
}

func fiberNewRelic(cfg Config) fiber.Handler {
	var app *newrelic.Application
	var err error
	if cfg.ErrorStatusCodeHandler == nil {
		cfg.ErrorStatusCodeHandler = ConfigDefault.ErrorStatusCodeHandler
	}
	if cfg.Application != nil {
		app = cfg.Application
	} else {
		if cfg.AppName == "" {
			cfg.AppName = ConfigDefault.AppName
		}
		if cfg.License == "" {
			panic(fmt.Errorf("unable to create New Relic Application -> License can not be empty"))
		}
		app, err = newrelic.NewApplication(
			newrelic.ConfigAppName(cfg.AppName),
			newrelic.ConfigLicense(cfg.License),
			newrelic.ConfigEnabled(cfg.Enabled),
		)
		if err != nil {
			panic(fmt.Errorf("unable to create New Relic Application -> %w", err))
		}
	}
	return func(c *fiber.Ctx) error {
		txn := app.StartTransaction(createTransactionName(c))
		defer txn.End()
		var (
			host   = utils.CopyString(c.Hostname())
			method = utils.CopyString(c.Method())
		)
		scheme := c.Request().URI().Scheme()
		txn.SetWebRequest(newrelic.WebRequest{
			Host:      host,
			Method:    method,
			Transport: transport(string(scheme)),
			URL: &url.URL{
				Host:     host,
				Scheme:   string(c.Request().URI().Scheme()),
				Path:     string(c.Request().URI().Path()),
				RawQuery: string(c.Request().URI().QueryString()),
			},
		})
		headersMap := make(map[string]string)
		c.Request().Header.VisitAll(func(key, value []byte) {
			headersMap[string(key)] = string(value)
		})
		headersJSON, err := json.Marshal(headersMap)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		handlerErr := c.Next()
		statusCode := c.Context().Response.StatusCode()
		txn.AddAttribute("#header", string(headersJSON))
		txn.AddAttribute("#response", string(c.Response().Body()))
		txn.AddAttribute("#status-code", c.Response().StatusCode())
		txn.AddAttribute("#body", string(c.Body()))
		c.SetUserContext(newrelic.NewContext(c.UserContext(), txn))
		if handlerErr != nil {
			statusCode = cfg.ErrorStatusCodeHandler(c, handlerErr)
			txn.NoticeError(handlerErr)
		}
		txn.SetWebResponse(nil).WriteHeader(statusCode)
		return handlerErr
	}
}

func FromContext(c *fiber.Ctx) *newrelic.Transaction {
	return newrelic.FromContext(c.UserContext())
}

func createTransactionName(c *fiber.Ctx) string {
	return fmt.Sprintf("%s %s", c.Request().Header.Method(), c.Request().URI().Path())
}

func transport(schema string) newrelic.TransportType {
	if strings.HasPrefix(schema, "https") {
		return newrelic.TransportHTTPS
	}
	if strings.HasPrefix(schema, "http") {
		return newrelic.TransportHTTP
	}
	return newrelic.TransportUnknown
}

func DefaultErrorStatusCodeHandler(c *fiber.Ctx, err error) int {
	if fiberErr, ok := err.(*fiber.Error); ok {
		return fiberErr.Code
	}
	return c.Context().Response.StatusCode()
}
