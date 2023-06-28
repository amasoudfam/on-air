package repository

import (
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
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

	suite.dbMock, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}))

	if err != nil {
		log.Fatal(err)
	}

	suite.sqlMock = sqlMock
}

func (suite *UserTestSuite) TestUser_RegisterUser_Success() {
	require := suite.Require()
	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("INSERT INTO `users`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.sqlMock.ExpectCommit()

	_, err := RegisterUser(suite.dbMock, UserEmail, UserPassword)
	require.NoError(err)
}

func (suite *UserTestSuite) TestPassenger_RegisterUser_Failure() {
	require := suite.Require()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("INSERT INTO `users`").
		WillReturnError(errors.New(""))
	suite.sqlMock.ExpectRollback()

	res, _ := RegisterUser(suite.dbMock, UserEmail, UserPassword)
	require.Empty(res)
}

func (suite *UserTestSuite) TestUser_GetUserByEmail_Success() {
	require := suite.Require()

	mockRow := suite.sqlMock.NewRows(
		[]string{
			"id", "email", "password",
		}).AddRow(UserID, UserEmail, UserPassword)
	suite.sqlMock.ExpectQuery("SELECT (.+) FROM `users` WHERE email = ?").WithArgs(UserEmail).WillReturnRows(mockRow)

	_, err := GetUserByEmail(suite.dbMock, UserEmail)
	require.NoError(err)
}

func (suite *UserTestSuite) TestUser_GetUserByEmail_Failure() {
	require := suite.Require()

	suite.sqlMock.ExpectQuery("SELECT (.+) FROM `users` WHERE email = ?").WithArgs(UserEmail).WillReturnError(errors.New(""))

	res, _ := GetUserByEmail(suite.dbMock, UserEmail)
	require.Empty(res)
}

func TestUser(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
