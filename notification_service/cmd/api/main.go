package main

import (
	"fmt"
	"log"
	"net/http"
	"notification_service/internal/config"
	"notification_service/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatal(err)
		return
	}

	r.HandleFunc("/test", TestHandler).Methods("GET")
	r.HandleFunc(
		"/api/persons/all",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.GetUsers(cfg, w, r)
		},
	).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}
