package repository

import (
	"errors"
	"on-air/models"
	"time"

	"gorm.io/gorm"
)

func ReserveTicket(db *gorm.DB, userID int, flightID int, unitPrice int, passengerIDs []int) (*models.Ticket, error) {
	var passengers []models.Passenger

	result := db.Model(&passengers).Where("ID IN ?", passengerIDs)
	if result.RowsAffected != int64(len(passengerIDs)) {
		return nil, errors.New("passengers not found")
	}

	//TODO: Third party call to mock server to reserve a flight

	ticket := models.Ticket{
		UserID:     uint(userID),
		UnitPrice:  unitPrice,
		FlightID:   flightID,
		Count:      len(passengerIDs),
		Passengers: passengers,
		Status:     "Reserved",
		CreatedAt:  time.Now(),
	}

	result = db.Create(&ticket)
	if err := result.Error; err != nil {
		return nil, err
	}

	return &ticket, nil
}

func ChangeTicketStatus(db *gorm.DB, ticketID uint, status string) error {
	var dbTicket models.Ticket

	result := db.First(&dbTicket, "ID = ?", ticketID)
	if result.RowsAffected == 0 {
		return errors.New("Ticket not found")
	}

	dbTicket.Status = status

	result = db.Save(dbTicket)
	if result.RowsAffected == 0 {
		return errors.New("Ticket not found")
	}

	return nil
}
