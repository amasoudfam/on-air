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
		configPath, _ := cmd.Flags().GetString("config")
		startServer(port, configPath)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("port", "", "Port number")
	serveCmd.Flags().String("config", "config.yaml", "config path")
	_ = viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
	_ = viper.BindPFlag("config", serveCmd.Flags().Lookup("config"))
}

func startServer(port string, configPath string) {
	cfg, err := config.InitConfig(configPath)
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	if port == "" {
		port = cfg.Server.Port
	}

	address := fmt.Sprintf(":%s", port)
	e.Start(address)
}
