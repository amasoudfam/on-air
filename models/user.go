package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `gorm:"type:varchar(50)"`
	LastName  string `gorm:"type:varchar(50)"`
	Email     string `gorm:"type:varchar(50)"`
	// NationalCode string `gorm:"type:varchar(50);unique"`
	PhoneNumber string `gorm:"type:varchar(15)"`
	Password    string `gorm:"type:varchar(128)"`
	Tickets     []Ticket
}
