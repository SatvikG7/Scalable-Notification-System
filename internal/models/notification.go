package models

import (
	"github.com/google/uuid"
)

type Notification struct {
	Id        uuid.UUID `json:"id"`           // Unique identifier
	Priority  string    `json:"priority"`     // Priority of the notification (low, medium, high)
	Recipient uuid.UUID `json:"recipient_id"` // Recipient of the notification (user)
	Message   Message   `json:"message"`      // Message to be sent
	Status    string    `json:"status"`       // Status of the notification (pending, sent, failed)
}
