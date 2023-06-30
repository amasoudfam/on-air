package repository

import (
	"errors"
	"on-air/models"

	"gorm.io/gorm"
)

func AddFlight(db *gorm.DB, flightNumber string) (*models.Flight, error) {

	//TODO: complete model
	flight := models.Flight{
		Number: flightNumber,
	}

	result := db.Create(&flight)
	if err := result.Error; err != nil {
		return nil, err
	}

	return &flight, nil
}

func FindFlight(db *gorm.DB, flightNumber string) (*models.Flight, error) {

	var flight models.Flight

	result := db.First(&flight, "Number = ?", flightNumber)
	if result.RowsAffected == 0 {
		return &models.Flight{}, errors.New("flight not found")
	}

	return &flight, nil
}
