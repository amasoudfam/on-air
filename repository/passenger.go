package repository

import (
	"errors"
	"on-air/models"

	"gorm.io/gorm"
)

func GetPassengerByNationalCodeAndUserID(db *gorm.DB, email string, userID int) (*models.Passenger, error) {
	var dbPassenger models.Passenger
	result := db.First(&dbPassenger, "national_code = ? and user_id = ?", email, userID)
	if result.RowsAffected == 0 {
		return nil, errors.New("Not found")
	}

	return &dbPassenger, nil
}

func AddPassenger(db *gorm.DB, userID int, nationalCode, firstName, lastName, gender string) (*models.Passenger, error) {
	passenger := models.Passenger{
		UserID:       uint(userID),
		NationalCode: nationalCode,
		FirstName:    firstName,
		LastName:     lastName,
		Gender:       gender,
	}

	result := db.Create(&passenger)
	if err := result.Error; err != nil {
		return nil, err
	}

	return &passenger, nil
}

func GetPassengersByUserID(db *gorm.DB, userID int) (*[]models.Passenger, error) {
	var dbPassengers []models.Passenger
	result := db.Find(&dbPassengers, "user_id = ?", userID)
	if result.RowsAffected == 0 {
		return nil, errors.New("Not found")
	}

	return &dbPassengers, nil
}
