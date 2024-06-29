package config

import (
	"os"
	"sync"
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func TestRDBMS(t *testing.T) {
	t.Run("load config", func(t *testing.T) {
		cnfOnce = sync.Once{}
		rdbmsOnce = sync.Once{}
		viper.Reset()
		os.Setenv(EnvConfigFileKey, "../config.example.yaml")

		cnf1, err := RDBMSCnf()
		assert.NoError(t, err)

		cnf2, err := RDBMSCnf()
		assert.NoError(t, err)

		assert.Equal(t, cnf1, cnf2)
	})

	t.Run("config error", func(t *testing.T) {
		rdbmsOnce = sync.Once{}
		viper.Reset()

		_, err := RDBMSCnf()
		assert.EqualError(t, err, `rdbms configuration error: {"dialect":"required","dsn":"required"}`)
	})
}
