/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"on-air/config"
	"on-air/repository"
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

	if !cfg.Worker.Enabled {
		log.Info("Worker: is disabled")
		return
	}

	go Run(&cfg.Worker)

}

func Run(worker *config.Worker) {
	//ticker := time.NewTicker(worker.Interval)
	//counter := 0
	for {
		//TODO : function to get pending ticket

	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	s := <-quit

	Action(worker)

	log.Infof("Worker: os signal recieved: %s", s)
}

func Action(worker *config.Worker) {
	var db *gorm.DB
	repository.ExpiringTicket(db)
}
