package handlers

import (
	"encoding/json"
	"errors"
	"io"
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

type GetTicketTestSuite struct {
	suite.Suite
	sqlMock  sqlmock.Sqlmock
	e        *echo.Echo
	endpoint string
	ticket   *Ticket
	UserID   int
}

func (suite *GetTicketTestSuite) CallHandler() (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodGet, suite.endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	c.Set("id", strconv.Itoa(suite.UserID))
	err := suite.ticket.GetTickets(c)
	return res, err
}

func (suite *GetTicketTestSuite) SetupSuite() {
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
			ExpiresIn: time.Minute * 3,
		},
	}
	suite.e = echo.New()
	suite.e.Validator = &utils.CustomValidator{Validator: validator.New()}
	suite.endpoint = "/tickets"
	suite.UserID = 1
}

func (suite *GetTicketTestSuite) TestGetTickets_Success() {
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

type GetTicketPDFTestSuite struct {
	suite.Suite
	sqlMock  sqlmock.Sqlmock
	e        *echo.Echo
	endpoint string
	ticket   *Ticket
	UserID   int
}

func (suite *GetTicketPDFTestSuite) SetupSuite() {
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
	suite.ticket = &Ticket{DB: db}
	suite.e = echo.New()
	suite.e.Validator = &utils.CustomValidator{Validator: validator.New()}
	suite.endpoint = "/ticketPDF"
	suite.UserID = 3
}

func (suite *GetTicketPDFTestSuite) CallGetHandler(ticketId int) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodGet, suite.endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	c.Set("id", strconv.Itoa(suite.UserID))
	if ticketId > 0 {
		c.QueryParams().Add("ticket_id", strconv.Itoa(ticketId))
	}
	err := suite.ticket.GetPDF(c)
	return res, err
}

func (suite *GetTicketPDFTestSuite) TestGetList_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	time := time.Now()

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

	patch := monkey.Patch(repository.GetTicket, func(db *gorm.DB, userID int, ticketID int) (models.Ticket, error) {
		return tickets[0], nil
	})

	defer patch.Unpatch()
	res, err := suite.CallGetHandler(2)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *GetTicketPDFTestSuite) TestGetList_Failure_InvalidTicketID() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedBody := "\"Invalid ticket_id\"\n"

	patch := monkey.Patch(repository.GetTicket, func(db *gorm.DB, userID int, ticketID int) (models.Ticket, error) {
		return models.Ticket{}, nil
	})
	defer patch.Unpatch()

	res, err := suite.CallGetHandler(0)
	require.NoError(err)
	body, _ := io.ReadAll(res.Body)
	require.Equal(expectedBody, string(body))
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *GetTicketPDFTestSuite) TestGetList_Failure_InternalError() {
	require := suite.Require()
	expectedStatusCode := http.StatusInternalServerError
	expectedBody := "\"Internal error\"\n"

	patch := monkey.Patch(repository.GetTicket, func(db *gorm.DB, userID int, ticketID int) (models.Ticket, error) {
		return models.Ticket{}, errors.New("Internal error")
	})
	defer patch.Unpatch()

	res, err := suite.CallGetHandler(2)
	require.NoError(err)
	body, _ := io.ReadAll(res.Body)
	require.Equal(expectedBody, string(body))
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *GetTicketPDFTestSuite) TestGetList_Failure_GenerateOutput() {
	require := suite.Require()
	expectedStatusCode := http.StatusInternalServerError
	expectedBody := "\"Internal error\"\n"

	patch1 := monkey.Patch(repository.GetTicket, func(db *gorm.DB, userID int, ticketID int) (models.Ticket, error) {
		return models.Ticket{}, nil
	})
	defer patch1.Unpatch()

	patch2 := monkey.Patch(utils.GeneratePDF, func(ticket models.Ticket) ([]byte, error) {
		return nil, errors.New("Internal error")
	})
	defer patch2.Unpatch()

	res, err := suite.CallGetHandler(2)
	require.NoError(err)
	body, _ := io.ReadAll(res.Body)
	require.Equal(expectedBody, string(body))
	require.Equal(expectedStatusCode, res.Code)
}

func TestGetTicket(t *testing.T) {
	suite.Run(t, new(GetTicketTestSuite))
}

func TestGetTicketPDF(t *testing.T) {
	suite.Run(t, new(GetTicketPDFTestSuite))
}
