package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"notification_service/internal/config"
	"notification_service/internal/db"
	"notification_service/internal/models"

	"github.com/gorilla/mux"
)

/*
GET /api/groups/all

# Endpoint to get all group names

@ Response body:

	status = 200 :  []string
	status!= 200 :  JsonStatusResponse
*/
func GetAllGroups(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), ENDPOINT_TIMEOUT)
	defer cancel()

	dbConn, err := db.CreatePostgresConnection(&cfg.Database)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	result, err := db.GetAllGroupNames(&dbConn, ctx)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(result)
}

/*
GET /api/groups/{groupName}

# Endpoint to get group by name

@ Response body:

	status = 200 :  models.Group
	status!= 200 :  JsonStatusResponse
*/
func GetGroupMyName(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), ENDPOINT_TIMEOUT)
	defer cancel()

	dbConn, err := db.CreatePostgresConnection(&cfg.Database)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	groupName := mux.Vars(r)["groupName"]

	result, err := db.GetGroupByName(groupName, &dbConn, ctx)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

/*
POST /api/groups/create

# Endpoint to create group

@ Request body:

	models.GroupCreate

@ Response body:

	JsonStatusResponse
*/
func PostGroupCreate(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), ENDPOINT_TIMEOUT)
	defer cancel()

	dbConn, err := db.CreatePostgresConnection(&cfg.Database)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	var reqBody models.GroupCreate
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	err = db.CreateGroup(reqBody, &dbConn, ctx)
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	ResponseWithOk(w)
}

/*
POST /api/groups/add_participants

# Endpoint to add participants in group

@ Request body:

	models.GroupCreate

@ Response body:

	JsonStatusResponse
*/
func PostGroupAddParticipants(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), ENDPOINT_TIMEOUT)
	defer cancel()

	dbConn, err := db.CreatePostgresConnection(&cfg.Database)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	var reqBody models.GroupCreate
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	err = db.AddPersonsToGroup(reqBody.Name, reqBody.ParticipantIds, &dbConn, ctx)
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	ResponseWithOk(w)
}

/*
DELETE /api/groups/delete/{groupName}

# Endpoint to delete group

@ Response body:

	JsonStatusResponse
*/
func DeleteGroup(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), ENDPOINT_TIMEOUT)
	defer cancel()

	dbConn, err := db.CreatePostgresConnection(&cfg.Database)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	groupName := mux.Vars(r)["groupName"]
	err = db.DeleteGroupByName(groupName, &dbConn, ctx)
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	ResponseWithOk(w)
}
