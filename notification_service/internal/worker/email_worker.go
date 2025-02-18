package worker

import (
	"fmt"
	"notification_service/internal/config"
	"notification_service/internal/email"
	"notification_service/protobuf"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

// Worker function to process email send requests
func emailWorker(msgs <-chan amqp.Delivery, cfg any) {
	emailSenderConfig := cfg.(*config.EmailSenderConfig)
	for d := range msgs {
		var emailSendProtobuf protobuf.EmailSend
		err := proto.Unmarshal(d.Body, &emailSendProtobuf)
		if err != nil {
			fmt.Println("Failed to process email send request message: ", err)
		}
		emailSend := email.FromProtobufEmailSend(&emailSendProtobuf)
		err = email.SendEmail(emailSenderConfig, *emailSend)
		if err != nil {
			fmt.Println("Failed to send email: ", err)
		}
	}
}
