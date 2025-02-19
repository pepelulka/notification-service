package main

import (
	"log"
	"net/http"
	"notification_service/internal/config"
	"notification_service/internal/handlers"

	"github.com/gorilla/mux"
)

func setupRouter(r *mux.Router, cfg *config.Config) {
	// =================================================
	// Person API:

	r.HandleFunc(
		"/api/persons/all",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.GetPersons(cfg, w, r)
		},
	).Methods("GET")

	r.HandleFunc(
		"/api/persons/{personId}",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.GetPerson(cfg, w, r)
		},
	).Methods("GET")

	r.HandleFunc(
		"/api/persons/create",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.PostCreatePerson(cfg, w, r)
		},
	).Methods("POST")

	r.HandleFunc(
		"/api/persons/delete",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.DeletePersons(cfg, w, r)
		},
	).Methods("DELETE")

	// =================================================
	// Group API:

	r.HandleFunc(
		"/api/groups/all",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.GetAllGroups(cfg, w, r)
		},
	).Methods("GET")

	r.HandleFunc(
		"/api/groups/{groupName}",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.GetGroupMyName(cfg, w, r)
		},
	).Methods("GET")

	r.HandleFunc(
		"/api/groups/create",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.PostGroupCreate(cfg, w, r)
		},
	).Methods("POST")

	r.HandleFunc(
		"/api/groups/add_participants",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.PostGroupAddParticipants(cfg, w, r)
		},
	).Methods("POST")

	r.HandleFunc(
		"/api/groups/delete/{groupName}",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.DeleteGroup(cfg, w, r)
		},
	).Methods("DELETE")

	// =================================================
	// Sending API:

	r.HandleFunc(
		"/api/send/email/addresses",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.PostSendEmailToAddresses(cfg, w, r)
		},
	).Methods("POST")

	r.HandleFunc(
		"/api/send/email/groups",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.PostSendEmailToGroups(cfg, w, r)
		},
	).Methods("POST")

}

func main() {
	r := mux.NewRouter()

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatal(err)
		return
	}

	setupRouter(r, &cfg)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}
