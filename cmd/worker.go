/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"on-air/config"
	"on-air/databases"
	"on-air/repository"
	"on-air/server/services"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/spf13/cobra"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Worker to handle pending request",

	Run: func(cmd *cobra.Command, args []string) {
		SetupWorker(configFlag)
		fmt.Println("worker called")
	},
}

func SetupWorker(configPath string) {
	cfg, err := config.InitConfig(configPath)
	if err != nil {
		panic(err)
	}

	db := databases.InitPostgres(cfg)

	if !cfg.Worker.Enabled {
		log.Info("Worker: is disabled")
		return
	}

	go Run(&cfg.Worker, db)

}

func Run(worker *config.Worker, db *gorm.DB) {
	var apiMock *services.APIMockClient
	//ticker := time.NewTicker(worker.Interval)
	//counter := 0
	for {
		//TODO : function to get pending ticket
		var tickets, _ = repository.ExpiringTicket(db)

		for _, ticket := range tickets {
			db.Transaction(func(tx *gorm.DB) error {
				var flight, _ = repository.FindFlightById(tx, ticket.FlightID)
				result, _ := apiMock.Refund(flight.Number)

				if result == true {
					repository.ChangeTicketStatus(tx, ticket.ID, "Expired")
					repository.ChangePaymentStatus(tx, ticket.ID, "Expired")
				}

				return nil
			})
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	s := <-quit

	log.Infof("Worker: os signal recieved: %s", s)
}
