package models

import (
	"gorm.io/gorm"
)

type Passenger struct {
	gorm.Model
	UserID       uint
	NationalCode string `gorm:"type:varchar(10)"`
	FirstName    string `gorm:"type:varchar(50)"`
	LastName     string `gorm:"type:varchar(50)"`
	Gender       string `gorm:"type:varchar(5)"`
}
