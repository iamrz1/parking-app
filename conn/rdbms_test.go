package conn

import (
	"supertal-tha-parking-app/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRDBMSIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Run("no error", func(t *testing.T) {
		cnf := &config.RDBMS{
			DSN:     "file::memory:?cache=shared",
			Dialect: "sqlite",
		}

		db, cleanUp, err := Connect(cnf)
		defer cleanUp()
		assert.NoError(t, err)
		assert.NotNil(t, db)
	})

	t.Run("connection error", func(t *testing.T) {
		_, _, err := Connect(&config.RDBMS{})
		assert.Error(t, err)
	})
}
