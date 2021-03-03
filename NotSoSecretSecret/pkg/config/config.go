package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config variables
var (
	DbName string
	DbURI  string
	Port   string
)

func init() {

	viper.SetDefault("DB_NAME", "very-secret-db")
	viper.SetDefault("DB_URI", "mongodb://mongo:27017")
	viper.SetDefault("PORT", "7777")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("No config file found, default varbiables will be load")
	}

	DbName = viper.GetString("DB_NAME")
	DbURI = viper.GetString("DB_URI")
	Port = viper.GetString("PORT")
}
