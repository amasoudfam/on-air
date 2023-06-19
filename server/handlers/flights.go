package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Flight struct {
	DB *gorm.DB
}

type ListRequest struct {
	Origin      string    `json:"origin" validate:"required"`
	Destination string    `json:"destination" validate:"required"`
	Date        time.Time `json:"date" query:"date" validate:"required,datetime=2006-01-02"`
}
type Airplane struct {
	Number      string
	Airplane    string
	Airline     string
	StartedAt   time.Time
	FinishedAt  time.Time
	Origin      string
	Destination string
	Capacity    int
}

type ListResponse struct {
	Flights []Airplane `json:"flights"`
}

func (f *Flight) List(ctx echo.Context) error {
	var req ListRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	// TODO Fetch flight data from the api
	// address := "https://api.example.com/flights"
	// origin := req.Origin
	// destination := req.Destination
	// date := req.Date

	// url := fmt.Sprintf("%s?origin=%s&destination=%s&date=%s", address, origin, destination, date)
	// flights, err := http.Get(url)
	// if err != nil {
	// 	return ctx.JSON(http.StatusInternalServerError, err.Error())
	// }

	response := ListResponse{
		Flights: []Airplane{},
	}

	return ctx.JSON(http.StatusOK, response)
}
