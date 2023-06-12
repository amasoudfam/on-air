package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Database Database
	Server   Server
}

type Database struct {
	Host     string
	Port     int
	Username string
	Password string
	Name     string
}

type Server struct {
	Port string
}

func InitConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}

	return &Config{
		Database: Database{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			Username: viper.GetString("database.username"),
			Password: viper.GetString("database.password"),
			Name:     viper.GetString("database.name"),
		},
		Server: Server{
			Port: viper.GetString("server.port"),
		},
	}, nil
}
