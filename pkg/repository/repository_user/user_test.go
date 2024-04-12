package repository_user_test

import (
	"testing"

	"github.com/infezek/app-chat/pkg/database"
	"github.com/infezek/app-chat/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		assert := assert.New(t)
		cfgTest := test.ConfigTest()
		db, err := database.New(cfgTest)
		assert.Nil(err)
		dbTest := test.NewDatabaseTest(db, cfgTest)
		dbTest.Delete()
		dbTest.Create()
	})

}
