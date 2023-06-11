package config

import "gorm.io/gorm"

type Store struct {
	DB *gorm.DB
}
