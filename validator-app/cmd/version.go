// Package cmd
/*
Copyright © 2024 RuntimeRacer
*/
package cmd

import (
	"fmt"
	"github.com/runtimeracer/ethereum-validator-go/constants"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Shows the version of %v", constants.AppName),
	Long:  fmt.Sprintf("Shows the version of %v", constants.AppName),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("%v %v", constants.AppName, constants.AppVersion))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
