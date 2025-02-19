package models

import rs "notification_service/internal/rabbit_senders"

type EmailSendRequestToSpecificAddressesBody rs.EmailSend

type EmailSendRequestToGroupBody struct {
	Email      rs.Email `json:"email"`
	GroupNames []string `json:"group_names"`
}
