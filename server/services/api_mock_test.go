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

func TestFlightService(t *testing.T) {
	suite.Run(t, new(FlightServiceTestSuite))
}
