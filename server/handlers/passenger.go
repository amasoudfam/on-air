package handlers

import (
	"net/http"
	"on-air/repository"
	"on-air/utils"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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

type CreateResponse struct {
	Status  bool
	Message string
}

func (p *Passenger) Create(ctx echo.Context) error {
	userID, _ := ctx.Get("user_id").(int)
	var req CreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
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
			logrus.Error("passenger_handler: Create failed when use repository.CreatePassenger, error:", err)
			return ctx.JSON(http.StatusInternalServerError, "Internal server error")
		}
	}

	return ctx.JSON(http.StatusCreated, CreateResponse{
		Status:  true,
		Message: "New user create successfully",
	})
}

type GetResponse struct {
	NationalCode string `json:"national_code" binding:"required" validate:"required"`
	FirstName    string `json:"first_name" binding:"required" validate:"required"`
	LastName     string `json:"last_name" binding:"required" validate:"required"`
	Gender       string `json:"gender" binding:"required" validate:"required"`
}

func (p *Passenger) Get(ctx echo.Context) error {
	userID, _ := ctx.Get("user_id").(int)
	passengers, err := repository.GetPassengersByUserID(p.DB, userID)

	if err != nil {
		logrus.Error("passenger_handler: Get failed when use repository.GetPassengersByUserID, error:", err)
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
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
