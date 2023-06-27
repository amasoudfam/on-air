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
	Origin        string `query:"org" param:"json" validate:"required"`
	Destination   string `query:"dest" param:"json" validate:"required"`
	Date          string `query:"date" param:"json" validate:"required,datetime=2006-01-02"`
	Airline       string `query:"Al"`
	Airplane      string `query:"AP"`
	Hour          int    `query:"HO"`
	EmptyCapacity bool   `query:"EC"`
	OrderBy       string `query:"order_by"`
	SortOrder     string `query:"sort_order"`
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

	if req.OrderBy != "" {
		switch req.OrderBy {
		case "price":
			flightsList = services.SortByPrice(flightsList, req.SortOrder)
		case "time":
			flightsList = services.SortByTime(flightsList, req.SortOrder)
		case "duration":
			flightsList = services.SortByDuration(flightsList, req.SortOrder)
		}
	}

	return ctx.JSON(http.StatusOK, ListResponse{
		Flights: flightsList,
	})
}
