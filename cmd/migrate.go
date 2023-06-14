package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"on-air/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate your database",
	Long:  `this command migrates all of your migration files`,
	Run: func(cmd *cobra.Command, args []string) {
		state, _ := cmd.Flags().GetString("state")
		if state == "up" {
			migrateDB(true)
		} else if state == "down" {
			migrateDB(false)
		} else {
			log.Fatal("Invalid state")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	serveCmd.Flags().String("state", "down", "write the state")
}

func migrateDB(isUpgrade bool) {
	conf, err := config.InitConfig(viper.ConfigFileUsed())
	if err != nil {
		log.Fatal(err)
		return
	}

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DB)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
		return
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}

	mig, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
		return
	}

	if isUpgrade {
		mig.Up()
	} else {
		mig.Down()
	}
}
