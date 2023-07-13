package repository

import (
	"errors"
	"log"
	"net/http"
	"on-air/models"
	"on-air/server/services"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eapache/go-resiliency/breaker"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CityTestSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock
	dbMock  *gorm.DB
	city    City
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
	suite.city = City{
		APIMockClient: &services.APIMockClient{
			Client:  &http.Client{},
			Breaker: &breaker.Breaker{},
			BaseURL: "http://example.com",
			Timeout: time.Second,
		},
		DB:         suite.dbMock,
		SyncPeriod: time.Duration(2 * time.Second),
	}
}

func (suite *CityTestSuite) Test_GetCityByName_Success() {
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
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "cities" WHERE name = \$1 (.+)`).
		WithArgs("Shiraz").
		WillReturnRows(mockCity)

	city, err := FindCityByName(suite.dbMock, "Shiraz")

	require.NoError(err)
	require.Equal(expectedCity, *city)
}

func (suite *CityTestSuite) Test_GetCityByName_Failure() {
	require := suite.Require()
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "cities" WHERE name = \$1 (.+)`).
		WithArgs("Shiraz").
		WillReturnError(errors.New("internal error"))

	_, err := FindCityByName(suite.dbMock, "Shiraz")

	require.Equal(err.Error(), "internal error")
}

func (suite *CityTestSuite) Test_SyncCities_Success() {
	require := suite.Require()

	getCitiesPatchCalledCount := 0
	getCitiesPatch := monkey.PatchInstanceMethod(
		reflect.TypeOf(suite.city.APIMockClient),
		"GetCities",
		func(_ *services.APIMockClient) ([]string, error) {
			getCitiesPatchCalledCount++
			return []string{"Tehran", "Tabriz", "Shiraz", "Kish", "Esfahan", "Qeshm", "Mashhad"}, nil
		},
	)
	defer getCitiesPatch.Unpatch()

	var c City
	storeCitiesCalled := false
	storeCitiesPatch := monkey.PatchInstanceMethod(reflect.TypeOf(&c), "StoreCities",
		func(c *City, cities []string) error {
			storeCitiesCalled = true
			return nil
		})
	defer storeCitiesPatch.Unpatch()

	go suite.city.SyncCities()
	time.Sleep(5 * time.Second)
	require.Equal(3, getCitiesPatchCalledCount)
	require.Equal(true, storeCitiesCalled)
}

func TestCityRepository(t *testing.T) {
	suite.Run(t, new(CityTestSuite))
}
