package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"notification_service/internal/config"
	"notification_service/internal/db"
	"time"
)

const DB_TIMEOUT = 1 * time.Second

type JsonErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func ResponseWithError(w http.ResponseWriter, e error, status int) {
	errorResponse := JsonErrorResponse{Status: "error", Message: e.Error()}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse)
}

func GetUsers(cfg config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)

	dbConn, err := db.CreatePostgresConnection(cfg)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		cancel()
		return
	}
	defer dbConn.Close()

	persons, err := db.GetAllPersons(&dbConn, ctx)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		cancel()
		return
	}

	cancel()

	json.NewEncoder(w).Encode(persons)
}
