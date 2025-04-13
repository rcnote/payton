package config

import (
	"github.com/spf13/viper"
	"log"
)

func Init() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config.yml")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading config file: %s\n", err.Error())
	}
}

func GetAppVersion() string {
	return "0.0.2"
}
