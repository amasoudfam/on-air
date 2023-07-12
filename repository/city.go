package repository

import (
	"on-air/models"

	"gorm.io/gorm"
)

func FindCityByName(db *gorm.DB, Name string) (*models.City, error) {
	var city models.City
	err := db.Where("Name = ?", Name).First(&city).Error
	if err != nil {
		return nil, err
	}

	return &city, nil
}
