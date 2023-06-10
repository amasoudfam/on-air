/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"on-air/cmd"
	"on-air/config"

	"github.com/spf13/pflag"
)

func main() {
	configFile := pflag.String("config", "config.yaml", "Path to config file")
	pflag.Parse()

	err := config.InitConfig(*configFile)
	if err != nil {
		panic(err)
	}

	// Use the configuration values
	// fmt.Printf("Database host: %s\n", cfg.Database.Host)
	cmd.Execute()
}
