package config

import (
	"github.com/spf13/viper"
)

func NewConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	return viper.ReadInConfig()
}
