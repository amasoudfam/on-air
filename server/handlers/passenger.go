package handlers

import (
	"net/http"
	"on-air/repository"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Passenger struct {
	DB *gorm.DB
}

type PassengerAddRequest struct {
	NationalCode string `json:"nationalcode" binding:"required"`
	FirstName    string `json:"firstname" binding:"required"`
	LastName     string `json:"lastname" binding:"required"`
	UserID       int    `json:"userid" binding:"required"`
	Gender       string `json:"gender" binding:"required"`
}

type UserPassengerListRequest struct {
	UserID int `json:"userid" binding:"required"`
}

type UserPassengerListResponse struct {
	NationalCode string `json:"nationalcode" binding:"required"`
	FirstName    string `json:"firstname" binding:"required"`
	LastName     string `json:"lastname" binding:"required"`
	Gender       string `json:"gender" binding:"required"`
}

func (p *Passenger) PassengerAdd(ctx echo.Context) error {
	passenger := new(PassengerAddRequest)
	if err := ctx.Bind(passenger); err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}

	if err := ctx.Validate(passenger); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	dbPassenger, _ := repository.GetPassengerByNationalCodeAndUserID(p.DB, passenger.NationalCode, passenger.UserID)
	if dbPassenger != nil {
		return ctx.JSON(http.StatusBadRequest, "Passenger exists for this user")
	}

	_, err := repository.AddPassenger(p.DB, passenger.UserID, passenger.NationalCode, passenger.FirstName, passenger.LastName, passenger.Gender)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusCreated, nil)
}

func (p *Passenger) PassengerListByUser(ctx echo.Context) error {
	req := new(UserPassengerListRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	dbPassengers, _ := repository.GetPassengersByUserID(p.DB, req.UserID)
	return ctx.JSON(http.StatusCreated, dbPassengers)
}
