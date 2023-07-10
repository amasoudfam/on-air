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

type TicketTestSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock
	dbMock  *gorm.DB
	UserID  int
}

func (suite *TicketTestSuite) SetupSuite() {
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

func (suite *TicketTestSuite) TestTicket_ReserveTicket_Success() {
	require := suite.Require()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "tickets"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	suite.sqlMock.ExpectCommit()
	_, err := ReserveTicket(suite.dbMock, suite.UserID, 1, 10000, []int{1, 2, 3})
	require.NoError(err)
}

func (suite *TicketTestSuite) TestTicket_ReserveTicket_Failure() {
	require := suite.Require()
	expectedError := "internal error"

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "tickets"`)).
		WillReturnError(errors.New("internal error"))
	suite.sqlMock.ExpectRollback()
	res, err := ReserveTicket(suite.dbMock, suite.UserID, 1, 0, []int{1, 2, 3})
	require.Equal(expectedError, string(err.Error()))
	require.Empty(res)
}

func (suite *TicketTestSuite) TestTicket_GetExpiredTickets_Success() {
	require := suite.Require()

	mockPassenger := suite.sqlMock.NewRows(
		[]string{
			"user_id", "unit_price", "count", "flight_id", "status",
		}).
		AddRow("1", "100000", "9", "1", "Expired")
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "tickets" WHERE user_id = (.+)`).
		WillReturnRows(mockPassenger)

	_, err := GetExpiredTickets(suite.dbMock)
	require.NoError(err)
}

func (suite *TicketTestSuite) TestTicket_GetExpiredTickets_Failure() {
	require := suite.Require()

	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "tickets" WHERE user_id = (.+)`).
		WillReturnError(errors.New("internal error"))

	_, err := GetPassengersByUserID(suite.dbMock, suite.UserID)
	require.Equal(err.Error(), "internal error")
}

func TestTicket(t *testing.T) {
	suite.Run(t, new(TicketTestSuite))
}
