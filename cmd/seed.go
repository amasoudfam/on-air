/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"on-air/config"
	"on-air/databases"
	"on-air/models"
	"on-air/utils"

	"github.com/spf13/cobra"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed database",
	Long:  "this command seeds your database",
	Run: func(cmd *cobra.Command, args []string) {
		addUser(configFlag)
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}

func addUser(configPath string) error {
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

	return db.Create(&user).Error
}
