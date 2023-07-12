package cmd

import (
	"log"
	"on-air/config"
	"on-air/databases"
	"on-air/models"
	"on-air/utils"
	"time"

	"github.com/spf13/cobra"
	"gorm.io/datatypes"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed database",
	Long:  "this command seeds your database",
	Run: func(cmd *cobra.Command, args []string) {
		fakeFlag, _ := cmd.Flags().GetBool("fake")
		seed(configFlag, fakeFlag)
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
	seedCmd.Flags().Bool("fake", false, "If you want fill tables for test add this flag")
}

func seed(configPath string, fakeFlag bool) error {
	cfg, err := config.InitConfig(configPath)
	if err != nil {
		panic(err)
	}
	password, _ := utils.HashPassword("12345678")
	db := databases.InitPostgres(cfg)
	user := models.User{
		FirstName:   "user",
		LastName:    "test",
		Email:       "test@example.com",
		PhoneNumber: "09122222222",
		Password:    password,
	}

	err = db.FirstOrCreate(&user, models.User{Email: user.Email}).Error
	if err != nil {
		log.Fatal(err)
		return err
	}

	if fakeFlag {
		country := models.Country{
			Name: "Iran",
		}

		err = db.FirstOrCreate(&country, models.Country{Name: "Iran"}).Error
		if err != nil {
			log.Fatal(err)
			return err
		}

		cities := []models.City{
			{
				Name:      "Tehran",
				CountryID: country.ID,
			},
			{
				Name:      "Shiraz",
				CountryID: country.ID,
			},
			{
				Name:      "Esfahan",
				CountryID: country.ID,
			},
		}

		for i := range cities {
			err := db.FirstOrCreate(&cities[i], models.City{Name: cities[i].Name}).Error
			if err != nil {
				continue
			}
		}

		flight := models.Flight{
			Number:     "FL005",
			FromCityID: 1,
			ToCityID:   2,
			Airplane:   "Boeing 777",
			Airline:    "Singapore Airlines",
			StartedAt:  time.Now().Add(time.Hour * 74),
			FinishedAt: time.Now().Add(time.Hour * 80),
			Penalties: datatypes.JSON([]byte(`[{
					"Start":   "",
					"End":     "` + time.Now().Add(-48*time.Hour).Format(time.RFC3339) + `",
					"Percent": 20
				},
				{
					"Start":   "` + time.Now().Add(-48*time.Hour).Format(time.RFC3339) + `",
					"End":     "` + time.Now().Add(-24*time.Hour).Format(time.RFC3339) + `",
					"Percent": 20
				},
				{
					"Start":   "` + time.Now().Add(-24*time.Hour).Format(time.RFC3339) + `",
					"End":     "` + time.Now().Add(-1*time.Minute).Format(time.RFC3339) + `",
					"Percent": 40
				}]`)),
		}
		err = db.FirstOrCreate(&flight, models.Flight{Number: "FL005"}).Error
		if err != nil {
			log.Fatal(err)
			return err
		}

		passengers := []models.Passenger{
			{
				UserID:       user.ID,
				NationalCode: "2550000000",
				FirstName:    "Ghazanfar",
				LastName:     "Ghazanfari",
				Gender:       "male",
			},
			{
				UserID:       user.ID,
				NationalCode: "2550000001",
				FirstName:    "Mohammad",
				LastName:     "Mohammadi",
				Gender:       "male",
			},
		}

		for i := range passengers {
			err := db.FirstOrCreate(&passengers[i], models.Passenger{NationalCode: passengers[i].NationalCode}).Error
			if err != nil {
				continue
			}
		}

		ticket := models.Ticket{
			UserID:     user.ID,
			UnitPrice:  2100000,
			Count:      2,
			FlightID:   flight.ID,
			Status:     "complete",
			Flight:     flight,
			Passengers: passengers,
		}

		err := db.FirstOrCreate(&ticket, models.Ticket{UserID: user.ID, FlightID: flight.ID}).Error
		if err != nil {
			return err
		}

	}

	return nil
}
