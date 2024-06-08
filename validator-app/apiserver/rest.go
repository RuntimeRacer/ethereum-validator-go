// Package apiserver
/*
Copyright © 2024 RuntimeRacer
*/
package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/runtimeracer/ethereum-validator-go/validation"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func GetApiRouter() *chi.Mux {
	router := chi.NewRouter()

	// Define basic Middleware stack
	router.Use(middleware.RequestID)
	// router.Use(middleware.RealIP) -> Flawed: https://github.com/go-chi/chi/issues/453
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Define this API to be a JSON API
	router.Use(render.SetContentType(render.ContentTypeJSON))

	return router
}

func AddCors(router *chi.Mux) *chi.Mux {
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Api-Key", "X-CSRF-Token", "X-Forwarded-For", "X-Real-IP"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	return router
}

func AddRoutes(router *chi.Mux) {
	// Blockreward Endpoint
	router.Route("/blockreward", func(r chi.Router) {
		r.Get("/{slot}", blockRewardGetSlot)
	})
	// Syncduties Endpoint
	router.Route("/syncduties", func(r chi.Router) {
		r.Get("/{slot}", syncDutiesGetSlot)
	})

	// Error 400 if Route is not found
	router.NotFound(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(400)
		errorHTTPResponse(w, INVALID_ROUTE, "")
	})
	// Error 405 if Method is not allowed
	router.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(405)
		errorHTTPResponse(w, METHOD_NOT_ALLOWED, "")
	})
}

func blockRewardGetSlot(w http.ResponseWriter, r *http.Request) {
	slot := chi.URLParam(r, "slot")
	slotNumber, errParseSlotNumber := strconv.ParseUint(slot, 10, 64)
	if errParseSlotNumber != nil {
		w.WriteHeader(404)
		errorHTTPResponse(w, NOT_FOUND, validation.ErrSlotDoesNotExist)
		return
	}

	slotDetails, errSlot := validation.GetBlockRewardSlot(slotNumber)
	if errSlot != nil {
		// Log error
		log.Errorf("failed to get slot reward details: %v", errSlot)
		//
		// return 500 to caller with generic info to avoid leaking backend data
		w.WriteHeader(500)
		errorHTTPResponse(w, INTERNAL_SERVER_ERROR, "")
		return
	}
	// 200 OK
	w.WriteHeader(200)
	// Return the slot details
	if errEncode := json.NewEncoder(w).Encode(slotDetails); errEncode != nil {
		log.Error(fmt.Errorf("failed to encode data: %v", errEncode))
	}
}

func syncDutiesGetSlot(w http.ResponseWriter, r *http.Request) {
	slot := chi.URLParam(r, "slot")
	slotNumber, errParseSlotNumber := strconv.ParseUint(slot, 10, 64)
	if errParseSlotNumber != nil {
		w.WriteHeader(404)
		errorHTTPResponse(w, NOT_FOUND, validation.ErrSlotDoesNotExist)
		return
	}

	syncDuties, errSlot := validation.GetSyncDuties(slotNumber)
	if errSlot != nil {
		// Log error
		log.Errorf("failed to get slot syncduties details: %v", errSlot)
		// return 500 to caller with generic info to avoid leaking backend data
		w.WriteHeader(500)
		errorHTTPResponse(w, INTERNAL_SERVER_ERROR, "")
		return
	}
	// 200 OK
	w.WriteHeader(200)
	// Return the slot details
	if errEncode := json.NewEncoder(w).Encode(syncDuties); errEncode != nil {
		log.Error(fmt.Errorf("failed to encode data: %v", errEncode))
	}
}
