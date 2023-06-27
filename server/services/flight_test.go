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

func (suite *FlightServiceTestSuite) TestGetFlightsFromRedis_Success() {
	suite.mockRedis.ExpectGet("redis_key").
		SetVal(`[{"Number": "FL001", "Airline": "AirlineA"}]`)

	flights, err := GetFlightsFromRedis(suite.redis, context.TODO(), "redis_key")

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), `[{"Number": "FL001", "Airline": "AirlineA"}]`, flights)

	err = suite.mockRedis.ExpectationsWereMet()
	assert.NoError(suite.T(), err)
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

func (suite *FlightServiceTestSuite) TestFilterByAirline() {
	flights := []FlightDetails{
		{Number: "FL001", Airline: "AirlineA"},
		{Number: "FL002", Airline: "AirlineB"},
		{Number: "FL003", Airline: "AirlineA"},
		{Number: "FL004", Airline: "AirlineC"},
	}
	airline := "AirlineA"
	filteredFlights := FilterByAirline(flights, airline)

	expectedFlights := []FlightDetails{
		{Number: "FL001", Airline: "AirlineA"},
		{Number: "FL003", Airline: "AirlineA"},
	}
	assert.Len(suite.T(), filteredFlights, 2)
	assert.Equal(suite.T(), expectedFlights, filteredFlights)
}

func (suite *FlightServiceTestSuite) TestFilterByAirplane() {
	flights := []FlightDetails{
		{Number: "FL001", Airline: "AirlineA", Airplane: "bb"},
		{Number: "FL002", Airline: "AirlineB", Airplane: "hh"},
		{Number: "FL003", Airline: "AirlineA", Airplane: "bb"},
		{Number: "FL004", Airline: "AirlineC", Airplane: "cc"},
	}
	airplane := "hh"
	filteredFlights := FilterByAirplane(flights, airplane)

	expectedFlights := []FlightDetails{
		{Number: "FL002", Airline: "AirlineB", Airplane: "hh"},
	}
	assert.Len(suite.T(), filteredFlights, 1)
	assert.Equal(suite.T(), expectedFlights, filteredFlights)
}

func (suite *FlightServiceTestSuite) TestFilterByCapacity() {
	flights := []FlightDetails{
		{Number: "FL001", Airline: "AirlineA", EmptyCapacity: 0},
		{Number: "FL002", Airline: "AirlineB", EmptyCapacity: 0},
		{Number: "FL003", Airline: "AirlineA", EmptyCapacity: 12},
		{Number: "FL004", Airline: "AirlineC", EmptyCapacity: 2},
	}
	filteredFlights := FilterByCapacity(flights)

	expectedFlights := []FlightDetails{
		{Number: "FL003", Airline: "AirlineA", EmptyCapacity: 12},
		{Number: "FL004", Airline: "AirlineC", EmptyCapacity: 2},
	}
	assert.Len(suite.T(), filteredFlights, 2)
	assert.Equal(suite.T(), expectedFlights, filteredFlights)
}

func (suite *FlightServiceTestSuite) TestFilterByHour() {
	flights := []FlightDetails{
		{Number: "FL001", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC)},
		{Number: "FL002", Airline: "AirlineB", StartedAt: time.Date(2023, 6, 1, 11, 0, 0, 0, time.UTC)},
		{Number: "FL003", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 30, 0, 0, time.UTC)},
		{Number: "FL004", Airline: "AirlineC", StartedAt: time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)},
	}
	hour := 10
	filteredFlights := FilterByHour(flights, hour)

	expectedFlights := []FlightDetails{
		{Number: "FL001", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC)},
		{Number: "FL003", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 30, 0, 0, time.UTC)},
	}
	assert.Len(suite.T(), filteredFlights, 2)
	assert.Equal(suite.T(), expectedFlights, filteredFlights)
}

func (suite *FlightServiceTestSuite) TestSortByPrice_desc() {
	flights := []FlightDetails{
		{Number: "FL001", Airline: "AirlineA", Price: 1800000},
		{Number: "FL002", Airline: "AirlineB", Price: 1300000},
		{Number: "FL003", Airline: "AirlineA", Price: 1900000},
		{Number: "FL004", Airline: "AirlineC", Price: 1100000},
	}
	filteredFlights := SortByPrice(flights, "desc")

	expectedFlights := []FlightDetails{
		{Number: "FL003", Airline: "AirlineA", Price: 1900000},
		{Number: "FL001", Airline: "AirlineA", Price: 1800000},
		{Number: "FL002", Airline: "AirlineB", Price: 1300000},
		{Number: "FL004", Airline: "AirlineC", Price: 1100000},
	}

	assert.Equal(suite.T(), expectedFlights, filteredFlights)
}

func (suite *FlightServiceTestSuite) TestSortByPrice_acs() {
	flights := []FlightDetails{
		{Number: "FL001", Airline: "AirlineA", Price: 1800000},
		{Number: "FL002", Airline: "AirlineB", Price: 1300000},
		{Number: "FL003", Airline: "AirlineA", Price: 1900000},
		{Number: "FL004", Airline: "AirlineC", Price: 1100000},
	}
	sortFlights := SortByPrice(flights, "asc")

	expectedFlights := []FlightDetails{
		{Number: "FL004", Airline: "AirlineC", Price: 1100000},
		{Number: "FL002", Airline: "AirlineB", Price: 1300000},
		{Number: "FL001", Airline: "AirlineA", Price: 1800000},
		{Number: "FL003", Airline: "AirlineA", Price: 1900000},
	}

	assert.Equal(suite.T(), expectedFlights, sortFlights)
}

func (suite *FlightServiceTestSuite) TestSortByTime_asc() {
	flights := []FlightDetails{
		{Number: "FL001", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC)},
		{Number: "FL003", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 2, 10, 30, 0, 0, time.UTC)},
		{Number: "FL002", Airline: "AirlineB", StartedAt: time.Date(2023, 6, 1, 11, 0, 0, 0, time.UTC)},
		{Number: "FL004", Airline: "AirlineC", StartedAt: time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)},
	}
	sortFlights := SortByTime(flights, "asc")

	expectedFlights := []FlightDetails{
		{Number: "FL001", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC)},
		{Number: "FL002", Airline: "AirlineB", StartedAt: time.Date(2023, 6, 1, 11, 0, 0, 0, time.UTC)},
		{Number: "FL003", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 2, 10, 30, 0, 0, time.UTC)},
		{Number: "FL004", Airline: "AirlineC", StartedAt: time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)},
	}

	assert.Equal(suite.T(), expectedFlights, sortFlights)
}

func (suite *FlightServiceTestSuite) TestSortByTime_desc() {
	flights := []FlightDetails{
		{Number: "FL001", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC)},
		{Number: "FL003", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 2, 10, 30, 0, 0, time.UTC)},
		{Number: "FL002", Airline: "AirlineB", StartedAt: time.Date(2023, 6, 1, 11, 0, 0, 0, time.UTC)},
		{Number: "FL004", Airline: "AirlineC", StartedAt: time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)},
	}
	sortFlights := SortByTime(flights, "desc")

	expectedFlights := []FlightDetails{
		{Number: "FL004", Airline: "AirlineC", StartedAt: time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)},
		{Number: "FL003", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 2, 10, 30, 0, 0, time.UTC)},
		{Number: "FL002", Airline: "AirlineB", StartedAt: time.Date(2023, 6, 1, 11, 0, 0, 0, time.UTC)},
		{Number: "FL001", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC)},
	}

	assert.Equal(suite.T(), expectedFlights, sortFlights)
}

func (suite *FlightServiceTestSuite) TestSortByDuration_asc() {
	flights := []FlightDetails{
		{Number: "FL001", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 11, 23, 0, 0, time.UTC)},
		{Number: "FL003", StartedAt: time.Date(2023, 6, 2, 8, 30, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 11, 50, 0, 0, time.UTC)},
		{Number: "FL002", StartedAt: time.Date(2023, 6, 1, 06, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 10, 20, 0, 0, time.UTC)},
		{Number: "FL004", StartedAt: time.Date(2023, 6, 2, 05, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 06, 0, 0, 0, time.UTC)},
	}
	sortFlights := SortByDuration(flights, "asc")

	expectedFlights := []FlightDetails{
		{Number: "FL004", StartedAt: time.Date(2023, 6, 2, 05, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 06, 0, 0, 0, time.UTC)},
		{Number: "FL001", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 11, 23, 0, 0, time.UTC)},
		{Number: "FL003", StartedAt: time.Date(2023, 6, 2, 8, 30, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 11, 50, 0, 0, time.UTC)},
		{Number: "FL002", StartedAt: time.Date(2023, 6, 1, 06, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 10, 20, 0, 0, time.UTC)},
	}

	assert.Equal(suite.T(), expectedFlights, sortFlights)
}

func (suite *FlightServiceTestSuite) TestSortByDuration_desc() {
	flights := []FlightDetails{
		{Number: "FL001", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 11, 23, 0, 0, time.UTC)},
		{Number: "FL003", StartedAt: time.Date(2023, 6, 2, 8, 30, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 11, 50, 0, 0, time.UTC)},
		{Number: "FL002", StartedAt: time.Date(2023, 6, 1, 06, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 10, 20, 0, 0, time.UTC)},
		{Number: "FL004", StartedAt: time.Date(2023, 6, 2, 05, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 06, 0, 0, 0, time.UTC)},
	}
	sortFlights := SortByDuration(flights, "desc")

	expectedFlights := []FlightDetails{
		{Number: "FL002", StartedAt: time.Date(2023, 6, 1, 06, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 10, 20, 0, 0, time.UTC)},
		{Number: "FL003", StartedAt: time.Date(2023, 6, 2, 8, 30, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 11, 50, 0, 0, time.UTC)},
		{Number: "FL001", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 11, 23, 0, 0, time.UTC)},
		{Number: "FL004", StartedAt: time.Date(2023, 6, 2, 05, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 06, 0, 0, 0, time.UTC)},
	}

	assert.Equal(suite.T(), expectedFlights, sortFlights)
}

func TestFlightService(t *testing.T) {
	suite.Run(t, new(FlightServiceTestSuite))
}
