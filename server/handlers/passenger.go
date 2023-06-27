package handlers

import (
	"net/http"
	"on-air/repository"
	"on-air/utils"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Passenger struct {
	DB *gorm.DB
}

type CreateRequest struct {
	NationalCode string `json:"nationalcode" binding:"required" validate:"required"`
	FirstName    string `json:"firstname" binding:"required" validate:"required"`
	LastName     string `json:"lastname" binding:"required" validate:"required"`
	Gender       string `json:"gender" binding:"required" validate:"required"`
}

func (p *Passenger) Create(ctx echo.Context) error {
	userID, _ := strconv.Atoi(ctx.Get("id").(string))
	var req CreateRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if !utils.ValidateNationalCode(req.NationalCode) {
		return ctx.JSON(http.StatusBadRequest, "Invalid national code")
	}

	_, err := repository.CreatePassenger(
		p.DB,
		userID,
		req.NationalCode,
		req.FirstName,
		req.LastName,
		req.Gender,
	)

	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return ctx.JSON(http.StatusBadRequest, "Passenger exists")
		} else {
			return ctx.JSON(http.StatusInternalServerError, "Internal error")
		}
	}

	return ctx.JSON(http.StatusCreated, nil)
}

type GetResponse struct {
	NationalCode string `json:"nationalcode" binding:"required" validate:"required"`
	FirstName    string `json:"firstname" binding:"required" validate:"required"`
	LastName     string `json:"lastname" binding:"required" validate:"required"`
	Gender       string `json:"gender" binding:"required" validate:"required"`
}

func (p *Passenger) Get(ctx echo.Context) error {
	userID, _ := strconv.Atoi(ctx.Get("id").(string))
	passengers, err := repository.GetPassengersByUserID(p.DB, userID)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	var response []GetResponse
	if len(*passengers) > 0 {
		response = make([]GetResponse, 0, len(*passengers))
		for _, p := range *passengers {
			response = append(response, GetResponse{
				FirstName:    p.FirstName,
				LastName:     p.LastName,
				NationalCode: p.NationalCode,
				Gender:       p.Gender,
			})
		}
	}
	return ctx.JSON(http.StatusOK, response)
}
