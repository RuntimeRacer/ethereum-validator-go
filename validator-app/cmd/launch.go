// Package cmd
/*
Copyright © 2024 RuntimeRacer
*/
package cmd

import (
	"context"
	"github.com/runtimeracer/ethereum-validator-go/apiserver"
	"github.com/runtimeracer/ethereum-validator-go/constants"
	"github.com/spf13/cobra"
)

// init sets up the command
func init() {
	rootCmd.AddCommand(launchCmd)
}

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launches " + constants.AppName + " based on config and environment",
	Run: func(cmd *cobra.Command, args []string) {
		// Init Validator API Server
		server := apiserver.Init(args)
		server.Start(context.Background())
	},
}
