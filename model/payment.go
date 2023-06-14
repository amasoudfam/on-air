package model

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	Amount   int
	Status   string
	TicketID uint
	PayedAt  time.Time
	Ticket   Ticket
}
