package worker

import (
	"context"
	"fmt"
	"notification_service/internal/config"
	"notification_service/protobuf"
	"strconv"
	"time"

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

// Send messages to specific users using data from etcd
func sendMessage(botApi *tgbotapi.BotAPI, etcdClient *etcdClient, tgSend rs.TelegramSend) {
	availableChatIds := make(map[int64]struct{}, 0)
	for _, recipient := range tgSend.Recipients {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		chatId, has, err := etcdClient.Get(ctx, recipient)
		cancel()
		if err != nil {
		} else if has {
			chatIdInt, err := strconv.ParseInt(chatId, 10, 64)
			if err != nil {
				continue
			}
			availableChatIds[chatIdInt] = struct{}{}
		}
	}

	for chatId := range availableChatIds {
		msg := tgbotapi.NewMessage(chatId, tgSend.Content)
		botApi.Send(msg)
	}

}

// Wait for messages to store mapping from username to chat id and store it in etcd.
func processMessages(bot *tgbotapi.BotAPI, etcdClient *etcdClient, quit chan struct{}) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			userName := update.Message.From.UserName
			chatId := strconv.FormatInt(update.Message.Chat.ID, 10)
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			etcdClient.Set(ctx, userName, chatId)
			fmt.Printf("Message got, %s : %s\n", userName, chatId)
			cancel()
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

	// Create etcd client
	etcdClient, err := createEtcdClient(&tgSenderConfig.Etcd)
	if err != nil {
		return
	}
	defer etcdClient.Close()

	// Running first worker - process messages
	// to store mapping username -> chatId in etcd

	quit := make(chan struct{})
	go processMessages(bot, &etcdClient, quit)

	// Running second worker that process requests from rabbitmq
	for d := range msgs {
		var tgSendPb protobuf.TelegramMessageSend
		err := proto.Unmarshal(d.Body, &tgSendPb)
		if err != nil {
			fmt.Println("Failed to process email send request message: ", err)
		}
		tgSend := fromProtobufTgSend(&tgSendPb)
		sendMessage(bot, &etcdClient, *tgSend)
	}

	<-quit
}
