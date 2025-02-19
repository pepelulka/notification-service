package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"notification_service/internal/config"
	"notification_service/internal/db"
	"notification_service/internal/models"
	"strconv"

	"github.com/gorilla/mux"
)

/*
GET /api/persons/all

# Endpoint to get list of all persons

@ Response body: []models.Person
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
GET /api/persons/{personId}

# Endpoint to get info about user by id

@ Response body:

	status  = 200: models.Person
	status != 200: JsonStatusResponse
*/
func GetPerson(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), ENDPOINT_TIMEOUT)
	defer cancel()

	dbConn, err := db.CreatePostgresConnection(&cfg.Database)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	personId, err := strconv.Atoi(mux.Vars(r)["personId"])
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	person, err := db.GetPerson(&dbConn, ctx, personId)
	if err != nil {
		ResponseWithError(w, errors.New("not found"), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(person)
}

/*
POST /api/persons/create

# Endpoint to create new user

@ Request body: models.PersonCreate
@ Response body: models.PersonCreateResult
*/
func PostCreatePerson(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), ENDPOINT_TIMEOUT)
	defer cancel()

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

	dbConn, err := db.CreatePostgresConnection(&cfg.Database)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	personId, err := db.CreatePerson(reqBody, &dbConn, ctx)
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	response := models.PersonCreateResult{
		PersonId: personId,
	}
	json.NewEncoder(w).Encode(&response)
}

/*
DELETE /api/persons/delete

# Endpoint to delete users

@ Request body: models.PersonsDelete
@ Response body: JsonStatusResponse
*/
func DeletePersons(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), ENDPOINT_TIMEOUT)
	defer cancel()

	var reqBody models.PersonsDelete
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	conn, err := db.CreatePostgresConnection(&cfg.Database)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}

	if err := db.DeletePersons(&conn, ctx, reqBody.PersonIds); err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	ResponseWithOk(w)
}
