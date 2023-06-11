/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"on-air/config"
	"on-air/db"
	"on-air/middleware"

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

	serveCmd.Flags().String("port", "", "Port number")
	_ = viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
}

func startServer(port string) {
	e := echo.New()

	cfg := config.GetConfig()
	if port == "" {
		port = cfg.Server.Port
	}

	// Get db instance
	db := db.DbManager()

	// Inject db by middleware
	e.Use(middleware.DbMiddleware(db))

	address := fmt.Sprintf(":%s", port)
	// Define your routes and handlers here
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	// start the server
	e.Start(address)
}
