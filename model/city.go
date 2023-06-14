package model

import (
	"gorm.io/gorm"
)

type City struct {
	gorm.Model
	Name      string
	CountryID uint
	Country   Country
}
