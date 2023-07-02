package server

import (
	"fmt"
	"on-air/config"
	"on-air/server/handlers"
	"on-air/server/middlewares"
	"on-air/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupServer(cfg *config.Config, db *gorm.DB, redis *redis.Client, port string) error {
	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	auth := &handlers.Auth{
		DB:  db,
		JWT: &cfg.JWT,
	}

	e.POST("/auth/login", auth.Login)
	e.POST("/auth/register", auth.Register)

	authMiddleware := &middlewares.Auth{
		JWT: &cfg.JWT,
	}

	passenger := &handlers.Passenger{
		DB: db,
	}

	e.POST("/passenger", passenger.Create, authMiddleware.AuthMiddleware)
	e.GET("/passenger", passenger.Get, authMiddleware.AuthMiddleware)

	return e.Start(fmt.Sprintf(":%s", port))
}
