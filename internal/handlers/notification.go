package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/SatvikG7/Scalable-Notification-System/internal/models"
	"github.com/SatvikG7/Scalable-Notification-System/internal/rabbitmq"
	"github.com/SatvikG7/Scalable-Notification-System/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func validateNotification(notification *models.Notification) error {
	if notification.Recipient.String() == "" {
		return fmt.Errorf("recipient id cannot be empty")
	}

	if notification.Priority != "low" && notification.Priority != "medium" && notification.Priority != "high" {
		return fmt.Errorf("invalid notification priority")
	}

	return nil
}

func saveNotification(notification *models.Notification) error {
	conn, ch := rabbitmq.ConnectRabbitMQ([]*rabbitmq.WorkerPool{})
	defer conn.Close()
	defer ch.Close()

	notification_bytes, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	user, err := utils.GetUser(notification.Recipient.String())
	if err != nil {
		return err
	}

	channel := user.Preference.GetChannel(notification.Priority)

	err = rabbitmq.PublishNotification(ch, "notifications_"+channel+"_"+notification.Priority, string(notification_bytes))
	if err != nil {
		return err
	}
	return nil
}

func CreateNotification(c *fiber.Ctx) error {
	notification := new(models.Notification)
	if err := c.BodyParser(notification); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := validateNotification(notification); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	notification.Id = uuid.New()

	notification.Status = "pending"

	err := saveNotification(notification)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to save notification",
		})
	}

	return c.Status(201).JSON(notification)
}
