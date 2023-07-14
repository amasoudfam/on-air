package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"on-air/config"
	"on-air/models"
	"on-air/server/handlers"
	"on-air/server/middlewares"
	"on-air/utils"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IntegrationTestSuite struct {
	suite.Suite
	db     *gorm.DB
	e      *echo.Echo
	JWT    config.JWT
	db_cfg config.Database
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.e = echo.New()
	suite.e.Validator = &utils.CustomValidator{Validator: validator.New()}

	suite.JWT = config.JWT{
		SecretKey: "mysecretkey",
		ExpiresIn: time.Minute * 60,
	}

	suite.db_cfg = config.Database{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "postgres",
		DB:       "fake_db",
	}

	suite.initDB()
	suite.initHandlers()
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.dropDB()
}

func (suite *IntegrationTestSuite) initHandlers() {
	authMiddleware := &middlewares.Auth{
		JWT: &suite.JWT,
	}

	auth := &handlers.Auth{
		DB:  suite.db,
		JWT: &suite.JWT,
	}
	suite.e.POST("/auth/register", auth.Register)
	suite.e.POST("/auth/login", auth.Login)

	passenger := &handlers.Passenger{
		DB: suite.db,
	}
	suite.e.POST("/passengers", passenger.Create, authMiddleware.AuthMiddleware)
}

func (suite *IntegrationTestSuite) initDB() {
	dsn := fmt.Sprintf("host= %s user=%s password=%s port=%d sslmode=disable",
		suite.db_cfg.Host,
		suite.db_cfg.Username,
		suite.db_cfg.Password,
		suite.db_cfg.Port,
	)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	result := conn.Exec("CREATE DATABASE fake_db")
	if result.Error != nil {
		panic(result.Error)
	}

	dsn = fmt.Sprintf("host= %s user=%s password=%s port=%d database=%s sslmode=disable",
		suite.db_cfg.Host,
		suite.db_cfg.Username,
		suite.db_cfg.Password,
		suite.db_cfg.Port,
		suite.db_cfg.DB,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Country{})
	db.AutoMigrate(&models.City{})
	db.AutoMigrate(&models.Flight{})
	db.AutoMigrate(&models.Ticket{})
	db.AutoMigrate(&models.Passenger{})
	db.AutoMigrate(&models.Payment{})
	suite.db = db
}

func (suite *IntegrationTestSuite) dropDB() {
	dsn := fmt.Sprintf("host= %s user=%s password=%s port=%d sslmode=disable",
		suite.db_cfg.Host,
		suite.db_cfg.Username,
		suite.db_cfg.Password,
		suite.db_cfg.Port,
	)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := suite.db.DB()
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB.Close()

	result := conn.Exec("DROP DATABASE fake_db")
	if result.Error != nil {
		panic(result.Error)
	}
}

func (suite *IntegrationTestSuite) TestRegister_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusCreated

	newUser := handlers.RegisterRequest{
		Email:    "zynab.sobhani@gmail.com",
		Password: "12345",
	}

	payload, err := json.Marshal(newUser)
	require.NoError(err)

	createReq := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(payload))

	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	res := httptest.NewRecorder()
	suite.e.ServeHTTP(res, createReq)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *IntegrationTestSuite) TestRegister_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedBody := "\"User exist\"\n"

	newUser := handlers.RegisterRequest{
		Email:    "masoud.aghdasifam@gmail.com",
		Password: "12345",
	}

	payload, err := json.Marshal(newUser)
	require.NoError(err)
	createReq := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(payload))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	suite.e.ServeHTTP(res, createReq)

	createReq = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(payload))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res = httptest.NewRecorder()
	suite.e.ServeHTTP(res, createReq)
	body, _ := io.ReadAll(res.Body)
	require.Equal(expectedBody, string(body))
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *IntegrationTestSuite) TestLogin_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK

	newUser := handlers.RegisterRequest{
		Email:    "amin.savari@gmail.com",
		Password: "12345",
	}

	payload, err := json.Marshal(newUser)
	require.NoError(err)
	createReq := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(payload))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	suite.e.ServeHTTP(res, createReq)

	createReq = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(payload))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res = httptest.NewRecorder()
	suite.e.ServeHTTP(res, createReq)
	require.Equal(expectedStatusCode, res.Code)
}

func (suite *IntegrationTestSuite) TestLogin_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusUnauthorized
	expectedBody := "\"Invalid credentials\"\n"

	newUser := handlers.RegisterRequest{
		Email:    "mohammad.serpush@gmail.com",
		Password: "12345",
	}

	payload, err := json.Marshal(newUser)
	require.NoError(err)
	createReq := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(payload))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	suite.e.ServeHTTP(res, createReq)

	newUser.Password = "123456"
	payload, _ = json.Marshal(newUser)
	createReq = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(payload))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res = httptest.NewRecorder()
	suite.e.ServeHTTP(res, createReq)

	body, _ := io.ReadAll(res.Body)
	require.Equal(expectedBody, string(body))
	require.Equal(expectedStatusCode, res.Code)
}

type Message struct {
	Access_token string
	Token_type   string
}

func (suite *IntegrationTestSuite) TestPassenger_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusCreated

	newUser := handlers.RegisterRequest{
		Email:    "amin.savari@gmail.com",
		Password: "12345",
	}

	newPassenger := handlers.CreateRequest{
		NationalCode: "1111111111",
		FirstName:    "Masoud",
		LastName:     "Aghdasifam",
		Gender:       "Male",
	}

	userPayload, err := json.Marshal(newUser)
	require.NoError(err)

	createReq := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(userPayload))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	suite.e.ServeHTTP(res, createReq)

	createReq = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(userPayload))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res = httptest.NewRecorder()
	suite.e.ServeHTTP(res, createReq)
	body, _ := io.ReadAll(res.Body)
	var v Message
	json.Unmarshal(body, &v)

	passengerPayload, err := json.Marshal(newPassenger)
	require.NoError(err)

	createReq = httptest.NewRequest(http.MethodPost, "/passengers", bytes.NewReader(passengerPayload))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	createReq.Header.Set("Authorization", fmt.Sprintf("%s %s", v.Token_type, v.Access_token))
	res = httptest.NewRecorder()
	suite.e.ServeHTTP(res, createReq)

	require.Equal(expectedStatusCode, res.Code)
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
