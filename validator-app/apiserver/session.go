// Package apiserver
/*
Copyright © 2024 RuntimeRacer
*/
package apiserver

import (
	"errors"
	"fmt"
	"net"
)

type EthereumValidatorHTTPSessionHandler struct {
	// Request Handling
	handlerId string
	server    *EthereumValidatorServer
	// Request Validation
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

func (h *EthereumValidatorHTTPSessionHandler) ValidateOrigin(requestOrigin string) error {
	requestOriginIP, _, err := net.SplitHostPort(requestOrigin)
	if err != nil {
		return fmt.Errorf("invalid request origin: %v", requestOrigin)
	}
	if h.originIP != requestOriginIP {
		return errors.New("origin IP and request IP mismatch")
	}
	return nil
}
