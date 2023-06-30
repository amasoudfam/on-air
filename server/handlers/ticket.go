package handlers

import (
	"net/http"
	"on-air/config"
	"on-air/repository"
	"on-air/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Ticket struct {
	DB  *gorm.DB
	JWT *config.JWT
}

type ReserveRequest struct {
	FlightNumber string `json:"flight_number" binding:"required"`
	PassengerIDs []int  `json:"passengers" binding:"required"`
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

	userId := ctx.Get("user_id").(int)

	//TODO: firstly find from DB
	var flighInfo, err = services.GetInfo(req.FlightNumber)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server")
	}

	//TODO: Get all data from mock server
	flight, err := repository.AddFlight(t.DB, req.FlightNumber)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server Error")
	}

	flightReserve, err := services.Reserve(req.FlightNumber, len(req.PassengerIDs))

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server Error")
	}

	if !flightReserve.Reserved {
		return ctx.JSON(http.StatusInternalServerError, "Sold out")
	}

	ticket, err := repository.ReserveTicket(
		t.DB,
		userId,
		int(flight.ID),
		flighInfo.Price,
		req.PassengerIDs,
	)

	if err != nil {
		services.Refund(req.FlightNumber, len(req.PassengerIDs))
		return ctx.JSON(http.StatusInternalServerError, "Internal server Error")
	}

	return ctx.JSON(http.StatusOK, ReserveResponse{
		Status: ticket.Status,
	})
}
