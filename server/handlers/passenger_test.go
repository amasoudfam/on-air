package handlers

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"on-air/utils"
	"strings"
	"testing"

	"bou.ke/monkey"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CreatePassengerTestSuite struct {
	suite.Suite
	sqlMock   sqlmock.Sqlmock
	e         *echo.Echo
	endpoint  string
	passenger *Passenger
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (suite *CreatePassengerTestSuite) SetupSuite() {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}))

	if err != nil {
		log.Fatal(err)
	}

	suite.sqlMock = sqlMock
	suite.passenger = &Passenger{DB: db}
	suite.e = echo.New()
	suite.e.Validator = &CustomValidator{validator: validator.New()}
	suite.endpoint = "/passenger"
}

func (suite *CreatePassengerTestSuite) CallCreateHandler(requestBody string) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodPost, suite.endpoint, strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	c.Set("id", "2")
	err := suite.passenger.Create(c)
	return res, err
}

func (suite *CreatePassengerTestSuite) CallGetHandler() (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodPost, suite.endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	c.Set("id", "2")
	err := suite.passenger.Get(c)
	return res, err
}

func (suite *CreatePassengerTestSuite) TestCreatePassenger_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusCreated

	monkey.Patch(utils.ValidateNationalCode, func(_ string) bool {
		return true
	})

	defer monkey.Unpatch(utils.ValidateNationalCode)

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("INSERT INTO `passengers`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.sqlMock.ExpectCommit()

	requestBody := `{"nationalcode": "0123456789","firstname": "name","lastname": "lname","gender": "f"}`
	res, err := suite.CallCreateHandler(requestBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *CreatePassengerTestSuite) TestCreatePassenger_ParseReq_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	requestBody := `{"nationalcode: "1000011111","firstname": "name","lastname": "lname","gender": "f"}`

	res, err := suite.CallCreateHandler(requestBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *CreatePassengerTestSuite) TestCreatePassenger_ValidateNationalCode_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedErr := "\"Invalid national code\"\n"
	requestBody := `{"nationalcode": "323","firstname": "name","lastname": "lname","gender": "f"}`

	res, err := suite.CallCreateHandler(requestBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	body, _ := io.ReadAll(res.Body)
	require.Equal(expectedErr, string(body))
}

func (suite *CreatePassengerTestSuite) TestListPassenger() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK

	mockRow := suite.sqlMock.NewRows(
		[]string{
			"nationalcode", "firstname", "lastname", "user_id", "gender",
		}).
		AddRow("1000011111", "name", "lname", 2, "f")

	suite.sqlMock.ExpectQuery("SELECT (.+) FROM `passengers`").
		WillReturnRows(mockRow)
	res, err := suite.CallGetHandler()
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func TestCreatePassenger(t *testing.T) {
	suite.Run(t, new(CreatePassengerTestSuite))
}
