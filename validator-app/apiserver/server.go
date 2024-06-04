// Package apiserver
/*
Copyright © 2024 RuntimeRacer
*/
package apiserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/runtimeracer/ethereum-validator-go/constants"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	// Program flow
	exitSignalReceived = make(chan bool)
	shutdownComplete   = make(chan bool)

	// Program instance
	server *EthereumValidatorServer
)

type ConnectionHandler interface {
	GetId() string
	ValidateOrigin(requestOrigin string) error
}

type EthereumValidatorServer struct {
	// Application Flow
	isServingRequests bool
	// Params
	port         string
	inlineServer http.Server
	// Connections
	activeHTTPSessions map[string]*EthereumValidatorHTTPSessionHandler
	connMtx            sync.RWMutex
}

// Init Command executed
func Init(args []string) *EthereumValidatorServer {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
	})

	// Write logs to stdout and file as well if enabled
	writers := make([]io.Writer, 0)
	writers = append(writers, os.Stdout)
	if viper.GetBool("LOG_FILE") {
		currentTime := time.Now()
		fileName := fmt.Sprintf("ethereum-validator_%v.log", currentTime.Format("2006-01-02_15-04-05"))
		file, errLogFile := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if errLogFile != nil {
			log.Warnf("Unable to write to log file '%v'; please check disk space and folder permissions", fileName)
		} else {
			writers = append(writers, file)
		}
	}
	multiWriter := io.MultiWriter(writers...)
	log.SetOutput(multiWriter)

	// init EthereumValidatorServer
	var err error
	server, err = initEthereumValidatorServer()
	if err != nil {
		log.Fatalf(constants.ErrInitFailed, err.Error())
	}

	return server
}

func initEthereumValidatorServer() (*EthereumValidatorServer, error) {
	// Parse mandatory env vars
	servicePort := viper.GetInt("PORT")
	if servicePort < 1024 {
		return nil, fmt.Errorf(constants.ErrConfigValue, "Port")
	}

	// Initialize EthereumValidatorServer
	eventServer := &EthereumValidatorServer{
		port:               strconv.Itoa(servicePort),
		inlineServer:       http.Server{},
		activeHTTPSessions: make(map[string]*EthereumValidatorHTTPSessionHandler),
		connMtx:            sync.RWMutex{},
	}

	// Init Router
	router := GetApiRouter()
	AddCors(router)
	AddRoutes(router)

	// Init request handler & register with event bus
	requestHandler := &validatorServerRequestHandler{
		server: eventServer,
	}
	eventServer.inlineServer.Handler = requestHandler

	// Add shutdown handler for inline server
	eventServer.inlineServer.RegisterOnShutdown(server.OnShutdown)

	return eventServer, nil
}

func (e *EthereumValidatorServer) Start(ctx context.Context) {
	// Init shutdown Hook for Ctrl+C / Interrupt shutdown
	go shutdownHook()

	// Start Request Handling
	if errComms := server.OpenComms(); errComms != nil {
		log.Fatalf(constants.ErrApiServerStart, errComms.Error())
	}

	// Run till cancelled
	for {
		select {
		case <-exitSignalReceived:
			if server != nil {
				err := server.shutdown()
				if err != nil {
					log.Warnf(constants.ErrShutdownFailed, err)
				}

			}
			return
		}
	}
}

func (e *EthereumValidatorServer) Stop(ctx context.Context) {
	// Create channel to wait for shutdown
	shutdownComplete = make(chan bool)
	// Push to shutdown channel
	exitSignalReceived <- true
	// Wait for shutdown to complete before returning
	<-shutdownComplete
}

func (e *EthereumValidatorServer) OpenComms() error {
	if e.isServingRequests {
		return errors.New("api server already active")
	}

	go func(k *EthereumValidatorServer) {
		// Open Port
		listener, err := net.Listen("tcp", fmt.Sprintf(":%v", k.port))
		if err != nil {
			log.Error(fmt.Errorf("failed to listen: %v", err))
		}

		// Start serving with the API server
		log.Infof("Starting API server on %v", listener.Addr())
		if err = k.inlineServer.Serve(listener); err != nil {
			log.Error(fmt.Errorf("failed to serve: %v", err))
		}
		k.isServingRequests = true
	}(e)

	return nil
}

func (e *EthereumValidatorServer) AddHTTPHandler(h *EthereumValidatorHTTPSessionHandler) {
	defer e.connMtx.Unlock()
	e.connMtx.Lock()
	h.handlerId = uuid.New().String()
	e.activeHTTPSessions[h.handlerId] = h
	log.Infof("Added new http session handler '%v'", h.handlerId)
}

func (e *EthereumValidatorServer) OnShutdown() {
	log.Warnf("Gracefully shutting down %v...", constants.AppName)
	// Close all active connections
	defer e.connMtx.Unlock()
	e.connMtx.Lock()
	for name, _ := range e.activeHTTPSessions {
		e.activeHTTPSessions[name] = nil
	}
	fmt.Println("... done! Stopping process in 5 seconds...")
	time.Sleep(time.Second * 5)
}

func shutdownHook() {

	// Initially define termination signal channel
	shutdownSignalOS := make(chan os.Signal)
	signal.Notify(shutdownSignalOS, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	signalReceived := <-shutdownSignalOS
	log.Debugf("Caught signal: %+v", signalReceived)

	// Push to shutdown channel
	exitSignalReceived <- true
}

func (e *EthereumValidatorServer) shutdown() error {
	log.Info("Shutting down API server...")

	if e.isServingRequests {
		// Stop the API server
		errShutdown := e.inlineServer.Shutdown(context.Background())
		return fmt.Errorf(constants.ErrApiServerStop, errShutdown.Error())
	}

	log.Info("Shutdown complete.")

	return nil
}
