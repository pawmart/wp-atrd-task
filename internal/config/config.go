package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thanhpk/randstr"
)

var (
	DEBUG              bool
	REDIS_HOST         string
	REDIS_PASSWORD     string
	REDIS_KEY_PREFIX   string
	REDIS_DEFAULT_TIME int
	AES_KEY            string
	HTTP_DOCS_DIR      string
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("REDIS_HOST", "localhost:6379")
	REDIS_HOST = viper.GetString("REDIS_HOST")
	viper.SetDefault("REDIS_PASSWORD", "")
	REDIS_PASSWORD = viper.GetString("REDIS_PASSWORD")
	viper.SetDefault("REDIS_KEY_PREFIX", "")
	REDIS_KEY_PREFIX = viper.GetString("REDIS_KEY_PREFIX")
	viper.SetDefault("REDIS_DEFAULT_TIME", 600)
	REDIS_DEFAULT_TIME = viper.GetInt("REDIS_DEFAULT_TIME")

	viper.SetDefault("HTTP_DOCS_DIR", "./api/swagger/")
	HTTP_DOCS_DIR = viper.GetString("HTTP_DOCS_DIR")

	viper.SetDefault("AES_KEY", "changeme")
	AES_KEY = viper.GetString("AES_KEY")
	if AES_KEY == "changeme" {
		logrus.Warn("no aes key set, generated random one")
		AES_KEY = randstr.String(16)
		logrus.Debugf("aes key generated: %s", AES_KEY)
	}
}
