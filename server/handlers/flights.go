package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"on-air/config"
	"on-air/server/services"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Flight struct {
	Redis         *redis.Client
	APIMockClient *services.APIMockClient
	Cache         *config.Redis
}

type FlightDetails struct {
	Number        string
	Airplane      string
	Airline       string
	Price         int
	Origin        string
	Destination   string
	Capacity      int
	EmptyCapacity int
	StartedAt     time.Time
	FinishedAt    time.Time
}

type ListRequest struct {
	Origin        string `query:"org" validate:"required"`
	Destination   string `query:"dest" validate:"required"`
	Date          string `query:"date" validate:"required,datetime=2006-01-02"`
	Airline       string `query:"airline"`
	Airplane      string `query:"airplane"`
	StartTime     string `query:"start_time" validate:"omitempty,CustomTimeValidator"`
	EndTime       string `query:"end_time" validate:"omitempty,CustomTimeValidator"`
	EmptyCapacity bool   `query:"empty_capacity"`
	OrderBy       string `query:"order_by"`
	SortOrder     string `query:"sort_order"`
}

type ListResponse struct {
	Flights []services.FlightResponse `json:"flights"`
}

func (f *Flight) List(ctx echo.Context) error {
	var req ListRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid query parameters")
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var flights []services.FlightResponse
	redisKey := fmt.Sprintf("flights_%s_%s_%s", req.Origin, req.Destination, req.Date)
	cashResult, err := f.Redis.Get(ctx.Request().Context(), redisKey).Result()

	if err != nil && err != redis.Nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	} else if err == redis.Nil {
		apiResult, err := f.APIMockClient.GetFlights(req.Origin, req.Destination, req.Date)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}

		if len(apiResult) > 0 {
			jsonData, err := json.Marshal(apiResult)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}

			if err := f.Redis.Set(ctx.Request().Context(), redisKey, jsonData, f.Cache.TTL).Err(); err != nil {
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}

			flights = apiResult
		}

	} else {
		if err := json.Unmarshal([]byte(cashResult), &flights); err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	if req.Airline != "" {
		flights = filterByAirline(flights, req.Airline)
	}

	if req.Airplane != "" {
		flights = filterByAirplane(flights, req.Airplane)
	}

	if req.StartTime != "" {
		flights = filterByTime(flights, req.StartTime, req.EndTime)
	}

	if req.EmptyCapacity {
		flights = filterByCapacity(flights)
	}

	if req.OrderBy != "" {
		switch req.OrderBy {
		case "price":
			flights = sortByPrice(flights, req.SortOrder)
		case "time":
			flights = sortByTime(flights, req.SortOrder)
		case "duration":
			flights = sortByDuration(flights, req.SortOrder)
		}
	}

	return ctx.JSON(http.StatusOK, ListResponse{
		Flights: flights,
	})
}

func filterByAirline(flights []services.FlightResponse, airline string) []services.FlightResponse {
	filteredFlights := make([]services.FlightResponse, 0)
	for _, flight := range flights {
		if flight.Airline == airline {
			filteredFlights = append(filteredFlights, flight)
		}
	}

	return filteredFlights
}

func filterByAirplane(flights []services.FlightResponse, airplane string) []services.FlightResponse {
	filteredFlights := make([]services.FlightResponse, 0)
	for _, flight := range flights {
		if flight.Airplane == airplane {
			filteredFlights = append(filteredFlights, flight)
		}
	}

	return filteredFlights
}

func filterByTime(flights []services.FlightResponse, startTime, endTime string) []services.FlightResponse {
	startTimeSplit := strings.Split(startTime, ":")
	startHour, _ := strconv.Atoi(startTimeSplit[0])
	startMinute, _ := strconv.Atoi(startTimeSplit[1])

	endTimeSplit := strings.Split(endTime, ":")
	endHour, _ := strconv.Atoi(endTimeSplit[0])
	endMinute, _ := strconv.Atoi(endTimeSplit[1])

	filteredFlights := make([]services.FlightResponse, 0)
	for _, flight := range flights {
		flightStartTime := flight.StartedAt

		if flightStartTime.Hour() > startHour ||
			(flightStartTime.Hour() == startHour && flightStartTime.Minute() >= startMinute) {
			if flightStartTime.Hour() < endHour ||
				(flightStartTime.Hour() == endHour && flightStartTime.Minute() <= endMinute) {
				filteredFlights = append(filteredFlights, flight)
			}
		}
	}

	return filteredFlights
}

func filterByCapacity(flights []services.FlightResponse) []services.FlightResponse {
	filteredFlights := make([]services.FlightResponse, 0)
	for _, flight := range flights {
		if flight.EmptyCapacity > 0 {
			filteredFlights = append(filteredFlights, flight)
		}
	}

	return filteredFlights
}

func sortByPrice(flights []services.FlightResponse, sortOrder string) []services.FlightResponse {
	sort.Slice(flights, func(i, j int) bool {
		if sortOrder == "desc" {
			return flights[i].Price > flights[j].Price
		} else {
			return flights[i].Price < flights[j].Price
		}
	})

	return flights
}

func sortByTime(flights []services.FlightResponse, sortOrder string) []services.FlightResponse {
	sort.Slice(flights, func(i, j int) bool {
		if sortOrder == "desc" {
			return flights[i].StartedAt.After(flights[j].StartedAt)
		} else {
			return flights[i].StartedAt.Before(flights[j].StartedAt)
		}
	})

	return flights
}

func sortByDuration(flights []services.FlightResponse, sortOrder string) []services.FlightResponse {
	sort.Slice(flights, func(i, j int) bool {
		durationA := flights[i].FinishedAt.Sub(flights[i].StartedAt)
		durationB := flights[j].FinishedAt.Sub(flights[j].StartedAt)
		if sortOrder == "asc" {
			return durationA < durationB
		} else {
			return durationA > durationB
		}
	})

	return flights
}
