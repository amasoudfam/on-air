package services

import (
	"context"
	"encoding/json"
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
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

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

// [91 123 34 78 117 109 98 101 114 34 58 34 70 76 49 50 51 34 44 34 65 105 114 112 108 97 110 101 34 58 34 34 44 34 65 105 114 108 105 110 101 34 58 34 69 120 97 109 112 108 101 32 65 105 114 108 105 110 101 115 34 44 34 80 114 105 99 101 34 58 49 48 48 44 34 79 114 105 103 105 110 34 58 34 34 44 34 68 101 115 116 105 110 97 116 105 111 110 34 58 34 34 44 34 67 97 112 97 99 105 116 121 34 58 48 44 34 69 109 112 116 121 67 97 112 97 99 105 116 121 34 58 48 44 34 83 116 97 114 116 101 100 65 116 34 58 34 48 48 48 49 45 48 49 45 48 49 84 48 48 58 48 48 58 48 48 90 34 44 34 70 105 110 105 115 104 101 100 65 116 34 58 34 48 48 48 49 45 48 49 45 48 49 84 48 48 58 48 48 58 48 48 90 34 125 93]
// [91 123 34 78 117 109 98 101 114 34 58 34 70 76 48 48 49 34 44 34 65 105 114 112 108 97 110 101 34 58 34 34 44 34 65 105 114 108 105 110 101 34 58 34 65 105 114 108 105 110 101 65 34 44 34 80 114 105 99 101 34 58 48 44 34 79 114 105 103 105 110 34 58 34 34 44 34 68 101 115 116 105 110 97 116 105 111 110 34 58 34 34 44 34 67 97 112 97 99 105 116 121 34 58 48 44 34 69 109 112 116 121 67 97 112 97 99 105 116 121 34 58 48 44 34 83 116 97 114 116 101 100 65 116 34 58 34 48 48 48 49 45 48 49 45 48 49 84 48 48 58 48 48 58 48 48 90 34 44 34 70 105 110 105 115 104 101 100 65 116 34 58 34 48 48 48 49 45 48 49 45 48 49 84 48 48 58 48 48 58 48 48 90 34 125 93]
