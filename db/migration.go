package db

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

func getMigrateObject() *migrate.Migrate {
	conf := config.GetConfig()
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DbName)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func Up() {
	getMigrateObject().Up()
}

func Down() {
	getMigrateObject().Down()
}
