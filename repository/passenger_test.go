package repository

import (
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PassengerTestSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock
	dbMock  *gorm.DB
	UserID  int
}

func (suite *PassengerTestSuite) SetupSuite() {
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
	suite.UserID = 3
}

func (suite *PassengerTestSuite) TestPassenger_CreatePassenger_Success() {
	require := suite.Require()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(
		regexp.QuoteMeta(`
		  INSERT INTO "passengers" ("created_at","updated_at","deleted_at","user_id","national_code","first_name","last_name","gender")
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 `)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	suite.sqlMock.ExpectCommit()
	_, err := CreatePassenger(suite.dbMock, suite.UserID, "0123456789", "fname", "lname", "f")
	require.NoError(err)
}

func (suite *PassengerTestSuite) TestPassenger_CreatePassenger_Failure() {
	require := suite.Require()
	expectedError := "internal error"

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(
		regexp.QuoteMeta(`
		  INSERT INTO "passengers" ("created_at","updated_at","deleted_at","user_id","national_code","first_name","last_name","gender")
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 `)).
		WillReturnError(errors.New("internal error"))
	suite.sqlMock.ExpectRollback()
	res, err := CreatePassenger(suite.dbMock, suite.UserID, "0123456789", "fname", "lname", "f")
	require.Equal(expectedError, string(err.Error()))
	require.Empty(res)
}

func (suite *PassengerTestSuite) TestGetPassenger_Success() {
	require := suite.Require()

	mockPassenger := suite.sqlMock.NewRows(
		[]string{
			"national_code", "first_name", "last_name", "gender",
		}).
		AddRow("1000011111", "name", "lname", "f")
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "passengers" WHERE user_id = (.+)`).
		WillReturnRows(mockPassenger)

	_, err := GetPassengersByUserID(suite.dbMock, suite.UserID)
	require.NoError(err)
}

func (suite *PassengerTestSuite) TestGetPassenger_Failure() {
	require := suite.Require()

	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "passengers" WHERE user_id = (.+)`).
		WillReturnError(errors.New("internal error"))

	_, err := GetPassengersByUserID(suite.dbMock, suite.UserID)
	require.Equal(err.Error(), "internal error")
}

func TestPassenger(t *testing.T) {
	suite.Run(t, new(PassengerTestSuite))
}
