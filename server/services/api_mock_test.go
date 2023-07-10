package services

import (
	"encoding/json"
	"errors"
	"time"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/jarcoal/httpmock"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type FlightServiceTestSuite struct {
	suite.Suite
	APIMockClient APIMockClient
}

func (suite *FlightServiceTestSuite) SetupSuite() {

	suite.APIMockClient = APIMockClient{
		Client:  &http.Client{},
		Breaker: &breaker.Breaker{},
		BaseURL: "http://example.com",
		Timeout: 10 * time.Second,
	}
}

type ListResponse struct {
	Flights []FlightResponse `json:"flights"`
}

func (suite *FlightServiceTestSuite) TestGetFlightsListFromApi_Success() {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(suite.T(), "/flights?origin=origin&destination=destination&date=date", r.URL.String())
		flights := []FlightResponse{
			{
				Number:        "FL001",
				Airplane:      "AirplaneA",
				Airline:       "AirlineA",
				Price:         100,
				Origin:        "OriginA",
				Destination:   "DestinationA",
				Capacity:      200,
				EmptyCapacity: 50,
				StartedAt:     time.Date(2023, 7, 1, 10, 0, 0, 0, time.UTC),
				FinishedAt:    time.Date(2023, 7, 1, 12, 0, 0, 0, time.UTC),
			},
		}
		respJSON, _ := json.Marshal(flights)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJSON)
	}))
	defer mockServer.Close()
	suite.APIMockClient.BaseURL = mockServer.URL
	expectedFlights := []FlightResponse{
		{
			Number:        "FL001",
			Airplane:      "AirplaneA",
			Airline:       "AirlineA",
			Price:         100,
			Origin:        "OriginA",
			Destination:   "DestinationA",
			Capacity:      200,
			EmptyCapacity: 50,
			StartedAt:     time.Date(2023, 7, 1, 10, 0, 0, 0, time.UTC),
			FinishedAt:    time.Date(2023, 7, 1, 12, 0, 0, 0, time.UTC),
		},
	}
	flights, err := suite.APIMockClient.GetFlights("origin", "destination", "date")
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), flights)
	require.Len(suite.T(), flights, 1)
	require.Equal(suite.T(), expectedFlights, flights)
}

type ErrorTransport struct{}

func (t *ErrorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("custom error")
}

func (suite *FlightServiceTestSuite) TestGetFlightsListFromApi_Unhandled_Response_Failure() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expectedURL := "http://example.com/flights?origin=origin&destination=destination&date=date"
	httpmock.RegisterResponder("GET", expectedURL, httpmock.NewStringResponder(http.StatusInternalServerError, "Internal Server Error"))

	flights, err := suite.APIMockClient.GetFlights("origin", "destination", "date")
	require.NotNil(suite.T(), err)
	require.Nil(suite.T(), flights)
}

func (suite *FlightServiceTestSuite) TestGetFlightsListFromApi_Request_Failure() {
	flights, err := suite.APIMockClient.GetFlights("origin", "destination", "date")
	require.NotNil(suite.T(), err)
	require.Nil(suite.T(), flights)
}

func (suite *FlightServiceTestSuite) TestGetFlightsCitiesFromApi_Success() {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(suite.T(), "/flights/cities", r.URL.String())
		cities := []string{"Shiraz", "Esfahan", "Tehran", "Tabriz"}
		respJSON, _ := json.Marshal(cities)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJSON)
	}))
	defer mockServer.Close()
	suite.APIMockClient.BaseURL = mockServer.URL
	expectedCities := []string{"Shiraz", "Esfahan", "Tehran", "Tabriz"}
	cities, err := suite.APIMockClient.GetCities()
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), cities)
	require.Len(suite.T(), cities, 4)
	require.Equal(suite.T(), expectedCities, cities)
}

func (suite *FlightServiceTestSuite) TestGetFlightsDatesFromApi_Success() {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(suite.T(), "/flights/dates", r.URL.String())
		cities := []string{"2023-03-07", "2023-03-08", "2023-03-10", "2023-03-11", "2023-03-20"}
		respJSON, _ := json.Marshal(cities)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJSON)
	}))
	defer mockServer.Close()
	suite.APIMockClient.BaseURL = mockServer.URL
	expectedDates := []string{"2023-03-07", "2023-03-08", "2023-03-10", "2023-03-11", "2023-03-20"}
	dates, err := suite.APIMockClient.GetDates()
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), dates)
	require.Len(suite.T(), dates, 5)
	require.Equal(suite.T(), expectedDates, dates)
}

func (suite *FlightServiceTestSuite) TestGetFlightFromApi_Success() {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(suite.T(), "/flights/FL001", r.URL.String())
		cities := FlightResponse{
			Number:        "FL001",
			Airplane:      "AirplaneA",
			Airline:       "AirlineA",
			Price:         100,
			Origin:        "OriginA",
			Destination:   "DestinationA",
			Capacity:      200,
			EmptyCapacity: 50,
			StartedAt:     time.Date(2023, 7, 1, 10, 0, 0, 0, time.UTC),
			FinishedAt:    time.Date(2023, 7, 1, 12, 0, 0, 0, time.UTC),
		}
		respJSON, _ := json.Marshal(cities)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJSON)
	}))
	defer mockServer.Close()
	suite.APIMockClient.BaseURL = mockServer.URL
	expectedFlight := FlightResponse{
		Number:        "FL001",
		Airplane:      "AirplaneA",
		Airline:       "AirlineA",
		Price:         100,
		Origin:        "OriginA",
		Destination:   "DestinationA",
		Capacity:      200,
		EmptyCapacity: 50,
		StartedAt:     time.Date(2023, 7, 1, 10, 0, 0, 0, time.UTC),
		FinishedAt:    time.Date(2023, 7, 1, 12, 0, 0, 0, time.UTC),
	}
	flight, err := suite.APIMockClient.GetFlight("FL001")
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), flight)
	require.Equal(suite.T(), &expectedFlight, flight)
}

func (suite *FlightServiceTestSuite) TestGetFlightReserveFromApi_Success() {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		require.Equal(suite.T(), "/flights/reserve", r.URL.String())
		cities := ReserveResponse{
			Status:  true,
			Message: "Flight reservation was successful.",
		}
		respJSON, _ := json.Marshal(cities)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJSON)
	}))
	defer mockServer.Close()
	suite.APIMockClient.BaseURL = mockServer.URL
	expectedRes := true
	reserveRes, err := suite.APIMockClient.Reserve("FL001", 1)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), reserveRes)
	require.Equal(suite.T(), expectedRes, reserveRes)
}

func (suite *FlightServiceTestSuite) TestGetFlightRefundFromApi_Success() {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		require.Equal(suite.T(), "/flights/refund", r.URL.String())
		cities := RefundResponse{
			Status:  true,
			Message: "Flight refund failed.",
		}
		respJSON, _ := json.Marshal(cities)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJSON)
	}))
	defer mockServer.Close()
	suite.APIMockClient.BaseURL = mockServer.URL
	expectedFlights := true
	reserveRes, err := suite.APIMockClient.Refund("FL001", 1)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), reserveRes)
	require.Equal(suite.T(), expectedFlights, reserveRes)
}

func TestFlightService(t *testing.T) {
	suite.Run(t, new(FlightServiceTestSuite))
}
