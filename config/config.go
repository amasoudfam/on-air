package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database Database
	Server   Server
	Redis    Redis
	JWT      JWT
	Services Services
}

type Database struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       string
}
type Redis struct {
	Host     string
	Port     int
	Password string
	DB       int
	TTL      time.Duration
}

type Server struct {
	Port string
}

type JWT struct {
	SecretKey string
	// TODO change name  expires_in
	LifeTime time.Duration
}

type Service struct {
	BaseURL string
	Timeout time.Duration
}

type Services struct {
	ApiMock Service
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
			DB:       viper.GetString("database.db"),
		},
		Redis: Redis{
			Host:     viper.GetString("redis.host"),
			Port:     viper.GetInt("redis.port"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
			TTL:      viper.GetDuration("redis.ttl"),
		},
		Server: Server{
			Port: viper.GetString("server.port"),
		},
		JWT: JWT{
			SecretKey: viper.GetString("auth.secret_key"),
			LifeTime:  viper.GetDuration("auth.lifetime"),
		},
		Services: Services{
			ApiMock: Service{
				BaseURL: viper.GetString("services.flights.url"),
				Timeout: viper.GetDuration("services.flights.timeout"),
			},
		},
	}, nil
}
