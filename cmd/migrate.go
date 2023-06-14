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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func doMigrate(isUpgrade bool) {
	conf, err := config.InitConfig(viper.ConfigFileUsed())
	if err != nil {
		log.Fatal(err)
		panic(err)
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
		panic(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	mig, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	if isUpgrade {
		mig.Up()
	} else {
		mig.Down()
	}
}

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate your database",
	Long:  `this command migrates all of your migration files`,
	Run: func(cmd *cobra.Command, args []string) {
		state, _ := cmd.Flags().GetString("state")
		if state == "up" {
			doMigrate(true)
		} else if state == "down" {
			doMigrate(false)
		} else {
			log.Fatal("Invalid state")
			panic("Invalid state")
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	serveCmd.Flags().String("state", "down", "write the state")
}
