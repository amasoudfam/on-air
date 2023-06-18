package models

import (
	"errors"
	"on-air/utils"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `gorm:"type:varchar(50)"`
	LastName  string `gorm:"type:varchar(50)"`
	Email     string `gorm:"type:varchar(50)"`
	// NationalCode string `gorm:"type:varchar(50);unique"`
	PhoneNumber string `gorm:"type:varchar(15)"`
	Password    string `gorm:"type:varchar(128)"`
	Tickets     []Ticket
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var dbUser User
	result := db.First(&dbUser, "email = ?", email)
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return &dbUser, nil
}

func RegisterUser(db *gorm.DB, email string, password string) (*User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := User{Email: email, Password: hashedPassword}
	result := db.Create(&user)
	if err = result.Error; err != nil {
		return nil, err
	}

	return &user, nil
}
