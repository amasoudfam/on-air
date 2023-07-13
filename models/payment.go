package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	Amount   int
	Status   string `gorm:"type:varchar(20)"`
	TicketID uint
	PayedAt  time.Time
	Ticket   Ticket
}

type PaymentStatus string

const (
	Requested      PaymentStatus = "Requested"
	PaymentPaid    PaymentStatus = "Paid"
	Verified       PaymentStatus = "Verified"
	PaymentExpired PaymentStatus = "Expired"
)
