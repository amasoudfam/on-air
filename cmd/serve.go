/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"on-air/config"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve command",
	Long:  "this command serve the project",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		startServer(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	cfg, _ := config.InitConfig()
	serveCmd.Flags().String("port", cfg.Server.Port, "Port number")
	_ = viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
}

func startServer(port string) {
	e := echo.New()

	// Define your routes and handlers here
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Start the server
	address := fmt.Sprintf(":%s", port)
	e.Start(address)
}
