// Package apiserver
/*
Copyright © 2024 RuntimeRacer
*/
package apiserver

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"strings"
)

type EthereumValidatorHTTPSessionHandler struct {
	// Request Handling
	handlerId string
	server    *EthereumValidatorServer
	// Request Metadata
	originIP string
}

func (h *EthereumValidatorHTTPSessionHandler) init(server *EthereumValidatorServer, requestOrigin string) error {
	// Parse Request Origin
	originIP, _, err := net.SplitHostPort(requestOrigin)
	if err != nil {
		return fmt.Errorf("invalid request origin: %v", requestOrigin)
	}

	// Store initial Origin IP
	h.originIP = originIP

	// Store server pointer
	h.server = server
	return nil
}

func (h *EthereumValidatorHTTPSessionHandler) GetId() string {
	return h.handlerId
}

func (h *EthereumValidatorHTTPSessionHandler) ValidateRequest(request *http.Request) error {
	// Get Required Headers and request info
	// We check for a Remote IP and the API key
	apiKey := request.Header.Get("Validator-Api-Key")
	if len(request.RemoteAddr) == 0 || len(apiKey) == 0 {
		return errors.New("session init headers invalid")
	}
	requestAPIKey := viper.GetString("DEFAULT_API_KEY")
	if requestAPIKey != apiKey {
		return errors.New("invalid api key")
	}
	// Get origin IP; ensure port is being stripped
	requestOriginIP := request.RemoteAddr
	if strings.Contains(requestOriginIP, ":") {
		var errHostPort error
		requestOriginIP, _, errHostPort = net.SplitHostPort(request.RemoteAddr)
		if errHostPort != nil {
			return fmt.Errorf("invalid request origin: %v", request.RemoteAddr)
		}
	}
	// Ensure IP is valid -> This will be null if the parsing fails
	ipAddress := net.ParseIP(requestOriginIP)
	if ipAddress == nil {
		return fmt.Errorf("invalid request origin: %v", request.RemoteAddr)
	}
	if h.originIP != requestOriginIP {
		return errors.New("origin IP and request IP mismatch")
	}
	return nil
}

func (h *EthereumValidatorHTTPSessionHandler) finalize() {
	h.server.RemoveHTTPHandler(h.handlerId)
}
