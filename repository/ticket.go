package repository

import (
	"on-air/models"

	"gorm.io/gorm"
)

func GetUserTickets(db *gorm.DB, userID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := db.Model(&models.Ticket{}).
		Where("user_id = ?", userID).
		Preload("Flight.FromCity.Country").
		Preload("Flight.ToCity.Country").
		Preload("User").
		Preload("Passengers").
		Find(&tickets).Error
	if err != nil {
		return nil, err
	}

	return tickets, nil
}
