package config

import (
	"fmt"

	"github.com/spf13/pflag"
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
}

type server struct {
	Port string
}

var cfg *Config

func InitConfig() (*Config, error) {

	configFile := pflag.String("config", "config.yaml", "Path to config file")
	pflag.Parse()

	viper.SetConfigFile(*configFile)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}

	cfg = &Config{
		Database: database{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			Username: viper.GetString("database.username"),
			Password: viper.GetString("database.password"),
		},
		Server: server{
			Port: viper.GetString("server.port"),
		},
	}

	return cfg, nil
}

func GetConfig() *Config {
	return cfg
}
