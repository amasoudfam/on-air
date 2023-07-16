package handlers

import (
	"net/http"
	"on-air/config"
	"on-air/repository"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Payment struct {
	DB  *gorm.DB
	IPG *config.IPG
}

type PayRequest struct {
	TicketID uint `json:"ticket_id" binding:"required"`
}

type PayResponse struct {
	Address string `json:"gate_way_url" binding:"required"`
}

func (t *Payment) Pay(ctx echo.Context) error {
	var req PayRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	address, err := repository.PayTicket(t.DB, t.IPG, req.TicketID)
	if err != nil {
		logrus.Error("payment_handler: Pay failed when use repository.PayTicket, error:", err)
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return ctx.JSON(http.StatusOK, PayResponse{
		Address: address,
	})
}

type CallBackRequest struct {
	PaymentID              int
	PaymentDate            time.Time
	TransactionReferenceID int
}

type CallBackResponse struct {
	Status string `json:"status" binding:"required"`
}

func (t *Payment) CallBack(ctx echo.Context) error {
	var req CallBackRequest
	req.PaymentID, _ = strconv.Atoi(ctx.Request().URL.Query().Get("iN"))
	req.PaymentDate, _ = time.Parse("2006/01/02", ctx.Request().URL.Query().Get("iD"))
	req.TransactionReferenceID, _ = strconv.Atoi(ctx.Request().URL.Query().Get("tref"))

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}

	status, err := repository.VerifyPayment(t.DB, t.IPG, req.PaymentID, req.PaymentDate, req.TransactionReferenceID)
	if err != nil {
		logrus.Error("payment_handler: CallBack failed when use repository.VerifyPayment, error:", err)
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return ctx.JSON(http.StatusOK, CallBackResponse{
		Status: status,
	})
}
