package repository

import (
	"on-air/models"

	"gorm.io/gorm"
)

func CreatePassenger(db *gorm.DB, userID int, nationalCode, firstName, lastName, gender string) (*models.Passenger, error) {
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
	var passengers []models.Passenger
	if err := db.Find(&passengers, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &passengers, nil
}
