package worker

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Starts worker. Must be called in the new goroutine
func startWorker(workerFunction func(msgs <-chan amqp.Delivery, cfg any), channel *amqp.Channel, queueName string, cfg any, quit chan struct{}) {
	q, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, fmt.Sprintf("Failed to declare a queue %s", queueName))
	msgs, err := channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, fmt.Sprintf("Failed to start consuming from %s", queueName))
	workerFunction(msgs, cfg)
	quit <- struct{}{}
	fmt.Printf("Worker for queue %s stopped\n", queueName)
}
