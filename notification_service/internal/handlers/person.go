package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"notification_service/internal/config"
	"notification_service/internal/db"
	"notification_service/internal/models"
)

/*
GET /api/persons/all

Endpoint to get list of all persons
*/
func GetPersons(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), ENDPOINT_TIMEOUT)
	defer cancel()

	dbConn, err := db.CreatePostgresConnection(&cfg.Database)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	persons, err := db.GetAllPersons(&dbConn, ctx)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(persons)
}

/*
POST /api/persons/create

# Endpoint to create new user

@ Request body: models.PersonCreate
@ Response body: models.PersonCreateResult
*/
func PostCreatePerson(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	var reqBody models.PersonCreate
	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		ResponseWithError(
			w,
			err,
			http.StatusBadRequest,
		)
		return
	}

	/*
		TO BE CONTINUED...
	*/
}
