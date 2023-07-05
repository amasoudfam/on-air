package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"on-air/config"
	"on-air/databases"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate your database",
	Long:  `this command migrates all of your migration files`,
	Run: func(cmd *cobra.Command, args []string) {
		state, _ := cmd.Flags().GetString("state")
		steps, _ := cmd.Flags().GetInt("steps")
		migrateDB(state, configFlag, steps)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().String("state", "down", "write the state")
	migrateCmd.Flags().Int("steps", 1, "write the steps that you need up or down")
}

func migrateDB(state string, configPath string, steps int) {
	conf, err := config.InitConfig(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	db := databases.InitPostgres(conf)
	sql, _ := db.DB()
	driver, err := postgres.WithInstance(sql, &postgres.Config{})
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

	switch state {
	case "up":
		err = mig.Up()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("migrate up has done")
	case "down":
		err = mig.Down()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("migrate down has done")
	case "drop":
		err = mig.Drop()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("migrate drop has done")
	case "steps":
		err = mig.Steps(steps)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("migration with steps has done")
	default:
		log.Fatal("nothing")
	}
}
