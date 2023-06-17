package models

import (
	"time"

	"gorm.io/gorm"
)

type Flight struct {
	gorm.Model
	Number     string `gorm:"type:varchar(20)"`
	FromCityID uint
	ToCityID   uint
	Airplane   string `gorm:"type:varchar(50)"`
	Airline    string `gorm:"type:varchar(50)"`
	StartedAt  time.Time
	FinishedAt time.Time
	FromCity   City `gorm:"foreignKey:FromCityID"`
	ToCity     City `gorm:"foreignKey:ToCityID"`
}
