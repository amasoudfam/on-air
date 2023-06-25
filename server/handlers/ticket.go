package handlers

import (
	"net/http"
	"on-air/config"
	"on-air/repository"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Ticket struct {
	DB  *gorm.DB
	JWT *config.JWT
}

type ReserveRequest struct {
	Price        int   `json:"price" binding:"required"`
	FlightID     int   `json:"flight_id" binding:"required"`
	PassengerIDs []int `json:"passengers" binding:"required"`
}

type ReserveResponse struct {
	Status string `json:"status" binding:"required"`
}

func (t *Ticket) Reserve(ctx echo.Context) error {
	var req ReserveRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	userId := ctx.Get("id").(int)

	dbUser, err := repository.ReserveTicket(
		t.DB,
		userId,
		req.FlightID,
		req.Price,
		req.PassengerIDs,
	)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server")
	}

	return ctx.JSON(http.StatusOK, ReserveResponse{
		Status: dbUser.Status,
	})
}
