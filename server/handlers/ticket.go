package handlers

import (
	"net/http"
	"on-air/config"
	"on-air/models"
	"on-air/repository"
	"on-air/server/services"
	"on-air/utils"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Ticket struct {
	DB            *gorm.DB
	JWT           *config.JWT
	APIMockClient *services.APIMockClient
}

type CountryResponse struct {
	Name string
}

type CityResponse struct {
	Name    string
	Country CountryResponse
}

type FlightResponse struct {
	Number     string
	Airplane   string
	Airline    string
	StartedAt  string
	FinishedAt string
	FromCity   CityResponse
	ToCity     CityResponse
}

type PassengerResponse struct {
	NationalCode string
	FirstName    string
	LastName     string
	Gender       string
}

type UserResponse struct {
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
}

type TicketResponse struct {
	ID         uint
	UnitPrice  int
	Count      int
	Status     string
	CreatedAt  string
	User       UserResponse
	Flight     FlightResponse
	Passengers []PassengerResponse
}

func (t *Ticket) GetTickets(ctx echo.Context) error {
	userID, _ := ctx.Get("user_id").(int)

	tickets, err := repository.GetUserTickets(t.DB, uint(userID))
	if err != nil {
		logrus.Error("ticket_handler: GetTickets failed when use repository.GetUserTickets, error:", err)
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	var ticketResponses []TicketResponse
	for _, ticket := range tickets {
		t := TicketResponse{
			ID:        ticket.ID,
			UnitPrice: ticket.UnitPrice,
			Count:     ticket.Count,
			Status:    ticket.Status,
			CreatedAt: ticket.CreatedAt.Format("2006-01-02 15:04"),
			User: UserResponse{
				FirstName:   ticket.User.FirstName,
				LastName:    ticket.User.LastName,
				Email:       ticket.User.Email,
				PhoneNumber: ticket.User.PhoneNumber,
			},
			Flight: FlightResponse{
				Number:     ticket.Flight.Number,
				Airplane:   ticket.Flight.Airplane,
				Airline:    ticket.Flight.Airline,
				StartedAt:  ticket.Flight.StartedAt.Format("2006-01-02 15:04"),
				FinishedAt: ticket.Flight.FinishedAt.Format("2006-01-02 15:04"),
				FromCity: CityResponse{
					Name: ticket.Flight.FromCity.Name,
					Country: CountryResponse{
						Name: ticket.Flight.FromCity.Country.Name,
					},
				},
				ToCity: CityResponse{
					Name: ticket.Flight.ToCity.Name,
					Country: CountryResponse{
						Name: ticket.Flight.ToCity.Country.Name,
					},
				},
			},
			Passengers: getPassengers(ticket.Passengers),
		}
		ticketResponses = append(ticketResponses, t)
	}

	return ctx.JSON(http.StatusOK, ticketResponses)
}

type ReserveRequest struct {
	FlightNumber string `json:"flight_number" binding:"required"`
	PassengerIDs []int  `json:"passengers" binding:"required"`
}

type ReserveResponse struct {
	TicketId int `json:"ticket_id" binding:"required"`
}

func (t *Ticket) Reserve(ctx echo.Context) error {
	userId, _ := ctx.Get("user_id").(int)
	var req ReserveRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	flightInfo, err := t.APIMockClient.GetFlight(req.FlightNumber)
	if err != nil {
		logrus.Error("ticket_handler: Reserve failed when use t.APIMockClient.GetFlight, error:", err)
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	flight, err := repository.FindFlight(t.DB, flightInfo.Number)
	if err != nil {
		logrus.Error("ticket_handler: Reserve failed when use repository.FindFlight, error:", err)
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	if flight == nil {
		flight, err = repository.AddFlight(t.DB,
			flightInfo.Number,
			flightInfo.Origin,
			flightInfo.Destination,
			flightInfo.Airplane,
			flightInfo.Airline,
			flightInfo.Penalties,
			flightInfo.StartedAt,
			flightInfo.FinishedAt,
		)
	}

	if err != nil {
		logrus.Error("ticket_handler: Reserve failed when use repository.AddFlight, error:", err)
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	flightReserve, err := t.APIMockClient.Reserve(req.FlightNumber, len(req.PassengerIDs))
	if err != nil {
		logrus.Error("ticket_handler: Reserve failed when use t.APIMockClient.Reserve, error:", err)
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	if !flightReserve {
		return ctx.JSON(http.StatusInternalServerError, "Sold out")
	}

	ticket, err := repository.ReserveTicket(
		t.DB,
		userId,
		int(flight.ID),
		flightInfo.Price,
		req.PassengerIDs,
	)
	if err != nil {
		t.APIMockClient.Refund(req.FlightNumber, len(req.PassengerIDs))
		logrus.Error("ticket_handler: Reserve failed when use repository.ReserveTicket, error:", err)
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return ctx.JSON(http.StatusOK, ReserveResponse{
		TicketId: int(ticket.ID),
	})
}

func getPassengers(passengers []models.Passenger) []PassengerResponse {
	var pass []PassengerResponse
	for _, passenger := range passengers {
		p := PassengerResponse{
			NationalCode: passenger.NationalCode,
			FirstName:    passenger.FirstName,
			LastName:     passenger.LastName,
			Gender:       passenger.Gender,
		}
		pass = append(pass, p)
	}

	return pass
}

func (t *Ticket) GetPDF(ctx echo.Context) error {
	userID, _ := strconv.Atoi(ctx.Get("id").(string))
	ticketID, err := strconv.Atoi(ctx.QueryParam("ticket_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid ticket_id")
	}

	ticket, err := repository.GetTicket(t.DB, userID, ticketID)
	if err != nil {
		logrus.Error("ticket_handler: GetPDF failed when use repository.GetTicket, error:", err)
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	result, err := utils.GeneratePDF(ticket)
	if err != nil {
		logrus.Error("ticket_handler: GetPDF failed when use utils.GeneratePDF, error:", err)
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	ctx.Response().Header().Set("Content-Type", "application/pdf")
	ctx.Response().Header().Set("Content-Disposition", "attachment; filename=myfile.pdf")
	ctx.Response().Header().Set("Content-Length", strconv.Itoa(len(result)))

	return ctx.Blob(http.StatusOK, "application/pdf", result)
}
