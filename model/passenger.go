package model

import (
	"gorm.io/gorm"
)

type Passenger struct {
	gorm.Model
	UserID       uint
	NationalCode string
	FirstName    string
	LastName     string
	Gender       string
}
