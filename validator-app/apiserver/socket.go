// Package apiserver
/*
Copyright © 2024 RuntimeRacer
*/
package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

const (
	// Basic Connection Handling
	WEBSOCKET_NOT_SUPPORTED = "WEBSOCKET_NOT_SUPPORTED" // Default result if the client attempts to upgrade the connection to WebSocket Protocol
	INIT_SESSION_FAILED     = "INIT_SESSION_FAILED"     // Default result if the request could not properly initialize a HTTP Session
	HANDLE_REQUEST_FAILED   = "HANDLE_REQUEST_FAILED"   // Default result if the request could not properly handle a HTTP Request
	INVALID_ROUTE           = "INVALID_ROUTE"           // Default result if the route is not defined
	METHOD_NOT_ALLOWED      = "METHOD_NOT_ALLOWED"      // Default result if the method is not implemented
)

type ValidatorHttpError struct {
	Error  string          `json:"result"`
	Params json.RawMessage `json:"params,omitempty"`
}

type validatorServerRequestHandler struct {
	// Server reference
	server *EthereumValidatorServer
	router *chi.Mux
}

func (h *validatorServerRequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.ToLower(req.Header.Get("Upgrade")) == "websocket" &&
		strings.ToLower(req.Header.Get("Connection")) == "upgrade" {
		// Upgrade the HTTP connection to a WebSocket connection -> Can be added later if needed
		// v.UpgradeToWebSocket(w, req)
		errorMessage := fmt.Sprintf("websocket is currently not supported")
		log.Warning(errorMessage)
		// Bad Request
		w.WriteHeader(400)
		errorHTTPResponse(w, WEBSOCKET_NOT_SUPPORTED, errorMessage)
		return
	} else {
		// Handle normal HTTP Request
		var handler *EthereumValidatorHTTPSessionHandler

		// Check for session; create a new one if it doesn't exist
		// TODO: Super simple, do that better at a later point
		var okSession bool
		handler, okSession = server.activeHTTPSessions[req.Header.Get("Validator-Session-Id")]
		if !okSession {
			// Create new Validator Session
			var errSession error
			handler, errSession = h.InitializeHTTPSession(req)
			if errSession != nil {
				errorMessage := fmt.Sprintf("failed to initialize session: %v", errSession.Error())
				log.Warning(errorMessage)
				// Bad Request
				w.WriteHeader(400)
				errorHTTPResponse(w, INIT_SESSION_FAILED, errorMessage)
				return
			}
		}

		// Validate Request against Session Info - Very rudimentary IP check; should be done more safely
		if errValidate := h.ValidateRequestForHTTPSession(req, handler); errValidate != nil {
			errorMessage := fmt.Sprintf("failed to validate session: %v", errValidate.Error())
			log.Warning(errorMessage)
			// Unauthorized
			w.WriteHeader(401)
			errorHTTPResponse(w, HANDLE_REQUEST_FAILED, errorMessage)
			return
		}

		// Return Validator Session ID as part of the response header
		w.Header().Add("Validator-Session-Id", handler.GetId())

		// Handle the requests based on Path and Method
		h.router.ServeHTTP(w, req)
	}
}

func (h *validatorServerRequestHandler) ValidateRequestForHTTPSession(req *http.Request, handler *EthereumValidatorHTTPSessionHandler) error {
	// Get Required Headers
	requestOrigin := req.RemoteAddr
	apiKey := req.Header.Get("Validator-Api-Key")
	if len(requestOrigin) == 0 || len(apiKey) == 0 {
		return errors.New("session init headers invalid")
	}
	requestAPIKey := viper.GetString("DEFAULT_API_KEY")
	if requestAPIKey != apiKey {
		return errors.New("invalid api key")
	}
	if errOrigin := handler.ValidateOrigin(requestOrigin); errOrigin != nil {
		return errOrigin
	}
	return nil
}

func (h *validatorServerRequestHandler) InitializeHTTPSession(req *http.Request) (*EthereumValidatorHTTPSessionHandler, error) {
	// Get Required Headers
	requestOrigin := req.RemoteAddr
	apiKey := req.Header.Get("Validator-Api-Key")
	if len(requestOrigin) == 0 || len(apiKey) == 0 {
		return nil, errors.New("session init headers invalid")
	}

	// Verify API Key
	requestAPIKey := viper.GetString("DEFAULT_API_KEY")
	if requestAPIKey != apiKey {
		return nil, errors.New("invalid api key")
	}

	// Init Session handler
	sessionHandler := &EthereumValidatorHTTPSessionHandler{}
	if errInit := sessionHandler.init(h.server, req.RemoteAddr); errInit != nil {
		log.Errorf("failed to initialize session handler: %v", errInit.Error())
	}

	// Add to list of handlers
	h.server.AddHTTPHandler(sessionHandler)
	return sessionHandler, nil
}

func buildErrorHTTPResponse(errorType, errorMessage string) *ValidatorHttpError {
	messageJSON, errEncode := json.Marshal(errorMessage)
	if errEncode != nil {
		messageJSON = []byte("")
	}
	return &ValidatorHttpError{
		Error:  errorType,
		Params: messageJSON,
	}
}

func errorHTTPResponse(w http.ResponseWriter, errorType, errorMessage string) {
	if errEncode := json.NewEncoder(w).Encode(buildErrorHTTPResponse(errorType, errorMessage)); errEncode != nil {
		log.Error(fmt.Errorf("failed to encode data: %v", errEncode))
	}
}
