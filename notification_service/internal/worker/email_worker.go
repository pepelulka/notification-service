package worker

import (
	"fmt"
	"net/smtp"
	"notification_service/internal/config"
	"notification_service/protobuf"

	rs "notification_service/internal/rabbit_senders"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func fromProtobufEmailSend(emailSend *protobuf.EmailSend) *rs.EmailSend {
	return &rs.EmailSend{
		Email: rs.Email{
			Subject: emailSend.Content.Subject,
			Body:    emailSend.Content.Body,
		},
		Recipients: emailSend.Recipients,
	}
}

func sendEmail(cfg *config.EmailSenderConfig, email rs.EmailSend) error {
	message := fmt.Sprintf("Subject: %s\n\n%s", email.Email.Subject, email.Email.Body)
	auth := smtp.PlainAuth("", cfg.SenderAddress, cfg.SenderPassword, cfg.SmtpHost)

	return smtp.SendMail(cfg.SmtpHost+":"+cfg.SmtpPort, auth, cfg.SenderAddress, email.Recipients, []byte(message))
}

// Worker function to process email send requests
func emailWorker(msgs <-chan amqp.Delivery, cfg any) {
	emailSenderConfig := cfg.(*config.EmailSenderConfig)
	for d := range msgs {
		var emailSendProtobuf protobuf.EmailSend
		err := proto.Unmarshal(d.Body, &emailSendProtobuf)
		if err != nil {
			fmt.Println("Failed to process email send request message: ", err)
		}
		emailSend := fromProtobufEmailSend(&emailSendProtobuf)
		err = sendEmail(emailSenderConfig, *emailSend)
		if err != nil {
			fmt.Println("Failed to send email: ", err)
		}
	}
}
