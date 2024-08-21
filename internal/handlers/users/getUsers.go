package users

import (
	"net/http"

	"github.com/SatvikG7/Scalable-Notification-System/internal/db"
	"github.com/SatvikG7/Scalable-Notification-System/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	query := `
	SELECT id, username, email, phone, status, preference_low_channel, preference_medium_channel, preference_high_channel 
	FROM users;
	`
	rows, err := db.DB.Raw(query).Rows()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}
	defer rows.Close()
	
	var users []models.User
	for rows.Next() {
		var user models.User
		var low, medium, high string
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Phone, &user.Status, &low, &medium, &high)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse user data",
			})
		}
		user.Preference = models.Preference{
			Low:    low,
			Medium: medium,
			High:   high,
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error iterating over users",
		})
	}

	return c.JSON(users)
}
