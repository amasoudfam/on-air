package repository

import (
	"on-air/models"
	"time"

	"gorm.io/gorm"
)

func ReserveTicket(db *gorm.DB, userID int, flightID int, unitPrice int, passengerIDs []int) (*models.Ticket, error) {
	var passengers []models.Passenger

	result := db.Model(&passengers).Where("ID IN ?", passengerIDs)
	if result.RowsAffected == 0 || result.RowsAffected != int64(len(passengerIDs)) {
		return nil, result.Error
	}

	ticket := models.Ticket{
		UserID:     uint(userID),
		UnitPrice:  unitPrice,
		FlightID:   uint(flightID),
		Count:      len(passengerIDs),
		Passengers: passengers,
		Status:     string(models.Reserved),
	}

	result = db.Create(&ticket)
	if err := result.Error; err != nil {
		return nil, err
	}

	return &ticket, nil
}

func ChangeTicketStatus(db *gorm.DB, id uint, status string) error {
	var ticket models.Ticket

	err := db.First(&ticket, "ID = ?", id).Error
	if err != nil {
		return err
	}

	ticket.Status = status

	err = db.Save(ticket).Error
	if err != nil {
		return err
	}

	return nil
}

func GetExpiredTickets(db *gorm.DB) ([]models.Ticket, error) {
	var tickets []models.Ticket

	err := db.Model(&tickets).
		Where("Status = ? AND CreatedAt > ?", string(models.Reserved), time.Now().Add(-15*time.Minute)).Error

	if err != nil {
		return tickets, err
	}
	return tickets, nil
}

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
