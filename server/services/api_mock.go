package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"gorm.io/datatypes"
)

type APIMockClient struct {
	Client  *http.Client
	Breaker *breaker.Breaker
	BaseURL string
	Timeout time.Duration
}

type Penalties struct {
	Start   string
	End     string
	Percent int
}

type Penalty struct {
}
type FlightResponse struct {
	Number        string         `json:"number"`
	Airplane      string         `json:"airplane"`
	Airline       string         `json:"airline"`
	Price         int            `json:"price"`
	Origin        string         `json:"origin"`
	Destination   string         `json:"destination"`
	Capacity      int            `json:"capacity"`
	EmptyCapacity int            `json:"empty_capacity"`
	StartedAt     time.Time      `json:"started_at"`
	FinishedAt    time.Time      `json:"finished_at"`
	Penalties     datatypes.JSON `json:"penalties"`
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

func (c *APIMockClient) GetCities() ([]string, error) {
	url := c.BaseURL + "/flights/cities"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	var resp []string
	err = c.Breaker.Run(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
		defer cancel()

		req = req.WithContext(ctx)

		response, err := c.Client.Do(req)
		if err != nil {
			return fmt.Errorf("api_mock_get_flights_cities: request failed, error: %w", err)
		}

		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("api_mock_get_flights_cities: read response body failed, error: %w", err)
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("api_mock_get_flights_cities: unhandled response, status: %d, response: %s", response.StatusCode, responseBody)
		}

		err = json.Unmarshal(responseBody, &resp)
		if err != nil {
			return fmt.Errorf("api_mock_get_flights_cities: parse response body failed, response: %s, error: %w", responseBody, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *APIMockClient) GetDates() ([]string, error) {
	url := c.BaseURL + "/flights/dates"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	var resp []string
	err = c.Breaker.Run(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
		defer cancel()

		req = req.WithContext(ctx)

		response, err := c.Client.Do(req)
		if err != nil {
			return fmt.Errorf("api_mock_get_flights_dates: request failed, error: %w", err)
		}

		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("api_mock_get_flights_dates: read response body failed, error: %w", err)
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("api_mock_get_flights_dates: unhandled response, status: %d, response: %s", response.StatusCode, responseBody)
		}

		err = json.Unmarshal(responseBody, &resp)
		if err != nil {
			return fmt.Errorf("api_mock_get_flights_dates: parse response body failed, response: %s, error: %w", responseBody, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *APIMockClient) GetFlight(number string) (*FlightResponse, error) {
	url := c.BaseURL + "/flights/" + number

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	var resp FlightResponse
	err = c.Breaker.Run(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
		defer cancel()

		req = req.WithContext(ctx)

		response, err := c.Client.Do(req)
		if err != nil {
			return fmt.Errorf("api_mock_get_flight: request failed, error: %w", err)
		}

		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("api_mock_get_flight: read response body failed, error: %w", err)
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("api_mock_get_flight: unhandled response, status: %d, response: %s", response.StatusCode, responseBody)
		}

		err = json.Unmarshal(responseBody, &resp)
		if err != nil {
			return fmt.Errorf("api_mock_get_flight: parse response body failed, response: %s, error: %w", responseBody, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type ReserveResponse struct {
	Status  bool
	Message string
}

type ReserveRequestParameters struct {
	Number string
	Count  int
}

func (c *APIMockClient) Reserve(flightNumber string, ticketCount int) (bool, error) {
	baseUrl := c.BaseURL + "/flights/reserve"
	data := ReserveRequestParameters{
		Number: flightNumber,
		Count:  ticketCount,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return false, err
	}

	req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")

	var resp ReserveResponse
	err = c.Breaker.Run(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
		defer cancel()

		req = req.WithContext(ctx)

		response, err := c.Client.Do(req)
		if err != nil {
			return fmt.Errorf("api_mock_reserve_flight: request failed, error: %w", err)
		}

		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("api_mock_reserve_flight: read response body failed, error: %w", err)
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("api_mock_reserve_flight: unhandled response, status: %d, response: %s", response.StatusCode, responseBody)
		}

		err = json.Unmarshal(responseBody, &resp)
		if err != nil {
			return fmt.Errorf("api_mock_reserve_flight: parse response body failed, response: %s, error: %w", responseBody, err)
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return resp.Status, nil
}

type RefundResponse struct {
	Status  bool
	Message string
}

type RefundRequestParameters struct {
	Number string
	Count  int
}

func (c *APIMockClient) Refund(flightNumber string, ticketCount int) (bool, error) {
	baseUrl := c.BaseURL + "/flights/refund"
	data := RefundRequestParameters{
		Number: flightNumber,
		Count:  ticketCount,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return false, err
	}

	req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")

	var resp RefundResponse
	err = c.Breaker.Run(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
		defer cancel()

		req = req.WithContext(ctx)

		response, err := c.Client.Do(req)
		if err != nil {
			return fmt.Errorf("api_mock_refund_ticket: request failed, error: %w", err)
		}

		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("api_mock_refund_ticket: read response body failed, error: %w", err)
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("api_mock_refund_ticket: unhandled response, status: %d, response: %s", response.StatusCode, responseBody)
		}

		err = json.Unmarshal(responseBody, &resp)
		if err != nil {
			return fmt.Errorf("api_mock_refund_ticket: parse response body failed, response: %s, error: %w", responseBody, err)
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return resp.Status, nil
}
