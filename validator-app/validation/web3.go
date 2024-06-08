package validation

import (
	"fmt"
	"github.com/chenzhijie/go-web3"
	"github.com/spf13/viper"
)

const (
	ErrSlotDoesNotExist = "slot does not exist"
	ErrSlotInFuture     = "slot is in the future"
)

var web3client *web3.Web3

// getWeb3BackendClient returns a ready-to-use web3 client based on this application's env vars
func getWeb3BackendClient() (*web3.Web3, error) {
	// Return if already initialized
	if web3client != nil {
		return web3client, nil
	}

	// Build Endpoint URL
	rpcProviderURL := viper.GetString("BACKEND_ENDPOINT")
	rpcProviderToken := viper.GetString("BACKEND_TOKEN")
	rpcFullURL := fmt.Sprintf("%s/%s", rpcProviderURL, rpcProviderToken)

	// Init client
	client, errInit := web3.NewWeb3(rpcFullURL)
	if errInit != nil {
		return nil, errInit
	}

	// Set singleton var
	web3client = client
	return web3client, nil
}
