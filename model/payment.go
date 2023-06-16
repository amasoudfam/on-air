package model

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	Amount   int
	Status   string `gorm:type:varchar(20)`
	TicketID uint
	PayedAt  time.Time
	Ticket   Ticket
}
