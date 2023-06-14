package model

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	UserID     uint
	UnitPrice  int
	Count      int
	FlightID   int
	Status     string
	CreatedAt  time.Time
	User       User
	Flight     Flight
	Passengers []Passenger `gorm:"many2many:ticket_passengers;"`
}
