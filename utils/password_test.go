package utils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PasswordTestSuite struct {
	suite.Suite
}

func (suite *PasswordTestSuite) TestHashPassword() {
	require := suite.Require()

	hashedPassword, err := HashPassword("Aa123!@#456")
	require.NoError(err)
	require.NotNil(hashedPassword)

}

func TestPassword(t *testing.T) {
	suite.Run(t, new(NationalCodeTestSuite))
}
