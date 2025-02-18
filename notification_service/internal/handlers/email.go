package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"notification_service/internal/config"
	"notification_service/internal/email"
	"notification_service/internal/models"
)

/*
POST /api/send/email/addresses

# Endpoint to send email to specific addresses

@ Request body: models.EmailSendRequestToSpecificAddressesBody
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
		quit <- email.RabbitSendEmailRequest(email.EmailSend(reqBody), &cfg.Rabbit)
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
