package rabbit_senders

import (
	"notification_service/internal/config"
	"notification_service/protobuf"

	"google.golang.org/protobuf/proto"
)

type TelegramSend struct {
	Content    string   `json:"content"`
	Recipients []string `json:"recipients"`
}

func (tgSend *TelegramSend) ToProtobuf() *protobuf.TelegramMessageSend {
	return &protobuf.TelegramMessageSend{
		Content:    tgSend.Content,
		Recipients: tgSend.Recipients,
	}
}

func RabbitSendTgRequest(tgSend TelegramSend, cfg *config.RabbitConfig) error {
	// Prepare message with protobuf
	msg := tgSend.ToProtobuf()
	body, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return RabbitSend(cfg, body, cfg.TgQueue)
}
