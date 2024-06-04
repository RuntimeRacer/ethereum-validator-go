// Package constants
/*
Copyright © 2024 RuntimeRacer
*/
package constants

const (
	// Application Info
	AppName    = "RuntimeRacer's Ethereum Validator Service"
	AppVersion = "v0.1"

	// Generic constants
	EnvPrefix = "ETHVAL"
)

var (
	// Error Codes
	ErrInitFailed     = "initialization failed: %v"
	ErrShutdownFailed = "shutdown failed: %v"
	ErrMissingEnvVar  = "required env var not set: %v"
	ErrConfigValue    = "required config value not set or bad value: %v"
	ErrApiServerStart = "unable to start api server: %v"
	ErrApiServerStop  = "unable to stop api server: %v"
)
