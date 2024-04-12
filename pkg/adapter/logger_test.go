package adapter_test

import (
	"testing"

	"github.com/infezek/app-chat/pkg/adapter"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	t.Run("verificando o tipo nil", func(t *testing.T) {
		assert := assert.New(t)
		log, err := adapter.NewLogger(nil, nil)
		assert.Nil(err)
		log.Info("TEST")
	})

	t.Run("asdasd", func(t *testing.T) {
		assert := assert.New(t)
		l := logrus.NewEntry(logrus.New())
		log, err := adapter.NewLogger(l, nil)
		assert.Nil(err)
		log.Info("TEST")
	})

	t.Run("teste", func(t *testing.T) {
		assert := assert.New(t)
		l := logrus.NewEntry(logrus.New())
		log, err := adapter.NewLogger(l, nil)
		log.AddParam("key", "value")
		log.AddParam("key2", "value2")
		assert.Nil(err)
		log.Info("TEST")
	})

}
