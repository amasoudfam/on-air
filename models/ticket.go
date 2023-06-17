package models

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
	Status     string `gorm:"type:varchar(10)"`
	CreatedAt  time.Time
	User       User        `gorm:"foreignkey:UserID"`
	Flight     Flight      `gorm:"foreignkey:FlightID"`
	Passengers []Passenger `gorm:"many2many:ticket_passengers;"`
}
