package rabbitmq

import (
	"encoding/json"
	"log"
	"time"

	"github.com/SatvikG7/Scalable-Notification-System/internal/models"
	"github.com/SatvikG7/Scalable-Notification-System/internal/utils"
	"github.com/streadway/amqp"
	"golang.org/x/time/rate"
)

func Worker(queueName string, rateLimit *rate.Limiter, ch *amqp.Channel) {

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack, set to false for manual ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		if rateLimit.Allow() {
			if processNotification(d.Body) {
				d.Ack(false)
			} else {
				handleFailure(d.Body)
			}
		}
	}

}

func processNotification(notification []byte) bool {
	var UnmarshaledNotification models.Notification
	err := json.Unmarshal(notification, &UnmarshaledNotification)
	if err != nil {
		log.Println(err)
		return false
	}

	var user models.User

	user, err = utils.GetUser(UnmarshaledNotification.Recipient.String())
	if err != nil {
		log.Println(err)
		return false
	}

	time.Sleep(3 * time.Second)

	channel := user.Preference.GetChannel(UnmarshaledNotification.Priority)
	println("notifications_" + channel + "_" + UnmarshaledNotification.Priority)

	return true
}

func handleFailure(notification []byte) {
	// TODO: Log the failure to PostgreSQL after 3 retries
	log.Printf("Failed to process notification: %s", string(notification))
}

// Define a rate limiter with X requests per minute
func NewRateLimiter(requestsPerMinute int) *rate.Limiter {
	return rate.NewLimiter(rate.Every((time.Minute)/time.Duration(requestsPerMinute)), 1)
}
