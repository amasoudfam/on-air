package repository

import (
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserTestSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock
	dbMock  *gorm.DB
}

const (
	UserId       = 1
	UserEmail    = "admin@gmail.com"
	UserPassword = "password@789"
)

func (suite *UserTestSuite) SetupSuite() {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	suite.dbMock, err = gorm.Open(postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	}))

	if err != nil {
		log.Fatal(err)
	}

	suite.sqlMock = sqlMock
}

func (suite *UserTestSuite) TestUser_RegisterUser_Success() {
	require := suite.Require()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(`INSERT`).WillReturnRows(sqlmock.NewRows([]string{"id"}))
	suite.sqlMock.ExpectCommit()

	_, err := RegisterUser(suite.dbMock, UserEmail, UserPassword)
	require.NoError(err)
}

func (suite *UserTestSuite) TestUser_RegisterUser_Failure() {
	require := suite.Require()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(`INSERT`).
		WillReturnError(errors.New("user not found"))
	suite.sqlMock.ExpectRollback()

	res, _ := RegisterUser(suite.dbMock, UserEmail, UserPassword)
	require.Empty(res)
}

func (suite *UserTestSuite) TestUser_GetUserByEmail_Success() {
	require := suite.Require()

	mockUser := suite.sqlMock.NewRows(
		[]string{
			"id", "email", "password",
		}).AddRow(UserId, UserEmail, UserPassword)

	suite.sqlMock.ExpectQuery(`SELECT`).WillReturnRows(mockUser)

	_, err := GetUserByEmail(suite.dbMock, UserEmail)
	require.NoError(err)
}

func (suite *UserTestSuite) TestUser_GetUserByEmail_Failure() {
	require := suite.Require()

	suite.sqlMock.ExpectQuery("SELECT").WithArgs(UserEmail).WillReturnError(errors.New(""))

	res, _ := GetUserByEmail(suite.dbMock, UserEmail)
	require.Empty(res)
}

func TestUser(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
