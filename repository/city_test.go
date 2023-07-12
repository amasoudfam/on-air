package repository

import (
	"errors"
	"log"
	"on-air/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CityTestSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock
	dbMock  *gorm.DB
}

func (suite *CityTestSuite) SetupSuite() {
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

func (suite *CityTestSuite) TestTickets_GetCityByName_Success() {
	require := suite.Require()

	expectedCity := models.City{
		Name:      "Shiraz",
		CountryID: uint(1),
	}

	expectedCity.ID = uint(1)
	mockCity := suite.sqlMock.NewRows(
		[]string{
			"id", "Name", "country_id",
		}).
		AddRow(1, "Shiraz", 1)
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "cities" WHERE Name = \$1 (.+)`).
		WithArgs("Shiraz").
		WillReturnRows(mockCity)

	city, err := FindCityByName(suite.dbMock, "Shiraz")

	require.NoError(err)
	require.Equal(expectedCity, *city)
}

func (suite *CityTestSuite) TestTickets_GetCityByName_Failure() {
	require := suite.Require()
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "cities" WHERE Name = \$1 (.+)`).
		WithArgs("Shiraz").
		WillReturnError(errors.New("internal error"))

	_, err := FindCityByName(suite.dbMock, "Shiraz")

	require.Equal(err.Error(), "internal error")
}

func TestCityRepository(t *testing.T) {
	suite.Run(t, new(CityTestSuite))
}
