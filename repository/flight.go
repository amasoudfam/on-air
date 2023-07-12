package repository

import (
	"on-air/models"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func AddFlight(
	db *gorm.DB,
	flightNumber string,
	origin string,
	destination string,
	airline string,
	airplane string,
	penalties datatypes.JSON,
	start time.Time,
	finish time.Time) (*models.Flight, error) {
	fromCity, err := FindCityByName(db, origin)

	if err != nil {
		return nil, err
	}

	toCity, err := FindCityByName(db, destination)
	if err != nil {
		return nil, err
	}

	flight := models.Flight{
		Number:     flightNumber,
		FromCityID: fromCity.ID,
		ToCityID:   toCity.ID,
		Airplane:   airplane,
		Airline:    airline,
		Penalties:  penalties,
		StartedAt:  start,
		FinishedAt: finish,
	}

	result := db.Create(&flight)
	if err := result.Error; err != nil {
		return nil, err
	}

	return &flight, nil
}

func FindFlight(db *gorm.DB, flightNumber string) (*models.Flight, error) {
	var flight models.Flight

	err := db.Where("Number = ?", flightNumber).First(&flight).Error
	if err != nil {
		return nil, err
	}

	return &flight, nil
}

func FindFlightById(db *gorm.DB, id int) (*models.Flight, error) {
	var flight models.Flight

	err := db.Where("ID = ?", id).First(&flight).Error

	if err != nil {
		return nil, err
	}

	return &flight, nil
}
