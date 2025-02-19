package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"notification_service/internal/config"
	"notification_service/internal/db"
	"notification_service/internal/models"
	rs "notification_service/internal/rabbit_senders"
)

/*
POST /api/send/email/addresses

# Endpoint to send email to specific addresses

@ Request body:

	models.EmailSendRequestToSpecificAddressesBody

@ Response body:

	JsonStatusResponse
*/
func PostSendEmailToAddresses(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	var reqBody models.EmailSendRequestToSpecificAddressesBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		ResponseWithError(
			w,
			err,
			http.StatusBadRequest,
		)
		return
	}

	// Send email request with timeout:
	ctx, cancel := context.WithTimeout(context.Background(), ENDPOINT_TIMEOUT)
	defer cancel()

	quit := make(chan error)
	go func() {
		quit <- rs.RabbitSendEmailRequest(rs.EmailSend(reqBody), &cfg.Rabbit)
	}()
	select {
	case err = <-quit:
		if err != nil {
			ResponseWithError(w, err, http.StatusInternalServerError)
			return
		}
	case <-ctx.Done():
		ResponseWithError(w, ErrTimeout, http.StatusInternalServerError)
		return
	}

	ResponseWithOk(w)
}

/*
POST /api/send/email/groups

# Endpoint to send email to specific groups

@ Request body:

	models.EmailSendRequestToGroupBody

@ Response body:

	JsonStatusResponse
*/
func PostSendEmailToGroups(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	var reqBody models.EmailSendRequestToGroupBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), ENDPOINT_TIMEOUT)
	defer cancel()

	dbConn, err := db.CreatePostgresConnection(&cfg.Database)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Construct email.EmailSend object:
	persons, err := db.GetPersonsByGroupsFilter(&dbConn, ctx, reqBody.GroupNames)
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}
	addresses := make([]string, 0)
	for _, person := range persons {
		if person.Email.Raw.Valid {
			addresses = append(addresses, person.Email.Raw.String)
		}
	}
	emailSend := rs.EmailSend{
		Email:      reqBody.Email,
		Recipients: addresses,
	}

	quit := make(chan error)
	go func() {
		quit <- rs.RabbitSendEmailRequest(rs.EmailSend(emailSend), &cfg.Rabbit)
	}()
	select {
	case err = <-quit:
		if err != nil {
			ResponseWithError(w, err, http.StatusInternalServerError)
			return
		}
	case <-ctx.Done():
		ResponseWithError(w, ErrTimeout, http.StatusInternalServerError)
		return
	}

	ResponseWithOk(w)
}
