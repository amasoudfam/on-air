package repository

import (
	"on-air/models"

	"gorm.io/gorm"
)

func GetTicket(db *gorm.DB, userID int, ticketID int) (models.Ticket, error) {
	var ticket models.Ticket
	err := db.Model(&models.Ticket{}).
		Where("user_id = ? and id = ?", userID, ticketID).
		Preload("Passengers").
		Preload("Flight").
		Preload("Flight.FromCity.Country").
		Preload("Flight.ToCity.Country").
		Find(&ticket).Error
	if err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

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
