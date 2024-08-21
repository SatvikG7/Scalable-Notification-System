package users

import (
	"github.com/SatvikG7/Scalable-Notification-System/internal/db"
	"github.com/gofiber/fiber/v2"
)

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	query := `
	DELETE FROM users WHERE id = ?;
	`
	rows, err := db.DB.Raw(query, id).Rows()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error iterating over user",
		})
	}

	return c.JSON(
		fiber.Map{
			"message": "User deleted successfully",
		},
	)
}
