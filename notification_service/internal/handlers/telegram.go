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
POST /api/send/tg/usernames

# Endpoint to send tg message to specific users

@ Request body:

	models.TgSendRequestToSpecificUsernamesBody

@ Response body:

	JsonStatusResponse
*/
func PostSendTgToUsernames(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	var reqBody models.TgSendRequestToSpecificUsernamesBody
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
		quit <- rs.RabbitSendTgRequest(rs.TelegramSend(reqBody), &cfg.Rabbit)
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
POST /api/send/tg/groups

# Endpoint to send tg message to specific groups

@ Request body:

	models.TgSendRequestToSpecificGroupsBody

@ Response body:

	JsonStatusResponse
*/
func PostSendTgToGroups(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	var reqBody models.TgSendRequestToSpecificGroupsBody
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

	dbConn, err := db.CreatePostgresConnection(&cfg.Database)
	if err != nil {
		ResponseWithError(w, err, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Construct TelegramSend object:
	persons, err := db.GetPersonsByGroupsFilter(&dbConn, ctx, reqBody.GroupNames)
	if err != nil {
		ResponseWithError(w, err, http.StatusBadRequest)
		return
	}
	usernames := make([]string, 0)
	for _, person := range persons {
		if person.TelegramId.Raw.Valid {
			usernames = append(usernames, person.TelegramId.Raw.String)
		}
	}
	telegramSend := rs.TelegramSend{
		Content:    reqBody.Content,
		Recipients: usernames,
	}

	quit := make(chan error)
	go func() {
		quit <- rs.RabbitSendTgRequest(rs.TelegramSend(telegramSend), &cfg.Rabbit)
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
