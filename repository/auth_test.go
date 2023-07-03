package repository

import (
	"log"
	"on-air/config"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthTestSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock
	dbMock  *gorm.DB
	jwt     *config.JWT
}

func (suite *AuthTestSuite) SetupSuite() {
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
	suite.jwt = &config.JWT{
		SecretKey: "superSafeSecretKey",
		LifeTime:  time.Duration(60 * time.Minute),
	}
}

func (suite *AuthTestSuite) TestAuth_CreateToken_Success() {
	require := suite.Require()

	_, err := CreateToken(suite.jwt, 3)
	require.NoError(err)
}

func (suite *AuthTestSuite) TestAuth_VerifyToken_Success() {
	require := suite.Require()

	accessToken, _ := CreateToken(suite.jwt, 1)

	payload, err := VerifyToken(suite.jwt, accessToken)
	require.NoError(err)
	require.Equal(payload.UserID, 1)
}

func (suite *AuthTestSuite) TestAuth_VerifyToken_Failure() {
	require := suite.Require()

	payload, err := VerifyToken(suite.jwt, "supermozakhraftoken")
	require.Error(err)
	require.Nil(payload)
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
