package main

import (
	"fmt"

	"github.com/SatvikG7/Scalable-Notification-System/config"
	"github.com/SatvikG7/Scalable-Notification-System/internal/db"
	"github.com/SatvikG7/Scalable-Notification-System/internal/rabbitmq"
	"github.com/SatvikG7/Scalable-Notification-System/internal/server"
)

func main() {
	if err := config.ConfigENV(); err != nil {
		fmt.Println("Error loading environment variables")
		return
	}

	if err := db.Init(); err != nil {
		fmt.Println("Error initializing database")
		return
	}

	go server.Init()

	// Define worker pools
	pools := []*rabbitmq.WorkerPool{
		rabbitmq.NewWorkerPool("notifications_email_high", rabbitmq.NewRateLimiter(50), 3, 5),
		rabbitmq.NewWorkerPool("notifications_sms_high", rabbitmq.NewRateLimiter(10), 3, 5),
		rabbitmq.NewWorkerPool("notifications_push_high", rabbitmq.NewRateLimiter(100), 3, 5),
		rabbitmq.NewWorkerPool("notifications_email_medium", rabbitmq.NewRateLimiter(25), 2, 3),
		rabbitmq.NewWorkerPool("notifications_sms_medium", rabbitmq.NewRateLimiter(6), 2, 3),
		rabbitmq.NewWorkerPool("notifications_push_medium", rabbitmq.NewRateLimiter(75), 2, 3),
		rabbitmq.NewWorkerPool("notifications_email_low", rabbitmq.NewRateLimiter(10), 1, 2),
		rabbitmq.NewWorkerPool("notifications_sms_low", rabbitmq.NewRateLimiter(3), 1, 2),
		rabbitmq.NewWorkerPool("notifications_push_low", rabbitmq.NewRateLimiter(50), 1, 2),
	}

	conn, ch := rabbitmq.ConnectRabbitMQ(pools)
	defer conn.Close()
	defer ch.Close()

	rabbitmq.Scheduler(pools, ch)

	select {}
}
