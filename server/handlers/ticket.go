package handlers

import (
	"net/http"
	"on-air/config"
	"on-air/repository"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Ticket struct {
	DB  *gorm.DB
	JWT *config.JWT
}

type ReserveRequest struct {
	UnitPrice    int   `json:"UnitPrice" binding:"required"`
	FlightID     int   `json:"FlightID" binding:"required"`
	PassengerIDs []int `json:"PassengerIDs" binding:"required"`
}

type ReserveResponse struct {
	Status string `json:"Status" binding:"required"`
}

func (t *Ticket) Reserve(ctx echo.Context) error {
	var req ReserveRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	reqToken := ctx.Request().Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	payLoad, err := repository.VerifyToken(t.JWT, reqToken)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, "Internal server error")
	}

	// TODO repository
	// TODO error package
	dbUser, err := repository.ReserveTicket(
		t.DB,
		payLoad.UserID,
		req.FlightID,
		req.UnitPrice,
		req.PassengerIDs,
	)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server")
	}

	return ctx.JSON(http.StatusOK, ReserveResponse{
		Status: dbUser.Status,
	})
}
