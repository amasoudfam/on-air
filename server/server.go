package server

import (
	"fmt"
	"net/http"
	"on-air/config"
	"on-air/repository"
	"on-air/server/handlers"
	"on-air/server/middlewares"
	"on-air/server/services"
	"on-air/utils"

	"github.com/eapache/go-resiliency/breaker"
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
	customValidator := &utils.CustomValidator{
		Validator: validator.New(),
	}

	_ = customValidator.Validator.RegisterValidation("CustomTimeValidator", utils.CustomTimeValidator)
	e.Validator = customValidator

	apiMock := &services.APIMockClient{
		Client:  &http.Client{},
		Breaker: &breaker.Breaker{},
		BaseURL: cfg.Services.ApiMock.BaseURL,
		Timeout: cfg.Services.ApiMock.Timeout,
	}

	cityRepo := &repository.City{
		APIMockClient: apiMock,
		DB:            db,
		SyncPeriod:    cfg.Services.ApiMock.CitiesSyncPeriod,
	}

	go cityRepo.SyncCities()

	authMiddleware := &middlewares.Auth{
		JWT: &cfg.JWT,
	}

	auth := &handlers.Auth{
		DB:  db,
		JWT: &cfg.JWT,
	}

	e.POST("/auth/login", auth.Login)
	e.POST("/auth/register", auth.Register)

	ticket := &handlers.Ticket{
		DB:            db,
		JWT:           &cfg.JWT,
		APIMockClient: apiMock,
	}

	e.GET("/tickets", ticket.GetTickets, authMiddleware.AuthMiddleware)
	e.POST("/tickets/reserve", ticket.Reserve, authMiddleware.AuthMiddleware)
	e.GET("/tickets/pdf", ticket.GetPDF, authMiddleware.AuthMiddleware)

	payment := &handlers.Payment{
		DB:  db,
		IPG: &cfg.IPG,
	}

	e.POST("/payments/pay", payment.Pay, authMiddleware.AuthMiddleware)
	e.POST("/payments/callBack", payment.CallBack, authMiddleware.AuthMiddleware)

	flight := &handlers.Flight{
		Redis: redis,
		APIMockClient: &services.APIMockClient{
			Client:  &http.Client{},
			Breaker: &breaker.Breaker{},
			BaseURL: cfg.Services.ApiMock.BaseURL,
			Timeout: cfg.Services.ApiMock.Timeout,
		},
		Cache: &cfg.Redis,
	}

	e.GET("/flights", flight.GetFlights)

	passenger := &handlers.Passenger{
		DB: db,
	}

	e.GET("/passengers", passenger.Get, authMiddleware.AuthMiddleware)
	e.POST("/passengers", passenger.Create, authMiddleware.AuthMiddleware)

	return e.Start(fmt.Sprintf(":%s", port))
}
