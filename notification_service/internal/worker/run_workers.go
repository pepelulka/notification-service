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

// Runs all the workers concurrently. If one of the workers failed program shuts down
func WorkerRoutine(cfg *config.RabbitConfig, workerCfg *config.WorkerConfig) {
	conn := tryConnectForever(cfg)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	fmt.Println("Connection to RabbitMQ established. Starting up workers...")

	// Workers start
	quit := make(chan struct{})

	go startWorker(emailWorker, ch, cfg.EmailQueue, &workerCfg.EmailWorkerConfig, quit)
	go startWorker(tgWorker, ch, cfg.TgQueue, &workerCfg.TgSenderConfig, quit)

	fmt.Println("Workers are running.")

	<-quit

	fmt.Println("Shutting down...")
}
