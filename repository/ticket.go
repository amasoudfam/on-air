package repository

import (
	"fmt"
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
		fmt.Println(err)
	}

	return ticket, nil
}
