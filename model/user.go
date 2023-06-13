package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string
	LastName     string
	Email        string `gorm:"unique;"`
	NationalCode string `gorm:"unique;"`
	PhoneNumber  string
	Password     string
	Tickets      []Ticket
}
