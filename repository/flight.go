package repository

import (
	"on-air/models"

	"gorm.io/gorm"
)

func AddFlight(
	db *gorm.DB,
	flightNumber string,
	origin string,
	destination string,
	airLine string,
	airPlane string) (*models.Flight, error) {

	var fromCity models.City
	var toCity models.City

	err := db.First(&fromCity, "Name = ?", origin).Error

	if err != nil {
		return nil, err
	}

	err = db.First(&toCity, "Name = ?", destination).Error

	if err != nil {
		return nil, err
	}

	flight := models.Flight{
		Number:     flightNumber,
		FromCityID: fromCity.ID,
		ToCityID:   toCity.ID,
		Airplane:   airPlane,
		Airline:    airLine,
	}

	result := db.Create(&flight)
	if err := result.Error; err != nil {
		return nil, err
	}

	return &flight, nil
}

func FindFlight(db *gorm.DB, flightNumber string) (*models.Flight, error) {
	var flight models.Flight

	err := db.First(&flight, "Number = ?", flightNumber).Error
	if err != nil {
		return &models.Flight{}, err
	}

	return &flight, nil
}

func FindFlightById(db *gorm.DB, id int) (*models.Flight, error) {
	var flight models.Flight

	err := db.First(&flight, "ID = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &flight, nil
}
