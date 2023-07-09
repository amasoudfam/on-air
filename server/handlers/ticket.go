package handlers

import (
	"net/http"
	"on-air/config"
	"on-air/models"
	"on-air/repository"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Ticket struct {
	DB  *gorm.DB
	JWT *config.JWT
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
		return ctx.JSON(http.StatusInternalServerError, err.Error())
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
