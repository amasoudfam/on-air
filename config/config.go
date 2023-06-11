package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Database database
	Server   server
}

type database struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
}

type server struct {
	Port string
}

var cfg *Config

func InitConfig(configPath string) error {
	viper.SetConfigFile(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config file: %s", err)
	}

	cfg = &Config{
		Database: database{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			Username: viper.GetString("database.username"),
			Password: viper.GetString("database.password"),
			DbName:   viper.GetString("database.dbname"),
		},
		Server: server{
			Port: viper.GetString("server.port"),
		},
	}

	return nil
}

func GetConfig() *Config {
	return cfg
}
