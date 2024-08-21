package utils

import (
	"github.com/SatvikG7/Scalable-Notification-System/internal/db"
	"github.com/SatvikG7/Scalable-Notification-System/internal/models"
)

func GetUser(id string) (models.User, error) {
	var user models.User
	query := `
	SELECT id, username, email, phone, status, preference_low_channel, preference_medium_channel, preference_high_channel 
	FROM users WHERE id = ?;
	`
	rows, err := db.DB.Raw(query, id).Rows()
	if err != nil {
		return user, err
	}
	defer rows.Close()

	for rows.Next() {
		var low, medium, high string
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Phone, &user.Status, &low, &medium, &high)
		if err != nil {
			return user, err
		}
		user.Preference = models.Preference{
			Low:    low,
			Medium: medium,
			High:   high,
		}
	}

	if err = rows.Err(); err != nil {
		return user, err
	}

	return user, nil
}
