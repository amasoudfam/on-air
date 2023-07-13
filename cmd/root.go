package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "on-air",
	Short: "A brief description of your application",
	Long:  "A longer description that spans multiple lines",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, Cobra!")
	},
}

var (
	configFlag string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "config.yaml", "config path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
