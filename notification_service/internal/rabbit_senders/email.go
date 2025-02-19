package rabbit_senders

import (
	"notification_service/internal/config"
	"notification_service/protobuf"

	"google.golang.org/protobuf/proto"
)

type Email struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailSend struct {
	Email      Email    `json:"email"`
	Recipients []string `json:"recipients"`
}

func (email *EmailSend) ToProtobuf() *protobuf.EmailSend {
	return &protobuf.EmailSend{
		Content: &protobuf.EmailSend_EmailContent{
			Subject: email.Email.Subject,
			Body:    email.Email.Body,
		},
		Recipients: email.Recipients,
	}
}

func RabbitSendEmailRequest(email EmailSend, cfg *config.RabbitConfig) error {
	// Prepare message with protobuf
	msg := email.ToProtobuf()
	body, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return RabbitSend(cfg, body, cfg.EmailQueue)
}
