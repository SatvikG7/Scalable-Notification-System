package users

import (
	"github.com/SatvikG7/Scalable-Notification-System/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := utils.GetUser(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return c.JSON(user)
}
