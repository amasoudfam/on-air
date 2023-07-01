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
	NationalCode string `json:"national_code" binding:"required" validate:"required"`
	FirstName    string `json:"first_name" binding:"required" validate:"required"`
	LastName     string `json:"last_name" binding:"required" validate:"required"`
	Gender       string `json:"gender" binding:"required" validate:"required"`
}

func (p *Passenger) Create(ctx echo.Context) error {
	userID, _ := strconv.Atoi(ctx.Get("id").(string))
	var req CreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Failed to bind")
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid json key")
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
	NationalCode string `json:"national_code" binding:"required" validate:"required"`
	FirstName    string `json:"first_name" binding:"required" validate:"required"`
	LastName     string `json:"last_name" binding:"required" validate:"required"`
	Gender       string `json:"gender" binding:"required" validate:"required"`
}

func (p *Passenger) Get(ctx echo.Context) error {
	userID, _ := strconv.Atoi(ctx.Get("id").(string))
	passengers, err := repository.GetPassengersByUserID(p.DB, userID)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal error")
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
