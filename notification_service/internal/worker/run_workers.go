package worker

import (
	"fmt"
	"notification_service/internal/config"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Trying connect to rabbitmq forever with cooldown = 1s
func tryConnectForever(cfg *config.RabbitConfig) *amqp.Connection {
	connString := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	var conn *amqp.Connection
	var err error
	for {
		if conn, err = amqp.Dial(connString); err == nil {
			break
		}
		fmt.Printf("Can't connect to RabbitMQ: %s. Trying again...\n", err)
		time.Sleep(time.Second)
	}
	return conn
}

func WorkerRoutine(cfg *config.RabbitConfig, workerCfg *config.WorkerConfig) {
	conn := tryConnectForever(cfg)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	forever := make(chan struct{})

	fmt.Println("Connection to RabbitMQ established. Starting up workers...")

	// Workers start
	startWorker(emailWorker, ch, cfg.EmailQueue, &workerCfg.EmailWorkerConfig)
	startWorker(tgWorker, ch, cfg.TgQueue, &workerCfg.TgSenderConfig)

	fmt.Println("Workers are running.")

	<-forever

	fmt.Println("Shutting down...")
}
