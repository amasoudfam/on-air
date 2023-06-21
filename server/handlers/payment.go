package handlers

import (
	"net/http"
	"on-air/config"
	"on-air/repository"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Payment struct {
	DB  *gorm.DB
	JWT *config.JWT
	IPG *config.IPG
}

type PayRequest struct {
	TicketID uint `json:"TicketID" binding:"required"`
}

type PayResponse struct {
	Address string `json:"Token" binding:"required"`
}

func (t *Payment) Pay(ctx echo.Context) error {
	var req PayRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	reqToken := ctx.Request().Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	_, err := repository.VerifyToken(t.JWT, reqToken)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, "Internal server error")
	}

	// TODO repository
	// TODO error package
	address, err := repository.PayTicket(t.DB, t.IPG, req.TicketID)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return ctx.JSON(http.StatusOK, PayResponse{
		Address: address,
	})
}

type CallBackRequest struct {
	PaymentID uint `json:"PaymentID" binding:"required"`
}

type CallBackResponse struct {
	Status string `json:"Status" binding:"required"`
}

func (t *Payment) CallBack(ctx echo.Context) error {
	var req CallBackRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	// TODO repository
	// TODO error package
	status, err := repository.VerifyPayment(t.DB, t.IPG, req.PaymentID)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return ctx.JSON(http.StatusOK, CallBackResponse{
		Status: status,
	})
}
