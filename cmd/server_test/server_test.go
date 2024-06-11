package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infezek/app-chat/pkg/domain/entity"
	"github.com/infezek/app-chat/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestShould(t *testing.T) {
	t.Run("should to return message ok", func(t *testing.T) {
		assert := assert.New(t)
		params := test.Implementations(test.ConfigTest())
		app := test.Server(params)
		req := httptest.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		assert.Nil(err)
		assert.Equal(200, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		assert.Nil(err)
		assert.Equal(`{"message":"ok"}`, string(body))
	})
	t.Run("should create user", func(t *testing.T) {
		assert := assert.New(t)
		params := test.Implementations(test.ConfigTest())
		db := test.NewDatabaseTest(params.DB, test.ConfigTest())
		db.Trucate()
		_, total, err := params.RepoUser.List(entity.Pagination{
			Page:    1,
			PerPage: 10,
		})
		assert.Equal(int64(0), total)
		assert.Nil(err)
		app := test.Server(params)
		b, _ := json.Marshal(map[string]string{
			"username": "ezequiel lopes",
			"email":    "ezequie.jn@outlook.com",
			"password": "123456645",
			"platform": "ios",
			"location": "test",
			"language": "english",
		})
		fmt.Println(string(b))
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.Equal(http.StatusOK, resp.StatusCode)
	})
}
