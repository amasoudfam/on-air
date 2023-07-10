package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"on-air/config"
	"on-air/models"
	"on-air/repository"
	"on-air/utils"
	"strconv"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TicketTestSuite struct {
	suite.Suite
	sqlMock  sqlmock.Sqlmock
	e        *echo.Echo
	endpoint string
	ticket   *Ticket
	UserID   int
}

func (suite *TicketTestSuite) CallHandler() (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodGet, suite.endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	c.Set("id", strconv.Itoa(suite.UserID))
	err := suite.ticket.GetTickets(c)
	return res, err
}

func (suite *TicketTestSuite) SetupSuite() {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}))

	if err != nil {
		log.Fatal(err)
	}

	suite.sqlMock = sqlMock
	suite.ticket = &Ticket{
		DB: db,
		JWT: &config.JWT{
			SecretKey: "testSecret",
			LifeTime:  time.Minute * 3,
		},
	}
	suite.e = echo.New()
	suite.e.Validator = &utils.CustomValidator{Validator: validator.New()}
	suite.endpoint = "/tickets"
	suite.UserID = 1
}

func (suite *TicketTestSuite) TestGetTickets_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	time := time.Now()
	data := []TicketResponse{
		{
			ID:        1,
			UnitPrice: 1200000,
			Count:     1,
			Status:    "complete",
			CreatedAt: time.Format("2006-01-02 15:04"),
			User: UserResponse{
				FirstName:   "user_fname",
				LastName:    "user_lname",
				Email:       "user@example.com",
				PhoneNumber: "09122222222",
			},
			Flight: FlightResponse{
				Number:     "F102",
				Airplane:   "F12",
				Airline:    "Aseman",
				StartedAt:  time.Format("2006-01-02 15:04"),
				FinishedAt: time.Format("2006-01-02 15:04"),
				FromCity: CityResponse{
					Name: "Tehran",
					Country: CountryResponse{
						Name: "Iran",
					},
				},
				ToCity: CityResponse{
					Name: "Shiraz",
					Country: CountryResponse{
						Name: "Iran",
					},
				},
			},
			Passengers: []PassengerResponse{
				{
					NationalCode: "2550000000",
					FirstName:    "p1_fname",
					LastName:     "p1_lname",
					Gender:       "gmail",
				},
			},
		},
	}

	tickets := []models.Ticket{
		{
			UnitPrice: 1200000,
			Count:     1,
			Status:    "complete",
			User: models.User{
				FirstName:   "user_fname",
				LastName:    "user_lname",
				Email:       "user@example.com",
				PhoneNumber: "09122222222",
			},
			Flight: models.Flight{
				Number:     "F102",
				Airplane:   "F12",
				Airline:    "Aseman",
				StartedAt:  time,
				FinishedAt: time,
				FromCity: models.City{
					Name: "Tehran",
					Country: models.Country{
						Name: "Iran",
					},
				},
				ToCity: models.City{
					Name: "Shiraz",
					Country: models.Country{
						Name: "Iran",
					},
				},
			},
			Passengers: []models.Passenger{
				{
					NationalCode: "2550000000",
					FirstName:    "p1_fname",
					LastName:     "p1_lname",
					Gender:       "gmail",
				},
			},
		},
	}

	tickets[0].ID = 1
	tickets[0].CreatedAt = time

	patch := monkey.Patch(repository.GetUserTickets, func(db *gorm.DB, userID uint) ([]models.Ticket, error) {
		return tickets, nil
	})

	defer patch.Unpatch()

	expectedJSON, _ := json.Marshal(data)
	expectedMsg := string(expectedJSON) + "\n"
	res, err := suite.CallHandler()
	require.NoError(err)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.Equal(expectedMsg, res.Body.String())
}

func TestTicket(t *testing.T) {
	suite.Run(t, new(TicketTestSuite))
}
