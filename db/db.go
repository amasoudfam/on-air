package db

import (
	"fmt"
	"on-air/config"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func Init() {
	config := config.GetConfig()

	conn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Database.Host,
		config.Database.Username,
		config.Database.Password,
		config.Database.DbName,
		config.Database.Port,
	)

	db, err = gorm.Open(postgres.Open(conn), &gorm.Config{})
}
func DbManager() *gorm.DB {
	return db
}
