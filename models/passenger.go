package models

import (
	"gorm.io/gorm"
)

type Passenger struct {
	gorm.Model
	UserID       uint   `gorm:"uniqueIndex:idx_name_location"`
	NationalCode string `gorm:"type:varchar(10);uniqueIndex:idx_name_location"`
	FirstName    string `gorm:"type:varchar(50)"`
	LastName     string `gorm:"type:varchar(50)"`
	Gender       string `gorm:"type:varchar(5)"`
}
