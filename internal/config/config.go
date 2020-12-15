package config

import (
	"github.com/spf13/viper"
)

var (
	DEBUG            bool
	REDIS_HOST       string
	REDIS_PASSWORD   string
	REDIS_KEY_PREFIX string
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("REDIS_HOST", "localhost:6379")
	REDIS_HOST = viper.GetString("REDIS_HOST")
	viper.SetDefault("REDIS_PASSWORD", "")
	REDIS_PASSWORD = viper.GetString("REDIS_PASSWORD")
	viper.SetDefault("REDIS_KEY_PREFIX", "")
	REDIS_KEY_PREFIX = viper.GetString("REDIS_KEY_PREFIX")
}
