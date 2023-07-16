package cmd

import (
	"context"
	"fmt"
	"net/http"
	"on-air/config"
	"on-air/databases"
	"on-air/models"
	"on-air/repository"
	"on-air/server/services"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/eapache/go-resiliency/breaker"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go Run(cfg, ctx, db, &wg)

	waitForShutdownSignal()
	cancel() // Signal the worker to stop

	wg.Wait() // Wait for the worker to finish processing

	log.Info("Worker has stopped")
}

func Run(cfg *config.Config, ctx context.Context, db *gorm.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	apiMock := &services.APIMockClient{
		Client:  &http.Client{},
		Breaker: &breaker.Breaker{},
		BaseURL: cfg.Services.ApiMock.BaseURL,
		Timeout: cfg.Services.ApiMock.Timeout,
	}

	ticker := time.NewTicker(cfg.Worker.Interval)
	counter := 0
	for {
		select {
		case <-ticker.C:
			var tickets, err = repository.GetExpiredTickets(db)
			if err != nil {
				log.Errorf("worker: Failed to get expired tickets: %v", err)
				continue
			}

			for _, ticket := range tickets {
				err := processTicket(db, apiMock, ticket)
				if err != nil {
					log.Errorf("worker: Failed to process ticket: %v", err)
				}
			}

			if cfg.Worker.Iteration > 0 {
				counter++
				if counter >= cfg.Worker.Iteration {
					return
				}
			}

		case <-ctx.Done():
			log.Info("worker: Done signal received")
			return
		}
	}
}

func processTicket(db *gorm.DB, apiMock *services.APIMockClient, ticket models.Ticket) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		flight, err := repository.FindFlightById(tx, int(ticket.FlightID))
		if err != nil {
			return fmt.Errorf("worker: failed to find flight: %w", err)
		}

		refundResult, err := apiMock.Refund(flight.Number, ticket.Count)
		if err != nil {
			return fmt.Errorf("worker: failed to refund ticket: %w", err)
		}

		if refundResult {
			err = repository.ChangeTicketStatus(tx, ticket.ID, string(models.TicketExpired))
			if err != nil {
				return fmt.Errorf("worker: failed to change ticket status: %w", err)
			}

			err = repository.ChangePaymentStatus(tx, ticket.ID, string(models.PaymentExpired))
			if err != nil {
				return fmt.Errorf("worker: failed to change payment status: %w", err)
			}
		}

		return nil
	})

	return err
}

func waitForShutdownSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	log.Info("worker: Received termination signal")
}
