package models

import (
	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	UserID     uint
	UnitPrice  int
	Count      int
	FlightID   uint
	Status     string      `gorm:"type:varchar(10)"`
	User       User        `gorm:"foreignkey:UserID"`
	Flight     Flight      `gorm:"foreignkey:FlightID"`
	Passengers []Passenger `gorm:"many2many:ticket_passengers;"`
}

type TicketStatus string

const (
	Reserved      TicketStatus = "Reserved"
	TicketPaid    TicketStatus = "Paid"
	TicketExpired TicketStatus = "Expired"
)
