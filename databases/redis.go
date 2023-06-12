package databases

import (
	"context"
	"fmt"
	"log"
	"on-air/config"

	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) *redis.Client {

	address := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return client

}
