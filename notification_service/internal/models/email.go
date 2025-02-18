package models

import "notification_service/internal/email"

type EmailSendRequestToSpecificAddressesBody email.EmailSend

type EmailSendRequestToGroupBody struct {
	Email      email.Email `json:"email"`
	GroupNames []string    `json:"group_names"`
}
