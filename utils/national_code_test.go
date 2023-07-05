package utils

import (
	"testing"

	"github.com/stretchr/testify/suite"
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
			"Invalid national code checksum character",
			"002525292a",
			false,
		},
		{
			"Invalid national code character",
			"a025252920",
			false,
		},
		{
			"Invalid national code",
			"1234567890",
			false,
		},
		{
			"Empty national code",
			"",
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
