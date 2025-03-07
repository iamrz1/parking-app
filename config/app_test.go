package config

import (
	"os"
	"sync"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	t.Run("load config", func(t *testing.T) {
		appOnce = sync.Once{}
		cnfOnce = sync.Once{}
		viper.Reset()
		os.Setenv(EnvConfigFileKey, "../config.example.yaml")

		cnf1, err := AppCnf()
		assert.NoError(t, err)

		cnf2, err := AppCnf()
		assert.NoError(t, err)

		assert.Equal(t, cnf1, cnf2)
	})

	t.Run("config error", func(t *testing.T) {
		appOnce = sync.Once{}
		viper.Reset()

		_, err := AppCnf()
		assert.EqualError(t, err, `app configuration error: {"daemon_key":"required","server_port":"required","system_port":"required","timeout":"required"}`)
	})
}
