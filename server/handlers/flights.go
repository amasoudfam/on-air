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
	FlightService *config.FlightService
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
	Origin        string `query:"org" param:"json" validate:"required"`
	Destination   string `query:"dest" param:"json" validate:"required"`
	Date          string `query:"date" param:"json" validate:"required,datetime=2006-01-02"`
	Airline       string `query:"airline"`
	Airplane      string `query:"airplane"`
	StartTime     string `query:"start_time" validate:"omitempty,CustomTimeValidator"`
	EndTime       string `query:"end_time" validate:"omitempty,CustomTimeValidator"`
	EmptyCapacity bool   `query:"empty_capacity"`
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
	cashFlights, err := f.Redis.Get(ctx.Request().Context(), redisKey).Result()

	if err != nil && err != redis.Nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	} else if err == redis.Nil {
		flights, err := services.GetFlightsListFromApi(f.Redis, f.FlightService, redisKey, ctx.Request().Context(), flightsList, req.Origin, req.Destination, req.Date)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}

		if len(flights) > 0 {
			jsonData, err := json.Marshal(flights)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}

			if err := f.Redis.Set(ctx.Request().Context(), redisKey, jsonData, time.Minute*10).Err(); err != nil {
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}

			flightsList = flights
		}

	} else {
		if err := json.Unmarshal([]byte(cashFlights), &flightsList); err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	// filters
	if req.Airline != "" {
		flightsList = filterByAirline(flightsList, req.Airline)
	}

	if req.Airplane != "" {
		flightsList = filterByAirplane(flightsList, req.Airplane)
	}

	if req.StartTime != "" {
		flightsList = filterByTime(flightsList, req.StartTime, req.EndTime)
	}

	if req.EmptyCapacity {
		flightsList = filterByCapacity(flightsList)
	}

	if req.OrderBy != "" {
		switch req.OrderBy {
		case "price":
			flightsList = sortByPrice(flightsList, req.SortOrder)
		case "time":
			flightsList = sortByTime(flightsList, req.SortOrder)
		case "duration":
			flightsList = sortByDuration(flightsList, req.SortOrder)
		}
	}

	return ctx.JSON(http.StatusOK, ListResponse{
		Flights: flightsList,
	})
}

func filterByAirline(flights []services.FlightDetails, airline string) []services.FlightDetails {
	filteredFlights := make([]services.FlightDetails, 0)
	for _, flight := range flights {
		if flight.Airline == airline {
			filteredFlights = append(filteredFlights, flight)
		}
	}

	return filteredFlights
}

func filterByAirplane(flights []services.FlightDetails, airplane string) []services.FlightDetails {
	filteredFlights := make([]services.FlightDetails, 0)
	for _, flight := range flights {
		if flight.Airplane == airplane {
			filteredFlights = append(filteredFlights, flight)
		}
	}

	return filteredFlights
}

func filterByTime(flights []services.FlightDetails, start_time, end_time string) []services.FlightDetails {
	startHourSplit := strings.Split(start_time, ":")
	startHour, _ := strconv.Atoi(startHourSplit[0])
	startMinute, _ := strconv.Atoi(startHourSplit[1])

	endTimeSplit := strings.Split(end_time, ":")
	endHour, _ := strconv.Atoi(endTimeSplit[0])
	endMinute, _ := strconv.Atoi(endTimeSplit[1])

	strings.Split(end_time, ":")
	filteredFlights := make([]services.FlightDetails, 0)
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

func filterByCapacity(flights []services.FlightDetails) []services.FlightDetails {
	filteredFlights := make([]services.FlightDetails, 0)
	for _, flight := range flights {
		if flight.EmptyCapacity > 0 {
			filteredFlights = append(filteredFlights, flight)
		}
	}

	return filteredFlights
}

func sortByPrice(flights []services.FlightDetails, sortOrder string) []services.FlightDetails {
	sort.Slice(flights, func(i, j int) bool {
		if sortOrder == "desc" {
			return flights[i].Price > flights[j].Price
		} else {
			return flights[i].Price < flights[j].Price
		}
	})

	return flights
}

func sortByTime(flights []services.FlightDetails, sortOrder string) []services.FlightDetails {
	sort.Slice(flights, func(i, j int) bool {
		if sortOrder == "desc" {
			return flights[i].StartedAt.After(flights[j].StartedAt)
		} else {
			return flights[i].StartedAt.Before(flights[j].StartedAt)
		}
	})

	return flights
}

func sortByDuration(flights []services.FlightDetails, sortOrder string) []services.FlightDetails {
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
