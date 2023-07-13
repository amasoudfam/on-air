package handlers

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/http/httptest"
	"on-air/config"
	"on-air/utils"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthTestSuite struct {
	suite.Suite
	sqlMock  sqlmock.Sqlmock
	e        *echo.Echo
	endpoint string
	auth     *Auth
}

func (suite *AuthTestSuite) SetupSuite() {
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
	suite.auth = &Auth{DB: db, JWT: &config.JWT{
		SecretKey: "testSecret",
		ExpiresIn: time.Minute * 3,
	}}

	suite.e = echo.New()
	suite.e.Validator = &utils.CustomValidator{Validator: validator.New()}
	suite.endpoint = "/auth"
}

func (suite *AuthTestSuite) SetupTest() {
	suite.e.Validator = &utils.CustomValidator{Validator: validator.New()}
	suite.e.Binder = &echo.DefaultBinder{}
}

func (suite *AuthTestSuite) CallRegisterHandler(requestBody string) (*httptest.ResponseRecorder, error) {
	endpoint := fmt.Sprintf("%s/register", suite.endpoint)
	req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	err := suite.auth.Register(c)
	return res, err
}

func (suite *AuthTestSuite) CallLoginHandler(requestBody string) (*httptest.ResponseRecorder, error) {
	endpoint := fmt.Sprintf("%s/login", suite.endpoint)

	req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	err := suite.auth.Login(c)
	return res, err
}

func (suite *AuthTestSuite) TestAuth_Register_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusCreated

	suite.e.Binder = &MockBinder{}

	suite.e.Validator = &MockValidator{}

	monkey.Patch(utils.HashPassword, func(_ string) (string, error) {
		return "superHashedPasswor", nil
	})
	defer monkey.Unpatch(utils.HashPassword)

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(
		`INSERT`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	suite.sqlMock.ExpectCommit()

	requestBody := `{"email": "admin@gmail.com" , "password" : "adminadmin"}`
	res, err := suite.CallRegisterHandler(requestBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *AuthTestSuite) TestAuth_Register_Failure_Invalid_Body() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest

	suite.e.Binder = &MockBinder{}

	suite.e.Validator = &MockValidator{}

	requestBody := `{}`
	res, err := suite.CallRegisterHandler(requestBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *AuthTestSuite) TestAuth_Register_Failure_Duplicate_User() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest

	suite.e.Binder = &MockBinder{}

	suite.e.Validator = &MockValidator{}

	pgErr := &pgconn.PgError{
		Message: "User exists",
		Code:    "23505",
	}

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(`INSERT`).
		WillReturnError(pgErr)
	suite.sqlMock.ExpectRollback()

	requestBody := `{"email" : "admin@gmail.com" , "password" : "admin"}`
	res, err := suite.CallRegisterHandler(requestBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)

}

func (suite *AuthTestSuite) TestAuth_Login_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK

	suite.e.Binder = &MockBinder{}

	suite.e.Validator = &MockValidator{}

	monkey.Patch(bcrypt.CompareHashAndPassword, func(_, _ []byte) error {
		return nil
	})
	defer monkey.Unpatch(bcrypt.CompareHashAndPassword)

	mockUser := suite.sqlMock.NewRows(
		[]string{
			"id", "email", "password",
		}).
		AddRow("1", "admin@gmail.com", "admin")
	suite.sqlMock.ExpectQuery(`SELECT`).
		WillReturnRows(mockUser)

	requestBody := `{"email": "admin@gmail.com" , "password" : "admin"}`
	res, err := suite.CallLoginHandler(requestBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
