package cmd

import (
	"log"
	"on-air/config"
	"on-air/databases"
	"on-air/server"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve command",
	Long:  "this command serve the project",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		startServer(port, configFlag)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("port", "", "Port number")
}

func startServer(port string, configPath string) {
	cfg, err := config.InitConfig(configPath)
	if err != nil {
		panic(err)
	}

	db := databases.InitPostgres(cfg)
	redis := databases.InitRedis(cfg)

	if port == "" {
		port = cfg.Server.Port
	}

	log.Fatal(server.SetupServer(cfg, db, redis, port))
}
