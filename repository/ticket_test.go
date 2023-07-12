package repository

import (
	"log"
	"on-air/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var penalties datatypes.JSON = datatypes.JSON([]byte(`{"test":"on-air"}`))

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
	suite.UserID = 1
}

// func (suite *TicketTestSuite) TestTicket_ReserveTicket_Success() {
// 	require := suite.Require()

// 	mockPassengerRows := suite.sqlMock.NewRows([]string{"id", "user_id", "national_code", "first_name", "last_name", "gender"}).
// 		AddRow(1, 1, "2550000000", "pfname1", "plname1", "male").
// 		AddRow(2, 1, "2550000001", "pfname2", "plname2", "male").
// 		AddRow(3, 1, "2550000002", "pfname3", "plname3", "male")
// 	suite.sqlMock.ExpectQuery(`(?i)SELECT\s+.+\s+FROM\s+"passengers"\s+WHERE\s+id\s+IN\s+\(\$1,\$2,\$3\)\s+AND\s+"passengers"\."deleted_at"\s+IS\s+NULL`).
// 		WithArgs(1, 2, 3).
// 		WillReturnRows(mockPassengerRows)

// 	// suite.sqlMock.ExpectBegin()
// 	// suite.sqlMock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "tickets"`)).
// 	// 	WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
// 	// 	WillReturnResult(sqlmock.NewResult(1, 1))

// 	suite.sqlMock.ExpectCommit()
// 	_, err := ReserveTicket(suite.dbMock, int(suite.UserID), 1, 10000, []int{1, 2, 3})
// 	require.NoError(err)
// }

func (suite *TicketTestSuite) TestTicket_GetExpiredTickets_Success() {
	require := suite.Require()

	mockPassenger := suite.sqlMock.NewRows(
		[]string{
			"id", "user_id", "unit_price", "count", "flight_id", "status", "created_at",
		}).
		AddRow(10, 1, "1000000", 9, 5, "Reserved", time.Now().Add(-20*time.Minute))

	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "tickets"`).
		WillReturnRows(mockPassenger)

	data, err := GetExpiredTickets(suite.dbMock)

	require.NoError(err)
	require.Len(data, 1)
}

func (suite *TicketTestSuite) TestTickets_GetTickets_Success() {
	require := suite.Require()
	data := []models.Ticket{
		{
			UserID:    1,
			UnitPrice: 100,
			Count:     2,
			FlightID:  1,
			Status:    "complete",
			User: models.User{
				FirstName:   "fname1",
				LastName:    "lname1",
				Email:       "email1@example.com",
				PhoneNumber: "09120000000",
				Password:    "12345678",
			},
			Flight: models.Flight{
				Number:     "F101",
				FromCityID: uint(1),
				ToCityID:   uint(2),
				Airplane:   "Aseman",
				Airline:    "f12",
				Penalties:  penalties,
				FromCity: models.City{
					Name:      "Tehran",
					CountryID: uint(1),
					Country: models.Country{
						Name: "Iran",
					},
				},
				ToCity: models.City{
					Name:      "Shiraz",
					CountryID: uint(1),
					Country: models.Country{
						Name: "Iran",
					},
				},
			},
			Passengers: []models.Passenger{
				{
					UserID:       uint(1),
					NationalCode: "2550000000",
					FirstName:    "pfname1",
					LastName:     "plname1",
					Gender:       "male",
				},
			},
		},
	}

	ticket1 := &data[0]
	ticket1.ID = uint(1)
	ticket1.User.ID = uint(1)
	ticket1.Flight.ID = uint(1)
	ticket1.Flight.FromCity.ID = uint(1)
	ticket1.Flight.FromCity.Country.ID = uint(1)
	ticket1.Flight.ToCity.ID = uint(2)
	ticket1.Flight.ToCity.Country.ID = uint(1)
	ticket1.Passengers[0].ID = uint(1)

	mockTicketRows := suite.sqlMock.NewRows([]string{"id", "unit_price", "flight_id", "count", "status", "user_id"}).
		AddRow(1, 100, 1, 2, "complete", 1)
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "tickets" WHERE user_id = \$1 AND (.+)`).
		WithArgs(suite.UserID).
		WillReturnRows(mockTicketRows)

	mockFlightRows := suite.sqlMock.NewRows([]string{"id", "number", "from_city_id", "to_city_id", "airplane", "airline", "penalties"}).
		AddRow(1, "F101", 1, 2, "Aseman", "f12", penalties)
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "flights" WHERE "flights"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(mockFlightRows)

	mockCityRows := suite.sqlMock.NewRows([]string{"id", "name", "country_id"}).
		AddRow(1, "Tehran", 1)
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "cities" WHERE "cities"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(mockCityRows)

	mockCountryRows := suite.sqlMock.NewRows([]string{"id", "name"}).
		AddRow(1, "Iran")
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "countries" WHERE "countries"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(mockCountryRows)

	mockToCityRows := suite.sqlMock.NewRows([]string{"id", "name", "country_id"}).
		AddRow(2, "Shiraz", 1)
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "cities" WHERE "cities"."id" = \$1`).
		WithArgs(2).
		WillReturnRows(mockToCityRows)

	mockToCountryRows := suite.sqlMock.NewRows([]string{"id", "name"}).
		AddRow(1, "Iran")
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "countries" WHERE "countries"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(mockToCountryRows)

	mockTicketPassengersRows := suite.sqlMock.NewRows([]string{"ticket_id", "passenger_id"}).
		AddRow(1, 1)
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "ticket_passengers" WHERE "ticket_passengers"."ticket_id" = \$1`).
		WithArgs(1).
		WillReturnRows(mockTicketPassengersRows)

	mockPassengerRows := suite.sqlMock.NewRows([]string{"id", "user_id", "national_code", "first_name", "last_name", "gender"}).
		AddRow(1, 1, "2550000000", "pfname1", "plname1", "male")
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "passengers" WHERE "passengers"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(mockPassengerRows)

	mockUserRows := suite.sqlMock.NewRows([]string{"id", "first_name", "last_name", "email", "phone_number", "password", "deleted_at"}).
		AddRow(1, "fname1", "lname1", "email1@example.com", "09120000000", "12345678", nil)
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "users" WHERE "users"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(mockUserRows)

	tickets, err := GetUserTickets(suite.dbMock, suite.UserID)
	require.NoError(err)
	require.NotNil(tickets)
	require.Len(tickets, 1)
	require.Equal(data, tickets)
}

func TestTicketsRepository(t *testing.T) {
	suite.Run(t, new(TicketTestSuite))
}
