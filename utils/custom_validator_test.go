package utils

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
)

type CustomeValidatorTestSuite struct {
	suite.Suite
}

type TestStruct struct {
	FullName string `validate:"required"`
	NickName string
}

func (suite *CustomeValidatorTestSuite) TestCustomeValidator_Success() {
	require := suite.Require()
	validator := &CustomValidator{Validator: validator.New()}

	mockStruct := TestStruct{
		FullName: "full",
		NickName: "nick",
	}
	err := validator.Validate(&mockStruct)

	require.NoError(err)
}

func (suite *CustomeValidatorTestSuite) TestCustomeValidator_NotRequired_Success() {
	require := suite.Require()
	validator := &CustomValidator{Validator: validator.New()}

	mockStruct := TestStruct{
		FullName: "full",
	}
	err := validator.Validate(&mockStruct)

	require.NoError(err)
}

func (suite *CustomeValidatorTestSuite) TestCustomeValidator_Failure() {
	require := suite.Require()
	validator := &CustomValidator{Validator: validator.New()}

	mockStruct := TestStruct{}
	err := validator.Validate(&mockStruct)

	require.Error(err)
}

func (suite *CustomeValidatorTestSuite) TestCustomeValidator_NotRequired_Failure() {
	require := suite.Require()
	validator := &CustomValidator{Validator: validator.New()}

	mockStruct := TestStruct{
		NickName: "nick",
	}
	err := validator.Validate(&mockStruct)

	require.Error(err)
}

func TestCustomeValidator(t *testing.T) {
	suite.Run(t, new(CustomeValidatorTestSuite))
}
