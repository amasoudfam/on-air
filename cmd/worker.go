/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"on-air/config"
	"on-air/databases"
	"on-air/models"
	"on-air/repository"
	"on-air/server/services"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	go Run(&cfg.Worker, context.Background(), db)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	s := <-quit
	log.Infof("Worker: os signal recieved: %s", s)
}

func Run(worker *config.Worker, ctx context.Context, db *gorm.DB) {
	var apiMock *services.APIMockClient
	ticker := time.NewTicker(worker.Interval)
	counter := 0
	for {
		var tickets, _ = repository.GetExpiredTickets(db)

		for _, ticket := range tickets {
			err := db.Transaction(func(tx *gorm.DB) error {

				var flight, err = repository.FindFlightById(tx, int(ticket.FlightID))
				if err != nil {
					return err
				}

				_, err = apiMock.Refund(flight.Number, 1)

				//err

				if err != nil {
					repository.ChangeTicketStatus(tx, ticket.ID, string(models.TicketExpired))
					repository.ChangePaymentStatus(tx, ticket.ID, string(models.PaymentExpired))
				}

				return nil
			})

			log.Fatal(err)
		}

		if worker.Iteration > 0 {
			counter++
			if counter >= worker.Iteration {
				break
			}
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			log.Info("done signal recieved")
			break
		}

	}

	log.Info("workwer finished succesfully")
}
