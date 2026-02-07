package model

import "message-pocket/internal/constants/message_box_enum"

type MessageBoxModel struct {
	ID              int32                            `json:"id" db:"id"`
	Status          int32                            `json:"status" db:"status"`
	Message         string                           `json:"message" db:"message"`
	SourceRequest   string                           `json:"source_request" db:"source_request"`
	SourceType      message_box_enum.SourceType      `json:"source_type" db:"source_type"`
	DestinationType message_box_enum.DestinationType `json:"destination_type" db:"destination_type"`
	CreatedAt       string                           `json:"created_at" db:"created_at"`
	LastedSentAt    string                           `json:"lasted_sent_at" db:"lasted_sent_at"`
}
