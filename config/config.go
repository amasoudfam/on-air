package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Database Database
	Server   Server
	Redis    Redis
}

type Database struct {
	Host     string
	Port     int
	Username string
	Password string
	Name     string
}
type Redis struct {
	Host     string
	Port     int
	Password string
	DB       int
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
		},
		Redis: Redis{
			Host:     viper.GetString("redis.host"),
			Port:     viper.GetInt("redis.port"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
		Server: Server{
			Port: viper.GetString("server.port"),
		},
	}, nil
}
