package repository

import (
	"errors"
	"log"
	"on-air/models"
	"regexp"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type FlightTestSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock
	dbMock  *gorm.DB
}

func (suite *FlightTestSuite) SetupSuite() {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	suite.dbMock, err = gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}))

	if err != nil {
		log.Fatal(err)
	}

	suite.sqlMock = sqlMock
}

func (suite *FlightTestSuite) TestTickets_AddFlight_Success() {
	require := suite.Require()

	city := getCity()
	monkey.Patch(FindCityByName, func(db *gorm.DB, Name string) (*models.City, error) {
		return city(), nil
	})
	defer monkey.Unpatch(FindCityByName)

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "flights"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	suite.sqlMock.ExpectCommit()

	_, err := AddFlight(suite.dbMock, "F101", "Shiraz", "Tehran", "Homa", "Airbus_360", datatypes.JSON([]byte(`{"type": "fine", "amount": 100}`)), time.Now(), time.Now())
	require.NoError(err)
}

func (suite *FlightTestSuite) TestTickets_AddFlightWhenGetFromCity_Failure() {
	require := suite.Require()

	monkey.Patch(FindCityByName, func(db *gorm.DB, Name string) (*models.City, error) {
		return nil, errors.New("Internal database error")
	})
	defer monkey.Unpatch(FindCityByName)

	_, err := AddFlight(suite.dbMock, "F101", "Shiraz", "Tehran", "Homa", "Airbus_360", datatypes.JSON([]byte(`{"type": "fine", "amount": 100}`)), time.Now(), time.Now())
	require.Equal(err.Error(), "Internal database error")
}

func (suite *FlightTestSuite) TestTickets_AddFlightWhenCreateFlight_Failure() {
	require := suite.Require()
	city := getCity()
	monkey.Patch(FindCityByName, func(db *gorm.DB, Name string) (*models.City, error) {
		return city(), nil
	})
	defer monkey.Unpatch(FindCityByName)

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "flights"`)).
		WillReturnError(errors.New("Internal database error"))
	suite.sqlMock.ExpectRollback()

	_, err := AddFlight(suite.dbMock, "F101", "Shiraz", "Tehran", "Homa", "Airbus_360", datatypes.JSON([]byte(`{"type": "fine", "amount": 100}`)), time.Now(), time.Now())
	require.Equal(err.Error(), "Internal database error")
}

func getCity() func() *models.City {
	interval := 1
	cities := []string{"Shiraz", "Tehran"}
	return func() *models.City {
		city := models.City{
			Name:      cities[interval-1],
			CountryID: uint(1),
			Country: models.Country{
				Name: "Iran",
			},
		}
		city.ID = uint(interval)
		city.Country.ID = uint(1)
		interval++
		return &city
	}
}

func (suite *FlightTestSuite) TestTickets_FindFlight_Success() {
	require := suite.Require()

	expectedFlight := models.Flight{
		Number:     "F101",
		FromCityID: uint(1),
		ToCityID:   uint(2),
		Airplane:   "Airbus_360",
		Airline:    "Homa",
	}

	expectedFlight.ID = uint(1)
	mockFlight := suite.sqlMock.NewRows(
		[]string{
			"id", "number", "from_city_id", "to_city_id", "airplane", "airline",
		}).
		AddRow(1, "F101", 1, 2, "Airbus_360", "Homa")
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "flights" WHERE Number = \$1 (.+)`).
		WithArgs("F101").
		WillReturnRows(mockFlight)

	flight, err := FindFlight(suite.dbMock, "F101")
	require.NoError(err)
	require.Equal(expectedFlight, *flight)
}

func (suite *FlightTestSuite) TestTickets_FindFlight_Failure() {
	require := suite.Require()
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "flights" WHERE Number = \$1 (.+)`).
		WithArgs("F101").
		WillReturnError(errors.New("internal error"))

	_, err := FindFlight(suite.dbMock, "F101")
	require.Equal(err.Error(), "internal error")
}

func (suite *FlightTestSuite) TestTickets_FindFlightById_Success() {
	require := suite.Require()

	expectedFlight := models.Flight{
		Number:     "F101",
		FromCityID: uint(1),
		ToCityID:   uint(2),
		Airplane:   "Airbus_360",
		Airline:    "Homa",
	}

	expectedFlight.ID = uint(1)
	mockFlight := suite.sqlMock.NewRows(
		[]string{
			"id", "number", "from_city_id", "to_city_id", "airplane", "airline",
		}).
		AddRow(1, "F101", 1, 2, "Airbus_360", "Homa")
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "flights" WHERE ID = \$1 (.+)`).
		WithArgs(1).
		WillReturnRows(mockFlight)

	flight, err := FindFlightById(suite.dbMock, 1)
	require.NoError(err)
	require.Equal(expectedFlight, *flight)
}

func (suite *FlightTestSuite) TestTickets_FindFlightById_Failure() {
	require := suite.Require()
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "flights" WHERE ID = \$1 (.+)`).
		WithArgs(1).
		WillReturnError(errors.New("internal error"))

	_, err := FindFlightById(suite.dbMock, 1)
	require.Equal(err.Error(), "internal error")
}

func TestFlightRepository(t *testing.T) {
	suite.Run(t, new(FlightTestSuite))
}
