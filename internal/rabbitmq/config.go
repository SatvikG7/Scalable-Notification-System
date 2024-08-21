package rabbitmq

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func ConnectRabbitMQ(pools []*WorkerPool) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	for _, pool := range pools {
		_, err := ch.QueueDeclare(
			pool.QueueName, // name
			false,          // durable
			false,          // delete when unused
			false,          // exclusive
			false,          // no-wait
			nil,            // arguments
		)
		failOnError(err, "Failed to declare a queue")
	}

	return conn, ch
}
