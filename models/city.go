package models

import (
	"gorm.io/gorm"
)

type City struct {
	gorm.Model
	Name      string `gorm:"type:varchar(50)"`
	CountryID uint
	Country   Country `gorm:"foreignKey:CountryID"`
}
