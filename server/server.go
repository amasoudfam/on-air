package server

import (
	"fmt"
	"net/http"
	"on-air/config"
	"on-air/server/handlers"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func SetupServer(cfg *config.Config, db *gorm.DB, redis *redis.Client, port string) error {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	auth := &handlers.Auth{
		DB:  db,
		JWT: &cfg.JWT,
	}

	e.POST("/auth/login", auth.Login)
	e.POST("/auth/register", auth.Register)

	ticket := &handlers.Ticket{
		DB:  db,
		JWT: &cfg.JWT,
	}

	e.POST("/ticket/reserve", ticket.Reserve)

	payment := &handlers.Payment{
		DB:  db,
		JWT: &cfg.JWT,
		IPG: &cfg.IPG,
	}

	e.POST("/Payment/reserve", payment.Pay)
	e.POST("/Payment/callBack", payment.CallBack)

	return e.Start(fmt.Sprintf(":%s", port))
}
