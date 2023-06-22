package handlers

import (
	"net/http"
	"on-air/repository"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Passenger struct {
	DB *gorm.DB
}

type CreateRequest struct {
	NationalCode string `json:"nationalcode" binding:"required"`
	FirstName    string `json:"firstname" binding:"required"`
	LastName     string `json:"lastname" binding:"required"`
	UserID       int    `json:"user_id" binding:"required"`
	Gender       string `json:"gender" binding:"required"`
}

func (p *Passenger) Create(ctx echo.Context) error {
	var req CreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	_, err := repository.CreatePassenger(
		p.DB,
		req.UserID,
		req.NationalCode,
		req.FirstName,
		req.LastName,
		req.Gender,
	)

	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return ctx.JSON(http.StatusBadRequest, "Passenger exists")
		} else {
			return ctx.JSON(http.StatusBadRequest, "Internal error")
		}
	}

	return ctx.JSON(http.StatusCreated, nil)
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

func (p *Passenger) Get(ctx echo.Context) error {
	var req UserPassengerListRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	passengers, _ := repository.GetPassengersByUserID(p.DB, req.UserID)
	response := make([]UserPassengerListResponse, 0, len(*passengers))
	for _, p := range *passengers {
		response = append(response, UserPassengerListResponse{
			FirstName:    p.FirstName,
			LastName:     p.LastName,
			NationalCode: p.NationalCode,
			Gender:       p.Gender,
		})
	}

	return ctx.JSON(http.StatusOK, response)
}
