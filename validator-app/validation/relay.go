package validation

import (
	"fmt"
	"github.com/metachris/flashbotsrpc"
	"github.com/spf13/viper"
)

var mev2RelayBackendClient *flashbotsrpc.FlashbotsRPC

func getMEV2RelayBackendClient() (*flashbotsrpc.FlashbotsRPC, error) {
	// Return if already initialized
	if mev2RelayBackendClient != nil {
		return mev2RelayBackendClient, nil
	}

	// Build Endpoint URL
	rpcProviderURL := viper.GetString("BACKEND_ENDPOINT")
	rpcProviderToken := viper.GetString("BACKEND_TOKEN")
	rpcFullURL := fmt.Sprintf("%s/%s", rpcProviderURL, rpcProviderToken)

	mev2RelayBackendClient = flashbotsrpc.New(rpcFullURL)
	if mev2RelayBackendClient == nil {
		return nil, fmt.Errorf("failed to create RPC client for '%s'", rpcFullURL)
	}
	return mev2RelayBackendClient, nil
}
