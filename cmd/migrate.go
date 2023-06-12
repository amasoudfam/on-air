package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"on-air/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func ExecuteMigrate(configPath string, isUpgrade bool) {
	conf, err := config.InitConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	mig, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	if isUpgrade {
		mig.Up()
	} else {
		mig.Down()

	}
}
