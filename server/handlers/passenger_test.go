package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"on-air/utils"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"bou.ke/monkey"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PassengerTestSuite struct {
	suite.Suite
	sqlMock   sqlmock.Sqlmock
	e         *echo.Echo
	endpoint  string
	passenger *Passenger
	UserID    int
}

type MockValidator struct {
}

func (mcv *MockValidator) Validate(_ interface{}) error {
	return nil
}

type MockBinder struct {
}

func (mcv *MockBinder) Bind(_ interface{}, _ echo.Context) error {
	return nil
}

func (suite *PassengerTestSuite) SetupSuite() {
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
	suite.passenger = &Passenger{DB: db}
	suite.e = echo.New()
	suite.e.Validator = &utils.CustomValidator{Validator: validator.New()}
	suite.endpoint = "/passenger"
	suite.UserID = 3
}

func (suite *PassengerTestSuite) SetupTest() {
	suite.e.Validator = &utils.CustomValidator{Validator: validator.New()}
	suite.e.Binder = &echo.DefaultBinder{}
}

func (suite *PassengerTestSuite) CallCreateHandler(requestBody string) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodPost, suite.endpoint, strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	c.Set("id", strconv.Itoa(suite.UserID))
	err := suite.passenger.Create(c)
	return res, err
}

func (suite *PassengerTestSuite) CallGetHandler() (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodGet, suite.endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	c.Set("id", strconv.Itoa(suite.UserID))
	err := suite.passenger.Get(c)
	return res, err
}

func (suite *PassengerTestSuite) TestCreatePassenger_CreatePassenger_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusCreated

	suite.e.Binder = &MockBinder{}

	suite.e.Validator = &MockValidator{}

	monkey.Patch(utils.ValidateNationalCode, func(_ string) bool {
		return true
	})
	defer monkey.Unpatch(utils.ValidateNationalCode)

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(
		regexp.QuoteMeta(`
		  INSERT INTO "passengers" ("created_at","updated_at","deleted_at","user_id","national_code","first_name","last_name","gender")
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 `)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	suite.sqlMock.ExpectCommit()

	requestBody := `{"national_code": "0123456789"}`
	res, err := suite.CallCreateHandler(requestBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *PassengerTestSuite) TestCreatePassenger_CreatePassenger_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusInternalServerError
	expectedBody := "\"Internal server error\"\n"

	suite.e.Binder = &MockBinder{}

	suite.e.Validator = &MockValidator{}

	monkey.Patch(utils.ValidateNationalCode, func(_ string) bool {
		return true
	})
	defer monkey.Unpatch(utils.ValidateNationalCode)

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(
		regexp.QuoteMeta(`
		  INSERT INTO "passengers" ("created_at","updated_at","deleted_at","user_id","national_code","first_name","last_name","gender")
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 `)).
		WillReturnError(errors.New("Internal server error"))
	suite.sqlMock.ExpectRollback()

	requestBody := `{"national_code": "0123456789"}`
	res, err := suite.CallCreateHandler(requestBody)
	body, _ := io.ReadAll(res.Body)
	require.Equal(expectedBody, string(body))
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *PassengerTestSuite) TestCreatePassenger_CreatePassenger_Duplicate_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedBody := "\"Passenger exists\"\n"

	suite.e.Binder = &MockBinder{}

	suite.e.Validator = &MockValidator{}

	monkey.Patch(utils.ValidateNationalCode, func(_ string) bool {
		return true
	})
	defer monkey.Unpatch(utils.ValidateNationalCode)

	pgErr := &pgconn.PgError{
		Message: "Passenger exists",
		Code:    "23505",
	}

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(
		regexp.QuoteMeta(`
		  INSERT INTO "passengers" ("created_at","updated_at","deleted_at","user_id","national_code","first_name","last_name","gender")
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 `)).
		WillReturnError(pgErr)
	suite.sqlMock.ExpectRollback()

	requestBody := `{"national_code": "0123456789"}`
	res, err := suite.CallCreateHandler(requestBody)
	body, _ := io.ReadAll(res.Body)
	require.Equal(expectedBody, string(body))
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *PassengerTestSuite) TestCreatePassenger_InvalidBody_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedBody := "\"Bind Error\"\n"

	requestBody := `{"national_code: "1000011111", "first_name": "name", "last_name": "lname", "gender": "f"}`

	res, err := suite.CallCreateHandler(requestBody)
	body, _ := io.ReadAll(res.Body)
	require.NoError(err)
	require.Equal(expectedBody, string(body))
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *PassengerTestSuite) TestCreatePassenger_InvalidValue_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedBody := "\"Bind Error\"\n"

	requestBody := `{"national_code": 1382122489, "firstname": "", "last_name": "lname", "gender": "f"}`

	res, err := suite.CallCreateHandler(requestBody)
	body, _ := io.ReadAll(res.Body)
	require.NoError(err)
	require.Equal(expectedBody, string(body))
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *PassengerTestSuite) TestCreatePassenger_InvalidKey_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedBody := "Error:Field validation for 'NationalCode'"

	suite.e.Binder = &MockBinder{}
	requestBody := `{"national_code": "1000011111", "firstname": "name", "last_name": "lname", "gender": "f"}`

	res, err := suite.CallCreateHandler(requestBody)
	body, err := io.ReadAll(res.Body)
	require.Contains(string(body), expectedBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *PassengerTestSuite) TestCreatePassenger_ValidateNationalCode_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedErr := "\"Invalid national code\"\n"

	suite.e.Binder = &MockBinder{}

	suite.e.Validator = &MockValidator{}

	requestBody := `{"national_code": "1234567890", "firstname": "fname", "last_name": 8, "gender": "f"}`

	res, err := suite.CallCreateHandler(requestBody)
	body, _ := io.ReadAll(res.Body)
	require.Equal(expectedErr, string(body))
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *PassengerTestSuite) TestGetPassenger_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	expectedBody := "[{\"national_code\":\"1000011111\",\"first_name\":\"name\",\"last_name\":\"lname\",\"gender\":\"f\"}"
	expectedBody += ",{\"national_code\":\"1002011111\",\"first_name\":\"fname\",\"last_name\":\"lname\",\"gender\":\"m\"}]\n"

	mockPassenger := suite.sqlMock.NewRows(
		[]string{
			"national_code", "first_name", "last_name", "gender",
		}).
		AddRow("1000011111", "name", "lname", "f").
		AddRow("1002011111", "fname", "lname", "m")
	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "passengers" WHERE user_id = (.+)`).
		WillReturnRows(mockPassenger)

	res, err := suite.CallGetHandler()
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	body, _ := io.ReadAll(res.Body)
	require.Equal(expectedBody, string(body))
}

func (suite *PassengerTestSuite) TestGetPassenger_Failure() {
	require := suite.Require()
	expectedBody := "\"Internal server error\"\n"
	expectedStatusCode := http.StatusInternalServerError

	suite.sqlMock.ExpectQuery(`SELECT (.+) FROM "passengers" WHERE user_id = (.+)`).
		WillReturnError(errors.New("Internal server error"))

	res, _ := suite.CallGetHandler()
	body, _ := io.ReadAll(res.Body)
	require.Equal(expectedBody, string(body))
	require.Equal(expectedStatusCode, res.Code)
}

func TestPassenger(t *testing.T) {
	suite.Run(t, new(PassengerTestSuite))
}
