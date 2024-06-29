package config

import (
	"os"
	"supertal-tha-parking-app/logger"
	"sync"

	"github.com/spf13/viper"
)

var cnfOnce = sync.Once{}

func read() {
	cnfOnce.Do(func() {
		configFileDir := os.Getenv(EnvConfigFileKey)
		if configFileDir == "" {
			configFileDir = DefaultConfigFileDir
		}
		viper.SetConfigFile(configFileDir)
		err := viper.ReadInConfig()
		if err != nil {
			logger.GetLogger().Infof("unable to read config file: %v", err)
		}
	})
}
