/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"on-air/cmd"
	"on-air/config"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	// Use the configuration values
	// fmt.Printf("Database host: %s\n", cfg.DBHost)
	_ = cfg
	cmd.Execute()
}
