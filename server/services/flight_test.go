package services

import (
	"context"
	"encoding/json"
	"log"
	"on-air/config"
	"time"

	"net/http"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/jarcoal/httpmock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FlightServiceTestSuite struct {
	suite.Suite
	mockRedis redismock.ClientMock
	redis     *redis.Client
}

func (suite *FlightServiceTestSuite) SetupSuite() {
	mockRedis, mock := redismock.NewClientMock()
	suite.redis = mockRedis
	suite.mockRedis = mock
}

type ListResponse struct {
	Flights []FlightDetails `json:"flights"`
}

func (suite *FlightServiceTestSuite) TestGetFlightsListFromApi_Success() {
	data := []FlightDetails{
		{
			Number:  "FL001",
			Airline: "AirlineA",
		},
	}

	jsonData, _ := json.Marshal(data)
	suite.mockRedis.ExpectSet("redis_key", jsonData, time.Minute*10).SetVal(string(jsonData))
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expectedURL := "https://api.example.com/flights?org=origin&dest=destination&date=date"
	response := ListResponse{
		Flights: data,
	}
	res, _ := json.Marshal(response)
	expectedResponse := string(res)

	httpmock.RegisterResponder("GET", expectedURL, httpmock.NewStringResponder(http.StatusOK, expectedResponse))

	flightService := &config.FlightService{
		Url: "https://api.example.com",
	}

	flightsList := []FlightDetails{}
	flights, err := GetFlightsListFromApi(suite.redis, flightService, "redis_key", context.TODO(), flightsList, "origin", "destination", "date")
	if err != nil {
		log.Fatal(err)
	}

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), flights)
	assert.NotEmpty(suite.T(), flights)

	err = suite.mockRedis.ExpectationsWereMet()
	assert.NoError(suite.T(), err)
}

func (suite *FlightServiceTestSuite) TestGetFlightsListFromApi_webService_Failure() {
	data := []FlightDetails{
		{
			Number:  "FL001",
			Airline: "AirlineA",
		},
	}

	jsonData, _ := json.Marshal(data)
	suite.mockRedis.ExpectSet("redis_key", jsonData, time.Minute*10).SetVal(string(jsonData))
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expectedURL := "https://api.example.com/flights?org=origin&dest=destination&date=date"
	httpmock.RegisterResponder("GET", expectedURL, httpmock.NewStringResponder(http.StatusInternalServerError, "Internal Server Error"))

	flightService := &config.FlightService{
		Url: "https://api2.example.com",
	}

	flightsList := []FlightDetails{}
	_, err := GetFlightsListFromApi(suite.redis, flightService, "redis_key", context.TODO(), flightsList, "origin", "destination", "date")

	assert.NotNil(suite.T(), err)
}

func (suite *FlightServiceTestSuite) TestGetFlightsListFromApi_webService_500_Failure() {
	data := []FlightDetails{
		{
			Number:  "FL001",
			Airline: "AirlineA",
		},
	}

	jsonData, _ := json.Marshal(data)
	suite.mockRedis.ExpectSet("redis_key", jsonData, time.Minute*10).SetVal(string(jsonData))
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expectedURL := "https://api.example.com/flights?org=origin&dest=destination&date=date"
	httpmock.RegisterResponder("GET", expectedURL, httpmock.NewStringResponder(http.StatusInternalServerError, "Internal Server Error"))

	flightService := &config.FlightService{
		Url: "https://api.example.com",
	}

	flightsList := []FlightDetails{}
	_, err := GetFlightsListFromApi(suite.redis, flightService, "redis_key", context.TODO(), flightsList, "origin", "destination", "date")

	assert.EqualError(suite.T(), err, "web service returned an error")
}

func TestFlightService(t *testing.T) {
	suite.Run(t, new(FlightServiceTestSuite))
}
