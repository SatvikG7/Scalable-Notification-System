package models

import "github.com/google/uuid"

type User struct {
	Id         uuid.UUID  `json:"id"`         // Unique identifier
	Username   string     `json:"username"`   // Username of the user
	Email      string     `json:"email"`      // Email of the user
	Phone      string     `json:"phone"`      // Phone number of the user
	Status     string     `json:"status"`     // Status of the user (active, inactive)
	Preference Preference `json:"preference"` // Notification preference of the user
}

// get channel
func (p *Preference) GetChannel(priority string) string {
	switch priority {
	case "low":
		return p.Low
	case "medium":
		return p.Medium
	case "high":
		return p.High
	default:
		return ""
	}
}
