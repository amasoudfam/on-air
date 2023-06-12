/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed database",
	Long:  "this command seeds your database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("seed called")
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
