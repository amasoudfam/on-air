package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUsername string
	DBPassword string

	ServerPort string
}

var cfg *Config

func InitConfig() (*Config, error) {
	viper.SetConfigFile("config.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}

	cfg := &Config{
		DBHost:     viper.GetString("database.host"),
		DBPort:     viper.GetInt("database.port"),
		DBUsername: viper.GetString("database.username"),
		DBPassword: viper.GetString("database.password"),
		ServerPort: viper.GetString("server.port"),
	}

	return cfg, nil
}

func GetConfig() *Config {
	return cfg
}
