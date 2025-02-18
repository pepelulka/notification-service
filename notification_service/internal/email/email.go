package email

import (
	"fmt"
	"net/smtp"
	"notification_service/internal/config"
	"notification_service/protobuf"
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

func FromProtobufEmailSend(emailSend *protobuf.EmailSend) *EmailSend {
	return &EmailSend{
		Email: Email{
			Subject: emailSend.Content.Subject,
			Body:    emailSend.Content.Body,
		},
		Recipients: emailSend.Recipients,
	}
}

func SendEmail(cfg *config.EmailSenderConfig, email EmailSend) error {
	message := fmt.Sprintf("Subject: %s\n\n%s", email.Email.Subject, email.Email.Body)
	auth := smtp.PlainAuth("", cfg.SenderAddress, cfg.SenderPassword, cfg.SmtpHost)

	return smtp.SendMail(cfg.SmtpHost+":"+cfg.SmtpPort, auth, cfg.SenderAddress, email.Recipients, []byte(message))
}
