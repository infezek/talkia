package util_url_test

import (
	"os"
	"testing"

	"github.com/infezek/app-chat/pkg/utils/util_url"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	t.Run("Development", func(t *testing.T) {
		assert := assert.New(t)
		os.Setenv("CHAT_ENVIRONMENT", "development")
		os.Setenv("CHAT_PREFIX_URL", "/test")
		url := util_url.New("/users")
		assert.Equal("/users", url)
	})
	t.Run("url prod 1", func(t *testing.T) {
		assert := assert.New(t)
		os.Setenv("CHAT_ENVIRONMENT", "prod")
		os.Setenv("CHAT_PREFIX_URL", "/test")
		url := util_url.New("/users")
		assert.Equal("/test/users", url)
	})

	t.Run("1 url prod", func(t *testing.T) {
		assert := assert.New(t)
		os.Setenv("CHAT_ENVIRONMENT", "prod")
		os.Setenv("CHAT_PREFIX_URL", "/test/")
		url := util_url.New("/users")
		assert.Equal("/test/users", url)
	})
	t.Run("3 url prod ", func(t *testing.T) {
		assert := assert.New(t)
		os.Setenv("CHAT_ENVIRONMENT", "prod")
		os.Setenv("CHAT_PREFIX_URL", "/test/")
		url := util_url.New("users")
		assert.Equal("/test/users", url)
	})

	t.Run("4 url prod ", func(t *testing.T) {
		assert := assert.New(t)
		os.Setenv("CHAT_ENVIRONMENT", "prod")
		os.Setenv("CHAT_PREFIX_URL", "test")
		url := util_url.New("users")
		assert.Equal("/test/users", url)
	})

	t.Run("** 4 url prod ", func(t *testing.T) {
		assert := assert.New(t)
		os.Setenv("CHAT_ENVIRONMENT", "prod")
		os.Setenv("CHAT_PREFIX_URL", "talkia/v2")
		url := util_url.New("users")
		assert.Equal("/talkiav2/users", url)
	})

}
