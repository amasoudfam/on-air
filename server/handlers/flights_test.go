package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"on-air/server/services"
	"on-air/utils"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redismock/v9"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type FlightHandlerTestSuite struct {
	suite.Suite
	mockRedis redismock.ClientMock
	redis     *redis.Client
	e         *echo.Echo
	endpoint  string
	flight    *Flight
}

func (suite *FlightHandlerTestSuite) CallHandler(queryString string) (*httptest.ResponseRecorder, error) {
	url := suite.endpoint + queryString
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := suite.e.NewContext(req, res)
	err := suite.flight.List(ctx)
	return res, err
}

func (suite *FlightHandlerTestSuite) SetupSuite() {
	mockRedis, mock := redismock.NewClientMock()
	suite.redis = mockRedis
	suite.mockRedis = mock
	suite.e = echo.New()
	suite.endpoint = "/flights"
	validator := validator.New()
	validator.RegisterValidation("CustomTimeValidator", utils.CustomTimeValidator)
	suite.e.Validator = &utils.CustomValidator{Validator: validator}
	suite.flight = &Flight{
		Redis:         suite.redis,
		APIMockClient: &services.APIMockClient{},
	}
}

func (suite *FlightHandlerTestSuite) TestGetFlightsList_WithCache_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	expectedMsg := "{\"flights\":[{\"number\":\"FL001\",\"airplane\":\"\",\"airline\":\"AirlineA\",\"price\":0,\"origin\":{\"id\":0,\"name\":\"\"},\"destination\":{\"id\":0,\"name\":\"\"},\"capacity\":0,\"empty_capacity\":0,\"started_at\":\"0001-01-01T00:00:00Z\",\"finished_at\":\"0001-01-01T00:00:00Z\"}]}\n"
	flights := []services.FlightResponse{
		{
			Number:  "FL001",
			Airline: "AirlineA",
		},
	}

	queryString := "?org=Shiraz&dest=Esfahan&date=2023-06-27"
	jsonData, _ := json.Marshal(flights)
	suite.mockRedis.ExpectGet("flights_Shiraz_Esfahan_2023-06-27").SetVal(string(jsonData))

	res, err := suite.CallHandler(queryString)
	require.NoError(err)
	require.Equal(expectedMsg, res.Body.String())
	require.Equal(expectedStatusCode, res.Code)

	err = suite.mockRedis.ExpectationsWereMet()
	require.NoError(err)
}

//func (suite *FlightHandlerTestSuite) TestGetFlightsList_WithoutCache_Success() {
//	require := suite.Require()
//	expectedStatusCode := http.StatusOK
//	flights := []services.FlightResponse{
//		{
//			Number:  "FL001",
//			Airline: "AirlineA",
//		},
//	}
//
//	// Set up the expectation on the Redis mock object
//	suite.mockRedis.ExpectGet("flights_Shiraz_Esfahan_2023-06-27").RedisNil()
//
//	monkey.Patch(services.GetFlights, func(_ *redis.Client, _ *config.APIMock, _ string, _ context.Context, _ []services.FlightResponse, _, _, _ string) ([]services.FlightResponse, error) {
//		return flights, nil
//	})
//	defer monkey.Unpatch(services.GetFlights)
//
//	jsonData, _ := json.Marshal(flights)
//	suite.mockRedis.ExpectSet("flights_Shiraz_Esfahan_2023-06-27", jsonData, time.Minute*10).SetVal(string(jsonData))
//
//	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27"
//	res, err := suite.CallHandler(queryParams)
//	require.NoError(err)
//	body, _ := io.ReadAll(res.Body)
//	var response ListResponse
//	err = json.Unmarshal(body, &response)
//	require.NoError(err)
//	require.Equal(flights, response.Flights)
//
//	require.NoError(err)
//	require.Equal(expectedStatusCode, res.Code)
//}

func (suite *FlightHandlerTestSuite) TestGetFlightsList_ParseReq_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest

	queryString := "?org=Shiraz&dest=Esfahan&&date2023-06-27&empty_capacity=test"
	res, err := suite.CallHandler(queryString)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *FlightHandlerTestSuite) TestGetFlightsList_Validation_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusUnprocessableEntity
	queryParams := ""
	res, err := suite.CallHandler(queryParams)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

//func (suite *FlightHandlerTestSuite) TestGetFlightsList_GetFromWebService_Failure() {
//	require := suite.Require()
//	expectedStatusCode := http.StatusInternalServerError
//
//	// Set up the expectation on the Redis mock object
//	suite.mockRedis.ExpectGet("flights_Shiraz_Esfahan_2023-06-27").RedisNil()
//
//	monkey.Patch(services.GetFlights, func(_ *redis.Client, _ *config.APIMock, _ string, _ context.Context, _ []services.FlightResponse, _, _, _ string) ([]services.FlightResponse, error) {
//		return nil, errors.New("error")
//	})
//	defer monkey.Unpatch(services.GetFlights)
//
//	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27"
//	res, err := suite.CallHandler(queryParams)
//	require.NoError(err)
//	require.Equal(expectedStatusCode, res.Code)
//}
//
//func (suite *FlightHandlerTestSuite) TestGetFlightsList_FilterByAirline() {
//	require := suite.Require()
//	expectedStatusCode := http.StatusOK
//	flights := []services.FlightResponse{}
//
//	suite.mockRedis.ExpectGet("flights_Shiraz_Esfahan_2023-06-27").RedisNil()
//	monkey.Patch(services.GetFlights, func(_ *redis.Client, _ *config.APIMock, _ string, _ context.Context, _ []services.FlightResponse, _, _, _ string) ([]services.FlightResponse, error) {
//		return flights, nil
//	})
//	defer monkey.Unpatch(services.GetFlights)
//
//	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&airline=AirlineA"
//	res, err := suite.CallHandler(queryParams)
//	require.NoError(err)
//	require.Equal(expectedStatusCode, res.Code)
//}
//
//func (suite *FlightHandlerTestSuite) TestGetFlightsList_FilterByAirplane() {
//	require := suite.Require()
//	expectedStatusCode := http.StatusOK
//	flights := []services.FlightResponse{}
//
//	suite.mockRedis.ExpectGet("flights_Shiraz_Esfahan_2023-06-27").RedisNil()
//	monkey.Patch(services.GetFlights, func(_ *redis.Client, _ *config.APIMock, _ string, _ context.Context, _ []services.FlightResponse, _, _, _ string) ([]services.FlightResponse, error) {
//		return flights, nil
//	})
//	defer monkey.Unpatch(services.GetFlights)
//
//	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&airplane=Airbus428"
//	res, err := suite.CallHandler(queryParams)
//	require.NoError(err)
//	require.Equal(expectedStatusCode, res.Code)
//}
//
//func (suite *FlightHandlerTestSuite) TestGetFlightsList_FilterTime() {
//	require := suite.Require()
//	expectedStatusCode := http.StatusOK
//	flights := []services.FlightResponse{}
//
//	suite.mockRedis.ExpectGet("flights_Shiraz_Esfahan_2023-06-27").RedisNil()
//	monkey.Patch(services.GetFlights, func(_ *redis.Client, _ *config.APIMock, _ string, _ context.Context, _ []services.FlightResponse, _, _, _ string) ([]services.FlightResponse, error) {
//		return flights, nil
//	})
//	defer monkey.Unpatch(services.GetFlights)
//
//	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&start_time:10:30&end_time:11:30"
//	res, err := suite.CallHandler(queryParams)
//	require.NoError(err)
//	require.Equal(expectedStatusCode, res.Code)
//}
//
//func (suite *FlightHandlerTestSuite) TestGetFlightsList_FilterByCapacity() {
//	require := suite.Require()
//	expectedStatusCode := http.StatusOK
//	flights := []services.FlightResponse{}
//
//	suite.mockRedis.ExpectGet("flights_Shiraz_Esfahan_2023-06-27").RedisNil()
//	monkey.Patch(services.GetFlights, func(_ *redis.Client, _ *config.APIMock, _ string, _ context.Context, _ []services.FlightResponse, _, _, _ string) ([]services.FlightResponse, error) {
//		return flights, nil
//	})
//	defer monkey.Unpatch(services.GetFlights)
//
//	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&empty_capacity=true"
//	res, err := suite.CallHandler(queryParams)
//	require.NoError(err)
//	require.Equal(expectedStatusCode, res.Code)
//}
//
//func (suite *FlightHandlerTestSuite) TestGetFlightsList_SortByPrice() {
//	require := suite.Require()
//	expectedStatusCode := http.StatusOK
//	flights := []services.FlightResponse{}
//
//	suite.mockRedis.ExpectGet("flights_Shiraz_Esfahan_2023-06-27").RedisNil()
//	monkey.Patch(services.GetFlights, func(_ *redis.Client, _ *config.APIMock, _ string, _ context.Context, _ []services.FlightResponse, _, _, _ string) ([]services.FlightResponse, error) {
//		return flights, nil
//	})
//	defer monkey.Unpatch(services.GetFlights)
//
//	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&order_by=price"
//	res, err := suite.CallHandler(queryParams)
//	require.NoError(err)
//	require.Equal(expectedStatusCode, res.Code)
//}
//
//func (suite *FlightHandlerTestSuite) TestGetFlightsList_SortByTime() {
//	require := suite.Require()
//	expectedStatusCode := http.StatusOK
//	flights := []services.FlightResponse{}
//
//	suite.mockRedis.ExpectGet("flights_Shiraz_Esfahan_2023-06-27").RedisNil()
//	monkey.Patch(services.GetFlights, func(_ *redis.Client, _ *config.APIMock, _ string, _ context.Context, _ []services.FlightResponse, _, _, _ string) ([]services.FlightResponse, error) {
//		return flights, nil
//	})
//	defer monkey.Unpatch(services.GetFlights)
//
//	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&order_by=time"
//	res, err := suite.CallHandler(queryParams)
//	require.NoError(err)
//	require.Equal(expectedStatusCode, res.Code)
//}
//
//func (suite *FlightHandlerTestSuite) TestGetFlightsList_SortByDuration() {
//	require := suite.Require()
//	expectedStatusCode := http.StatusOK
//	flights := []services.FlightResponse{}
//
//	suite.mockRedis.ExpectGet("flights_Shiraz_Esfahan_2023-06-27").RedisNil()
//	monkey.Patch(services.GetFlights, func(_ *redis.Client, _ *config.APIMock, _ string, _ context.Context, _ []services.FlightResponse, _, _, _ string) ([]services.FlightResponse, error) {
//		return flights, nil
//	})
//	defer monkey.Unpatch(services.GetFlights)
//
//	queryParams := "?org=Shiraz&dest=Esfahan&date=2023-06-27&order_by=duration"
//	res, err := suite.CallHandler(queryParams)
//	require.NoError(err)
//	require.Equal(expectedStatusCode, res.Code)
//}

func (suite *FlightHandlerTestSuite) TestFilterByAirline() {
	require := suite.Require()
	flights := []services.FlightResponse{
		{Number: "FL001", Airline: "AirlineA"},
		{Number: "FL002", Airline: "AirlineB"},
		{Number: "FL003", Airline: "AirlineA"},
		{Number: "FL004", Airline: "AirlineC"},
	}
	airline := "AirlineA"
	filteredFlights := filterByAirline(flights, airline)

	expectedFlights := []services.FlightResponse{
		{Number: "FL001", Airline: "AirlineA"},
		{Number: "FL003", Airline: "AirlineA"},
	}
	require.Len(filteredFlights, 2)
	require.Equal(expectedFlights, filteredFlights)
}

func (suite *FlightHandlerTestSuite) TestFilterByAirplane() {
	require := suite.Require()
	flights := []services.FlightResponse{
		{Number: "FL001", Airline: "AirlineA", Airplane: "bb"},
		{Number: "FL002", Airline: "AirlineB", Airplane: "hh"},
		{Number: "FL003", Airline: "AirlineA", Airplane: "bb"},
		{Number: "FL004", Airline: "AirlineC", Airplane: "cc"},
	}
	airplane := "hh"
	filteredFlights := filterByAirplane(flights, airplane)

	expectedFlights := []services.FlightResponse{
		{Number: "FL002", Airline: "AirlineB", Airplane: "hh"},
	}
	require.Len(filteredFlights, 1)
	require.Equal(expectedFlights, filteredFlights)
}

func (suite *FlightHandlerTestSuite) TestFilterByCapacity() {
	require := suite.Require()
	flights := []services.FlightResponse{
		{Number: "FL001", Airline: "AirlineA", EmptyCapacity: 0},
		{Number: "FL002", Airline: "AirlineB", EmptyCapacity: 0},
		{Number: "FL003", Airline: "AirlineA", EmptyCapacity: 12},
		{Number: "FL004", Airline: "AirlineC", EmptyCapacity: 2},
	}
	filteredFlights := filterByCapacity(flights)

	expectedFlights := []services.FlightResponse{
		{Number: "FL003", Airline: "AirlineA", EmptyCapacity: 12},
		{Number: "FL004", Airline: "AirlineC", EmptyCapacity: 2},
	}
	require.Len(filteredFlights, 2)
	require.Equal(expectedFlights, filteredFlights)
}

func (suite *FlightHandlerTestSuite) TestFilterByHour() {
	require := suite.Require()
	flights := []services.FlightResponse{
		{Number: "FL001", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC)},
		{Number: "FL002", Airline: "AirlineB", StartedAt: time.Date(2023, 6, 1, 11, 0, 0, 0, time.UTC)},
		{Number: "FL003", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 30, 0, 0, time.UTC)},
		{Number: "FL004", Airline: "AirlineC", StartedAt: time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)},
	}

	start_time := "10:00"
	end_time := "11:30"
	filteredFlights := filterByTime(flights, start_time, end_time)

	expectedFlights := []services.FlightResponse{
		{Number: "FL001", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC)},
		{Number: "FL002", Airline: "AirlineB", StartedAt: time.Date(2023, 6, 1, 11, 0, 0, 0, time.UTC)},
		{Number: "FL003", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 30, 0, 0, time.UTC)},
	}

	require.Len(filteredFlights, 3)
	require.Equal(expectedFlights, filteredFlights)
}

func (suite *FlightHandlerTestSuite) TestSortByPrice_desc() {
	require := suite.Require()
	flights := []services.FlightResponse{
		{Number: "FL001", Airline: "AirlineA", Price: 1800000},
		{Number: "FL002", Airline: "AirlineB", Price: 1300000},
		{Number: "FL003", Airline: "AirlineA", Price: 1900000},
		{Number: "FL004", Airline: "AirlineC", Price: 1100000},
	}
	filteredFlights := sortByPrice(flights, "desc")

	expectedFlights := []services.FlightResponse{
		{Number: "FL003", Airline: "AirlineA", Price: 1900000},
		{Number: "FL001", Airline: "AirlineA", Price: 1800000},
		{Number: "FL002", Airline: "AirlineB", Price: 1300000},
		{Number: "FL004", Airline: "AirlineC", Price: 1100000},
	}

	require.Equal(expectedFlights, filteredFlights)
}

func (suite *FlightHandlerTestSuite) TestSortByPrice_acs() {
	require := suite.Require()
	flights := []services.FlightResponse{
		{Number: "FL001", Airline: "AirlineA", Price: 1800000},
		{Number: "FL002", Airline: "AirlineB", Price: 1300000},
		{Number: "FL003", Airline: "AirlineA", Price: 1900000},
		{Number: "FL004", Airline: "AirlineC", Price: 1100000},
	}
	sortFlights := sortByPrice(flights, "asc")

	expectedFlights := []services.FlightResponse{
		{Number: "FL004", Airline: "AirlineC", Price: 1100000},
		{Number: "FL002", Airline: "AirlineB", Price: 1300000},
		{Number: "FL001", Airline: "AirlineA", Price: 1800000},
		{Number: "FL003", Airline: "AirlineA", Price: 1900000},
	}

	require.Equal(expectedFlights, sortFlights)
}

func (suite *FlightHandlerTestSuite) TestSortByTime_asc() {
	require := suite.Require()
	flights := []services.FlightResponse{
		{Number: "FL001", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC)},
		{Number: "FL003", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 2, 10, 30, 0, 0, time.UTC)},
		{Number: "FL002", Airline: "AirlineB", StartedAt: time.Date(2023, 6, 1, 11, 0, 0, 0, time.UTC)},
		{Number: "FL004", Airline: "AirlineC", StartedAt: time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)},
	}
	sortFlights := sortByTime(flights, "asc")

	expectedFlights := []services.FlightResponse{
		{Number: "FL001", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC)},
		{Number: "FL002", Airline: "AirlineB", StartedAt: time.Date(2023, 6, 1, 11, 0, 0, 0, time.UTC)},
		{Number: "FL003", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 2, 10, 30, 0, 0, time.UTC)},
		{Number: "FL004", Airline: "AirlineC", StartedAt: time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)},
	}

	require.Equal(expectedFlights, sortFlights)
}

func (suite *FlightHandlerTestSuite) TestSortByTime_desc() {
	require := suite.Require()
	flights := []services.FlightResponse{
		{Number: "FL001", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC)},
		{Number: "FL003", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 2, 10, 30, 0, 0, time.UTC)},
		{Number: "FL002", Airline: "AirlineB", StartedAt: time.Date(2023, 6, 1, 11, 0, 0, 0, time.UTC)},
		{Number: "FL004", Airline: "AirlineC", StartedAt: time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)},
	}
	sortFlights := sortByTime(flights, "desc")

	expectedFlights := []services.FlightResponse{
		{Number: "FL004", Airline: "AirlineC", StartedAt: time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)},
		{Number: "FL003", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 2, 10, 30, 0, 0, time.UTC)},
		{Number: "FL002", Airline: "AirlineB", StartedAt: time.Date(2023, 6, 1, 11, 0, 0, 0, time.UTC)},
		{Number: "FL001", Airline: "AirlineA", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC)},
	}

	require.Equal(expectedFlights, sortFlights)
}

func (suite *FlightHandlerTestSuite) TestSortByDuration_asc() {
	require := suite.Require()
	flights := []services.FlightResponse{
		{Number: "FL001", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 11, 23, 0, 0, time.UTC)},
		{Number: "FL003", StartedAt: time.Date(2023, 6, 2, 8, 30, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 11, 50, 0, 0, time.UTC)},
		{Number: "FL002", StartedAt: time.Date(2023, 6, 1, 06, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 10, 20, 0, 0, time.UTC)},
		{Number: "FL004", StartedAt: time.Date(2023, 6, 2, 05, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 06, 0, 0, 0, time.UTC)},
	}
	sortFlights := sortByDuration(flights, "asc")

	expectedFlights := []services.FlightResponse{
		{Number: "FL004", StartedAt: time.Date(2023, 6, 2, 05, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 06, 0, 0, 0, time.UTC)},
		{Number: "FL001", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 11, 23, 0, 0, time.UTC)},
		{Number: "FL003", StartedAt: time.Date(2023, 6, 2, 8, 30, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 11, 50, 0, 0, time.UTC)},
		{Number: "FL002", StartedAt: time.Date(2023, 6, 1, 06, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 10, 20, 0, 0, time.UTC)},
	}

	require.Equal(expectedFlights, sortFlights)
}

func (suite *FlightHandlerTestSuite) TestSortByDuration_desc() {
	require := suite.Require()
	flights := []services.FlightResponse{
		{Number: "FL001", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 11, 23, 0, 0, time.UTC)},
		{Number: "FL003", StartedAt: time.Date(2023, 6, 2, 8, 30, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 11, 50, 0, 0, time.UTC)},
		{Number: "FL002", StartedAt: time.Date(2023, 6, 1, 06, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 10, 20, 0, 0, time.UTC)},
		{Number: "FL004", StartedAt: time.Date(2023, 6, 2, 05, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 06, 0, 0, 0, time.UTC)},
	}
	sortFlights := sortByDuration(flights, "desc")

	expectedFlights := []services.FlightResponse{
		{Number: "FL002", StartedAt: time.Date(2023, 6, 1, 06, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 10, 20, 0, 0, time.UTC)},
		{Number: "FL003", StartedAt: time.Date(2023, 6, 2, 8, 30, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 11, 50, 0, 0, time.UTC)},
		{Number: "FL001", StartedAt: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 1, 11, 23, 0, 0, time.UTC)},
		{Number: "FL004", StartedAt: time.Date(2023, 6, 2, 05, 0, 0, 0, time.UTC), FinishedAt: time.Date(2023, 6, 2, 06, 0, 0, 0, time.UTC)},
	}

	require.Equal(expectedFlights, sortFlights)
}

func (suite *FlightHandlerTestSuite) TestCustomTimeValidator() {
	type testStruct struct {
		StartTime string `validate:"required,customTimeValidator"`
		EndTime   string `validate:"required,customTimeValidator,eqfield=StartTime"`
	}

	tests := []struct {
		Name      string
		StartTime string
		EndTime   string
		Expected  bool
	}{
		{"ValidTimes", "21:30", "21:30", true},
		{"InvalidTimes", "21:30", "22:00", false},
	}

	validator := validator.New()
	_ = validator.RegisterValidation("customTimeValidator", utils.CustomTimeValidator)

	for _, test := range tests {
		suite.Run(test.Name, func() {
			testData := testStruct{
				StartTime: test.StartTime,
				EndTime:   test.EndTime,
			}

			err := validator.Struct(testData)
			valid := err == nil
			suite.Equal(test.Expected, valid)
		})
	}
}

func TestFlightService(t *testing.T) {
	suite.Run(t, new(FlightHandlerTestSuite))
}
