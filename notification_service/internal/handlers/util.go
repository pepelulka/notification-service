package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type JsonStatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

const ENDPOINT_TIMEOUT = 1 * time.Second

var ErrTimeout error = errors.New("timeout exceeded")

func ResponseWithError(w http.ResponseWriter, e error, status int) {
	errorResponse := JsonStatusResponse{Status: "error", Message: e.Error()}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse)
}

func ResponseWithOk(w http.ResponseWriter) {
	okResponse := JsonStatusResponse{Status: "ok", Message: ""}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(okResponse)
}
