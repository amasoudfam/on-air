package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/eapache/go-resiliency/breaker"
)

type APIMockClient struct {
	Client  *http.Client
	Breaker *breaker.Breaker
	BaseURL string
	Timeout time.Duration
}

type FlightResponse struct {
	Number        string    `json:"number"`
	Airplane      string    `json:"airplane"`
	Airline       string    `json:"airline"`
	Price         int       `json:"price"`
	Origin        string    `json:"origin"`
	Destination   string    `json:"destination"`
	Capacity      int       `json:"capacity"`
	EmptyCapacity int       `json:"empty_capacity"`
	StartedAt     time.Time `json:"started_at"`
	FinishedAt    time.Time `json:"finished_at"`
}

type City struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *APIMockClient) GetFlights(origin, destination, date string) ([]FlightResponse, error) {
	url := c.BaseURL + "/flights" + fmt.Sprintf("?origin=%s&destination=%s&date=%s", origin, destination, date)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	var resp []FlightResponse
	err = c.Breaker.Run(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
		defer cancel()

		req = req.WithContext(ctx)

		response, err := c.Client.Do(req)
		if err != nil {
			return fmt.Errorf("api_mock_get_flights: request failed, error: %w", err)
		}

		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("api_mock_get_flights: read response body failed, error: %w", err)
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("api_mock_get_flights: unhandled response, status: %d, response: %s", response.StatusCode, responseBody)
		}

		err = json.Unmarshal(responseBody, &resp)
		if err != nil {
			return fmt.Errorf("api_mock_get_flights: parse response body failed, response: %s, error: %w", responseBody, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
