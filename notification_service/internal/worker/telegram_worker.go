package worker

import (
	"fmt"
	"log"
	"notification_service/internal/config"
	"notification_service/protobuf"

	rs "notification_service/internal/rabbit_senders"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/proto"

	amqp "github.com/rabbitmq/amqp091-go"
)

func fromProtobufTgSend(tgSend *protobuf.TelegramMessageSend) *rs.TelegramSend {
	return &rs.TelegramSend{
		Content:    tgSend.Content,
		Recipients: tgSend.Recipients,
	}
}

func sendMessage(botApi *tgbotapi.BotAPI, tgSend rs.TelegramSend) {

}

func processMessages(bot *tgbotapi.BotAPI, quit chan struct{}) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %d", update.Message.From.UserName, update.Message.Chat.ID)
		}
	}

	quit <- struct{}{}
}

// We have two workers here - one that store users chat id
// And other who send messages to these chats
func tgWorker(msgs <-chan amqp.Delivery, cfg any) {
	tgSenderConfig := cfg.(*config.TgSenderConfig)

	bot, err := tgbotapi.NewBotAPI(tgSenderConfig.Token)
	if err != nil {
		fmt.Println("Failed to connect to bot api ", err)
	}

	// Running first worker - process messages
	// to store mapping username -> chatId in etcd

	quit := make(chan struct{})
	go processMessages(bot, quit)

	for d := range msgs {
		var tgSendPb protobuf.TelegramMessageSend
		err := proto.Unmarshal(d.Body, &tgSendPb)
		if err != nil {
			fmt.Println("Failed to process email send request message: ", err)
		}
		tgSend := fromProtobufTgSend(&tgSendPb)
		sendMessage(bot, *tgSend)
	}

	<-quit
}
