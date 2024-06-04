// Package cmd
/*
Copyright © 2024 RuntimeRacer
*/
package cmd

import (
	"fmt"
	"github.com/runtimeracer/ethereum-validator-go/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ethereum-validator-go",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Check if a default command is defined
	// https://github.com/spf13/cobra/issues/823
	cmd, _, errArgs := rootCmd.Find(os.Args[1:])
	// default cmd if no cmd is given
	if errArgs != nil || cmd.Args == nil {
		args := append([]string{"launch"}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	if errExecute := rootCmd.Execute(); errExecute != nil {
		delayedShutdownWithExitCode(errExecute, "", 1, 10)
	}
}

func delayedShutdownWithExitCode(err error, solutionHint string, exitCode, shutdownSeconds int) {
	fmt.Println(err)
	if len(solutionHint) > 0 {
		fmt.Println(solutionHint)
	}
	fmt.Println(fmt.Sprintf("Shutting down in %v seconds...", shutdownSeconds))
	time.Sleep(time.Second * time.Duration(shutdownSeconds))
	os.Exit(exitCode)
}

func init() {
	// Disable mousetrap
	cobra.MousetrapHelpText = ""
	// REMARK: Uncomment this in case we require a config file later
	cobra.OnInitialize(initConfig)

	// Root CMD Persistend Flags
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.json)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Setup Enviroment variables
	viper.SetEnvPrefix(constants.EnvPrefix)
	viper.AutomaticEnv() // read in environment variables that match

	bindFlags(rootCmd)

	// Ensure default API Key is set
	defaultApiKey := viper.GetString("DEFAULT_API_KEY")
	if len(defaultApiKey) == 0 {
		fmt.Println("FATAL: No default API Key was provided.")
		solutionHint := "This application requires an API Key for security reasons. Please check the documentation for details."
		delayedShutdownWithExitCode(fmt.Errorf("unable to read config.json file (maybe it is missing): %v"), solutionHint, 1, 10)
	}
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			if err := viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", constants.EnvPrefix, envVarSuffix)); err != nil {
				fmt.Println(err)
			}
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			if err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
				fmt.Println(err)
			}
		}
	})
}
