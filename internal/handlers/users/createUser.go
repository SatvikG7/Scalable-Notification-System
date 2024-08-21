package users

import (
	"github.com/SatvikG7/Scalable-Notification-System/internal/db"
	"github.com/SatvikG7/Scalable-Notification-System/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	user.Id = uuid.New()

	err := db.DB.Exec("INSERT INTO users (id, username, email, phone, status, preference_low_channel, preference_medium_channel, preference_high_channel) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", user.Id, user.Username, user.Email, user.Phone, user.Status, user.Preference.Low, user.Preference.Medium, user.Preference.High).Error

	if err != nil {
		println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}
