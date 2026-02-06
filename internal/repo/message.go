package repo

import "message-pocket/internal/constants/message_box_enum"

type MessageBoxModel struct {
	ID              int    `json:"id"`
	Status          int    `json:"status"`
	Message         string `json:"message"`
	SourceRequest   string
	SourceType      message_box_enum.SourceType
	DestinationType int
	CreatedAt       string
	LastedSentAt    string
}
