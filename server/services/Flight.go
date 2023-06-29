package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"on-air/config"
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
