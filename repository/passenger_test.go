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

const UserID = 3

type PassengerTestSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock
	dbMock  *gorm.DB
}

func (suite *PassengerTestSuite) SetupSuite() {
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

func (suite *PassengerTestSuite) TestPassenger_CreatePassenger_Success() {
	require := suite.Require()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("INSERT INTO `passengers`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.sqlMock.ExpectCommit()

	_, err := CreatePassenger(suite.dbMock, 3, "", "", "", "")
	require.NoError(err)
}

func (suite *PassengerTestSuite) TestPassenger_CreatePassenger_Failure() {
	require := suite.Require()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("INSERT INTO `passengers`").
		WillReturnError(errors.New(""))
	suite.sqlMock.ExpectRollback()

	res, _ := CreatePassenger(suite.dbMock, 3, "", "", "", "")
	require.Empty(res)
}

func (suite *PassengerTestSuite) TestGetPassenger_Success() {
	require := suite.Require()

	mockRow := suite.sqlMock.NewRows(
		[]string{
			"nationalcode", "firstname", "lastname", "gender",
		}).AddRow("1000011111", "name", "lname", "f")
	suite.sqlMock.ExpectQuery("SELECT (.+) FROM `passengers` WHERE user_id = ?").WithArgs(UserID).WillReturnRows(mockRow)

	_, err := GetPassengersByUserID(suite.dbMock, UserID)
	require.NoError(err)
}

func (suite *PassengerTestSuite) TestGetPassenger_Failure() {
	require := suite.Require()

	suite.sqlMock.ExpectQuery("SELECT (.+) FROM `passengers` WHERE user_id = ?").WithArgs(UserID).WillReturnError(errors.New(""))

	res, _ := GetPassengersByUserID(suite.dbMock, UserID)
	require.Empty(res)
}

func TestPassenger(t *testing.T) {
	suite.Run(t, new(PassengerTestSuite))
}
