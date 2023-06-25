package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"on-air/config"
	"on-air/server/services"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Flight struct {
	DB            *gorm.DB
	Redis         *redis.Client
	FlightService *config.FlightService
}

type ListRequest struct {
	Origin        string `query:"org" validate:"required"`
	Destination   string `query:"dest" validate:"required"`
	Date          string `query:"date" validate:"required,datetime=2006-01-02"`
	Airline       string `query:"AL"`
	Airplane      string `query:"AP"`
	Hour          int    `query:"HO"`
	EmptyCapacity bool   `query:"EC"`
}

type ListResponse struct {
	Flights []services.FlightDetails `json:"flights"`
}

func (f *Flight) List(ctx echo.Context) error {
	var req ListRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid query parameters")
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var flightsList []services.FlightDetails
	redisKey := fmt.Sprintf("%s_%s_%s_%s", "flights", req.Origin, req.Destination, req.Date)
	cashFlights, err := services.GetFlightsFromRedis(f.Redis, ctx.Request().Context(), redisKey)
	if err != nil && err != redis.Nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	} else if err == redis.Nil {
		flightsList, err = services.GetFlightsListFromApi(f.Redis, f.FlightService, redisKey, ctx.Request().Context(), flightsList, req.Origin, req.Destination, req.Date)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		if err := json.Unmarshal([]byte(cashFlights), &flightsList); err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	// filters
	if req.Airline != "" {
		flightsList = services.FilterByAirline(flightsList, req.Airline)
	}

	if req.Airplane != "" {
		flightsList = services.FilterByAirplane(flightsList, req.Airplane)
	}

	if req.Hour != 0 {
		flightsList = services.FilterByHour(flightsList, req.Hour)
	}

	if req.EmptyCapacity {
		flightsList = services.FilterByCapacity(flightsList)
	}

	return ctx.JSON(http.StatusOK, flightsList)
}
