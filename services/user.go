package services

import (
	"errors"
	"on-air/models"

	"gorm.io/gorm"
)

func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var dbUser models.User
	result := db.First(&dbUser, "email = ?", email)
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return &dbUser, nil

}
