package model

import (
	"time"

	"gorm.io/gorm"
)

type Flight struct {
	gorm.Model
	Number     string
	FromCityID uint
	ToCityID   uint
	Airplane   string
	Airline    string
	StartedAt  time.Time
	FinishedAt time.Time
	FromCity   City `gorm:"foreignKey:FromCityID"`
	ToCity     City `gorm:"foreignKey:ToCityID"`
}
