package repository

import (
	"log"
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
	UserID  uint
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
	suite.UserID = uint(1)
}

func (suite *TicketTestSuite) TestTickets_GetTickets_Success() {
	require := suite.Require()

	mockTicketRows := suite.sqlMock.NewRows([]string{"id", "unit_price", "count", "status"}).
		AddRow(1, 100, 2, "complete").
		AddRow(2, 150, 1, "pending")

	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "tickets" WHERE user_id = \$1 AND (.+)`).
		WithArgs(suite.UserID).
		WillReturnRows(mockTicketRows)

	mockPassengerRows := suite.sqlMock.NewRows([]string{"NationalCode", "FirstName", "LastName", "Gender"}).
		AddRow("2550000000", "fname1", "lname1", "male").
		AddRow("2550000001", "fname2", "lname2", "male").
		AddRow("2550000002", "fname3", "lname3", "female")

	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "ticket_passengers" WHERE "ticket_passengers"."ticket_id" IN \(\$1,\$2\)`).
		WithArgs(1, 2).
		WillReturnRows(mockPassengerRows)

	mockUserRows := suite.sqlMock.NewRows([]string{"FirstName", "LastName", "Email", "PhoneNumber", "Password"}).
		AddRow("fname1", "lname1", "email1@example.com", "09120000000", "12345678").
		AddRow("fname2", "lname2", "email2@example.com", "09120000001", "12345678").
		AddRow("fname3", "lname3", "email3@example.com", "09120000002", "12345678")

	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "users" WHERE "users"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(mockUserRows)

	tickets, err := GetUserTickets(suite.dbMock, suite.UserID)
	require.NoError(err)
	require.NotNil(tickets)
	require.Len(tickets, 2)
}

func TestTicketsRepository(t *testing.T) {
	suite.Run(t, new(TicketTestSuite))
}
