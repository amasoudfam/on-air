package model

import (
	"gorm.io/gorm"
)

type Country struct {
	gorm.Model
	Name string `gorm:"type:varchar(50)"`
}
