package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"on-air/config"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Flight struct {
	DB            *gorm.DB
	FlightService *config.FlightService
}

type Airplane struct {
	Number   string
	Airplane string
	Airline  string
	Price    int

	Origin      string
	Destination string
	Capacity    int
	StartedAt   time.Time
	FinishedAt  time.Time
}

type ListRequest struct {
	Origin      string `json:"origin" validate:"required"`
	Destination string `json:"destination" validate:"required"`
	Date        string `json:"date" query:"date" validate:"required,datetime=2006-01-02"`
}

type ListResponse struct {
	Flights []Airplane `json:"flights"`
}

func (f *Flight) List(ctx echo.Context) error {
	var req ListRequest
	// FIXME: c.Bind(&req) does not work
	req.Origin = ctx.QueryParam("origin")
	req.Destination = ctx.QueryParam("destination")
	req.Date = ctx.QueryParam("date")

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	address := fmt.Sprintf("%s/%s", f.FlightService.Url, "flights")
	origin := req.Origin
	destination := req.Destination
	date := req.Date

	url := fmt.Sprintf("%s?origin=%s&destination=%s&date=%s", address, origin, destination, date)
	res, err := http.Get(url)
	fmt.Println(res.Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	var response ListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, response)
}
