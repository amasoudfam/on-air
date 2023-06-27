package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"on-air/config"
	"on-air/server/services"
	"testing"

	"bou.ke/monkey"
	"github.com/go-redis/redismock/v9"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type MockContextWithTrueValidation struct {
	echo.Context
}

func (m *MockContextWithTrueValidation) Validate(i interface{}) error {
	return nil
}

type MockContextWithErrorValidation struct {
	echo.Context
}

func (m *MockContextWithErrorValidation) Validate(i interface{}) error {
	return errors.New("some parameter are needed")
}

type FlightHandlerTestSuite struct {
	suite.Suite
	mockRedis redismock.ClientMock
	redis     *redis.Client
	e         *echo.Echo
	endpoint  string
	flight    *Flight
}

func (suite *FlightHandlerTestSuite) CallHandler(queryString string) (*http.Request, *httptest.ResponseRecorder) {

	url := suite.endpoint + queryString

	req := httptest.NewRequest(http.MethodGet, url, nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	res := httptest.NewRecorder()

	return req, res
}

func (suite *FlightHandlerTestSuite) SetupSuite() {
	mockRedis, mock := redismock.NewClientMock()
	suite.redis = mockRedis
	suite.mockRedis = mock
	suite.e = echo.New()
	suite.endpoint = "/flights"
	suite.flight = &Flight{}
}

func (suite *FlightHandlerTestSuite) TestGetFlightsList_GetFromRedis_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	flights := []services.FlightDetails{
		{
			Number:  "FL001",
			Airline: "AirlineA",
		},
	}

	jsonData, _ := json.Marshal(flights)
	monkey.Patch(services.GetFlightsFromRedis, func(_ *redis.Client, _ context.Context, _ string) (string, error) {
		return string(jsonData), nil
	})
	defer monkey.Unpatch(services.GetFlightsFromRedis)
	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27"
	req, res := suite.CallHandler(queryParams)
	ctx := &MockContextWithTrueValidation{
		Context: suite.e.NewContext(req, res),
	}
	err := suite.flight.List(ctx)

	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *FlightHandlerTestSuite) TestGetFlightsList_Validation_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusUnprocessableEntity

	queryParams := ""
	req, res := suite.CallHandler(queryParams)
	ctx := &MockContextWithErrorValidation{
		Context: suite.e.NewContext(req, res),
	}
	_ = suite.flight.List(ctx)

	require.Equal(expectedStatusCode, res.Code)
}

func (suite *FlightHandlerTestSuite) TestGetFlightsList_GetFromWebService_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	flights := []services.FlightDetails{
		{
			Number:  "FL001",
			Airline: "AirlineA",
		},
	}

	monkey.Patch(services.GetFlightsFromRedis, func(_ *redis.Client, _ context.Context, _ string) (string, error) {
		return "", redis.Nil
	})

	defer monkey.Unpatch(services.GetFlightsFromRedis)

	monkey.Patch(services.GetFlightsListFromApi, func(_ *redis.Client, _ *config.FlightService, _ string, _ context.Context, _ []services.FlightDetails, _, _, _ string) ([]services.FlightDetails, error) {
		return flights, nil
	})
	defer monkey.Unpatch(services.GetFlightsListFromApi)
	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27"
	req, res := suite.CallHandler(queryParams)
	ctx := &MockContextWithTrueValidation{
		Context: suite.e.NewContext(req, res),
	}
	err := suite.flight.List(ctx)

	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *FlightHandlerTestSuite) TestGetFlightsList_FilterByAirline() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	flights := []services.FlightDetails{}

	monkey.Patch(services.GetFlightsFromRedis, func(_ *redis.Client, _ context.Context, _ string) (string, error) {
		return "", redis.Nil
	})

	defer monkey.Unpatch(services.GetFlightsFromRedis)

	monkey.Patch(services.GetFlightsListFromApi, func(_ *redis.Client, _ *config.FlightService, _ string, _ context.Context, _ []services.FlightDetails, _, _, _ string) ([]services.FlightDetails, error) {
		return flights, nil
	})
	defer monkey.Unpatch(services.GetFlightsListFromApi)

	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&AL=AirlineA"
	req, res := suite.CallHandler(queryParams)
	ctx := &MockContextWithTrueValidation{
		Context: suite.e.NewContext(req, res),
	}
	err := suite.flight.List(ctx)

	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *FlightHandlerTestSuite) TestGetFlightsList_FilterByAirplane() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	flights := []services.FlightDetails{}

	monkey.Patch(services.GetFlightsFromRedis, func(_ *redis.Client, _ context.Context, _ string) (string, error) {
		return "", redis.Nil
	})

	defer monkey.Unpatch(services.GetFlightsFromRedis)

	monkey.Patch(services.GetFlightsListFromApi, func(_ *redis.Client, _ *config.FlightService, _ string, _ context.Context, _ []services.FlightDetails, _, _, _ string) ([]services.FlightDetails, error) {
		return flights, nil
	})
	defer monkey.Unpatch(services.GetFlightsListFromApi)

	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&AP=Airbus428"
	req, res := suite.CallHandler(queryParams)
	ctx := &MockContextWithTrueValidation{
		Context: suite.e.NewContext(req, res),
	}
	err := suite.flight.List(ctx)

	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *FlightHandlerTestSuite) TestGetFlightsList_FilterByHour() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	flights := []services.FlightDetails{}

	monkey.Patch(services.GetFlightsFromRedis, func(_ *redis.Client, _ context.Context, _ string) (string, error) {
		return "", redis.Nil
	})

	defer monkey.Unpatch(services.GetFlightsFromRedis)

	monkey.Patch(services.GetFlightsListFromApi, func(_ *redis.Client, _ *config.FlightService, _ string, _ context.Context, _ []services.FlightDetails, _, _, _ string) ([]services.FlightDetails, error) {
		return flights, nil
	})
	defer monkey.Unpatch(services.GetFlightsListFromApi)

	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&HO=10"
	req, res := suite.CallHandler(queryParams)
	ctx := &MockContextWithTrueValidation{
		Context: suite.e.NewContext(req, res),
	}
	err := suite.flight.List(ctx)

	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *FlightHandlerTestSuite) TestGetFlightsList_FilterByCapacity() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	flights := []services.FlightDetails{}

	monkey.Patch(services.GetFlightsFromRedis, func(_ *redis.Client, _ context.Context, _ string) (string, error) {
		return "", redis.Nil
	})

	defer monkey.Unpatch(services.GetFlightsFromRedis)

	monkey.Patch(services.GetFlightsListFromApi, func(_ *redis.Client, _ *config.FlightService, _ string, _ context.Context, _ []services.FlightDetails, _, _, _ string) ([]services.FlightDetails, error) {
		return flights, nil
	})
	defer monkey.Unpatch(services.GetFlightsListFromApi)

	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&EC=true"
	req, res := suite.CallHandler(queryParams)
	ctx := &MockContextWithTrueValidation{
		Context: suite.e.NewContext(req, res),
	}
	err := suite.flight.List(ctx)

	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *FlightHandlerTestSuite) TestGetFlightsList_SortByPrice() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	flights := []services.FlightDetails{}

	monkey.Patch(services.GetFlightsFromRedis, func(_ *redis.Client, _ context.Context, _ string) (string, error) {
		return "", redis.Nil
	})

	defer monkey.Unpatch(services.GetFlightsFromRedis)

	monkey.Patch(services.GetFlightsListFromApi, func(_ *redis.Client, _ *config.FlightService, _ string, _ context.Context, _ []services.FlightDetails, _, _, _ string) ([]services.FlightDetails, error) {
		return flights, nil
	})
	defer monkey.Unpatch(services.GetFlightsListFromApi)

	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&order_by=price"
	req, res := suite.CallHandler(queryParams)
	ctx := &MockContextWithTrueValidation{
		Context: suite.e.NewContext(req, res),
	}
	err := suite.flight.List(ctx)

	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *FlightHandlerTestSuite) TestGetFlightsList_SortByTime() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	flights := []services.FlightDetails{}

	monkey.Patch(services.GetFlightsFromRedis, func(_ *redis.Client, _ context.Context, _ string) (string, error) {
		return "", redis.Nil
	})

	defer monkey.Unpatch(services.GetFlightsFromRedis)

	monkey.Patch(services.GetFlightsListFromApi, func(_ *redis.Client, _ *config.FlightService, _ string, _ context.Context, _ []services.FlightDetails, _, _, _ string) ([]services.FlightDetails, error) {
		return flights, nil
	})
	defer monkey.Unpatch(services.GetFlightsListFromApi)

	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&order_by=time"
	req, res := suite.CallHandler(queryParams)
	ctx := &MockContextWithTrueValidation{
		Context: suite.e.NewContext(req, res),
	}
	err := suite.flight.List(ctx)

	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *FlightHandlerTestSuite) TestGetFlightsList_SortByDuration() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	flights := []services.FlightDetails{}

	monkey.Patch(services.GetFlightsFromRedis, func(_ *redis.Client, _ context.Context, _ string) (string, error) {
		return "", redis.Nil
	})

	defer monkey.Unpatch(services.GetFlightsFromRedis)

	monkey.Patch(services.GetFlightsListFromApi, func(_ *redis.Client, _ *config.FlightService, _ string, _ context.Context, _ []services.FlightDetails, _, _, _ string) ([]services.FlightDetails, error) {
		return flights, nil
	})
	defer monkey.Unpatch(services.GetFlightsListFromApi)

	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&order_by=duration"
	req, res := suite.CallHandler(queryParams)
	ctx := &MockContextWithTrueValidation{
		Context: suite.e.NewContext(req, res),
	}
	err := suite.flight.List(ctx)

	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func TestFlightService(t *testing.T) {
	suite.Run(t, new(FlightHandlerTestSuite))
}
