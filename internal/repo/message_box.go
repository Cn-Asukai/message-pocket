package repo

import (
	"fmt"
	"message-pocket/internal/constants/message_box_enum"
	"message-pocket/internal/define/model"
	"time"

	"github.com/pocketbase/dbx"
)

type IMessageBoxRepo interface {
	CreateMessage(in CreateMessageIn) (*model.MessageBoxModel, error)
}

type MessageBoxRepo struct {
	db dbx.Builder
}

func NewMessageBoxRepo(db dbx.Builder) *MessageBoxRepo {
	return &MessageBoxRepo{
		db: db,
	}
}

type CreateMessageIn struct {
	BizID           string `db:"biz_id"`
	Message         string
	SourceRequest   string
	SourceType      message_box_enum.SourceType
	DestinationType message_box_enum.DestinationType
}

func (m *MessageBoxRepo) CreateMessage(in CreateMessageIn) (*model.MessageBoxModel, error) {
	// 先创建 MessageBoxModel
	createdAt := time.Now().Unix()
	createdAtStr := fmt.Sprintf("%d", createdAt)
	
	messageBox := &model.MessageBoxModel{
		ID:              0, // 将在插入后更新
		Status:          1,
		Message:         in.Message,
		SourceRequest:   in.SourceRequest,
		SourceType:      in.SourceType,
		DestinationType: in.DestinationType,
		CreatedAt:       createdAtStr,
		LastedSentAt:    "",
	}

	// 使用 MessageBoxModel 的值构建 SQL
	result, err := m.db.NewQuery(`
		INSERT INTO message_box (
			biz_id,
			status,
			message,
			source_request,
			source_type,
			destination_type,
			created_at
		) VALUES (
			{:biz_id},
			{:status},
			{:message},
			{:source_request},
			{:source_type},
			{:destination_type},
			{:created_at}
		)
	`).
		Bind(map[string]any{
			"biz_id":           in.BizID,
			"status":           messageBox.Status,
			"message":          messageBox.Message,
			"source_request":   messageBox.SourceRequest,
			"source_type":      messageBox.SourceType.Val(),
			"destination_type": messageBox.DestinationType.Val(),
			"created_at":       createdAt,
		}).
		Execute()
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 更新 ID
	messageBox.ID = int32(id)
	return messageBox, nil
}
