package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"on-air/config"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
)

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

type ApiResponse struct {
	Flights []FlightDetails `json:"flights"`
}

func GetFlightsFromRedis(redisClient *redis.Client, ctx context.Context, redisKey string) (string, error) {
	return redisClient.Get(ctx, redisKey).Result()
}

func GetFlightsListFromApi(redisClient *redis.Client, flightService *config.FlightService, redisKey string, ctx context.Context, flightsList []FlightDetails, origin, destination, date string) ([]FlightDetails, error) {
	address := fmt.Sprintf("%s/%s", flightService.Url, "flights")
	url := fmt.Sprintf("%s?org=%s&dest=%s&date=%s", address, origin, destination, date)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode <= 599 {
		return nil, errors.New("web service returned an error")
	}

	body, _ := io.ReadAll(res.Body)

	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, err
	}

	if len(apiResponse.Flights) > 0 {
		jsonData, err := json.Marshal(apiResponse.Flights)
		if err != nil {
			return nil, err
		}

		if err := redisClient.Set(ctx, redisKey, jsonData, time.Minute*10).Err(); err != nil {
			return nil, err
		}

		flightsList = apiResponse.Flights
	}

	return flightsList, nil
}

func FilterByAirline(flights []FlightDetails, airline string) []FlightDetails {
	filteredFlights := make([]FlightDetails, 0)
	for _, flight := range flights {
		if flight.Airline == airline {
			filteredFlights = append(filteredFlights, flight)
		}
	}

	return filteredFlights
}

func FilterByAirplane(flights []FlightDetails, airplane string) []FlightDetails {
	filteredFlights := make([]FlightDetails, 0)
	for _, flight := range flights {
		if flight.Airplane == airplane {
			filteredFlights = append(filteredFlights, flight)
		}
	}

	return filteredFlights
}

func FilterByHour(flights []FlightDetails, hour int) []FlightDetails {
	filteredFlights := make([]FlightDetails, 0)
	for _, flight := range flights {
		if flight.StartedAt.Hour() == hour {
			filteredFlights = append(filteredFlights, flight)
		}
	}

	return filteredFlights
}

func FilterByCapacity(flights []FlightDetails) []FlightDetails {
	filteredFlights := make([]FlightDetails, 0)
	for _, flight := range flights {
		if flight.EmptyCapacity > 0 {
			filteredFlights = append(filteredFlights, flight)
		}
	}

	return filteredFlights
}

func SortByPrice(flights []FlightDetails, sortOrder string) []FlightDetails {
	sort.Slice(flights, func(i, j int) bool {
		if sortOrder == "desc" {
			return flights[i].Price > flights[j].Price
		} else {
			return flights[i].Price < flights[j].Price
		}
	})

	return flights
}

func SortByTime(flights []FlightDetails, sortOrder string) []FlightDetails {
	sort.Slice(flights, func(i, j int) bool {
		if sortOrder == "desc" {
			return flights[i].StartedAt.After(flights[j].StartedAt)
		} else {
			return flights[i].StartedAt.Before(flights[j].StartedAt)
		}
	})

	return flights
}

func SortByDuration(flights []FlightDetails, sortOrder string) []FlightDetails {
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
