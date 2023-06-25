package utils

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type NationalCodeTestSuite struct {
	suite.Suite
}

func (suite *NationalCodeTestSuite) TestValidateNationalCode() {
	require := suite.Require()
	testCases := []struct {
		desc           string
		nationalCode   string
		expectedResult bool
	}{
		{
			"Valid national code",
			"0025252925",
			true,
		},
		{
			"Invalid national code length",
			"1234",
			false,
		},
		{
			"Invalid national code",
			"123456789",
			false,
		},
	}

	for _, t := range testCases {
		res := ValidateNationalCode(t.nationalCode)
		require.EqualValues(t.expectedResult, res)
	}
}

func TestNationalCode(t *testing.T) {
	suite.Run(t, new(NationalCodeTestSuite))
}
