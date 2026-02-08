package repo

import (
	"context"
	"fmt"
	"message-pocket/internal/constants/message_box_enum"
	"message-pocket/internal/define/model"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/samber/do/v2"
)

type IMessageBoxRepo interface {
	Create(ctx context.Context, in CreateMessageIn) (*model.MessageBoxModel, error)
	ListFailedBefore(ctx context.Context, t time.Time) ([]*model.MessageBoxModel, error)
	UpdateByID(ctx context.Context, messageID int32, data map[string]any) error
}

type MessageBoxRepo struct {
	db dbx.Builder
}

func NewMessageBoxRepo(db dbx.Builder) *MessageBoxRepo {
	return &MessageBoxRepo{
		db: db,
	}
}

func ProvideMessageBoxRepo(i do.Injector) (*MessageBoxRepo, error) {
	db := do.MustInvoke[dbx.Builder](i)
	return NewMessageBoxRepo(db), nil
}

type CreateMessageIn struct {
	BizID           string `db:"biz_id"`
	Message         string
	SourceRequest   string
	SourceType      message_box_enum.SourceType
	DestinationType message_box_enum.DestinationType
}

func (m *MessageBoxRepo) Create(ctx context.Context, in CreateMessageIn) (*model.MessageBoxModel, error) {
	// 先创建 MessageBoxModel
	createdAt := time.Now().Unix()
	createdAtStr := fmt.Sprintf("%d", createdAt)

	messageBox := &model.MessageBoxModel{
		ID:              0, // 将在插入后更新
		Status:          1, // 默认状态为发送中(Pending)
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
		WithContext(ctx).
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

func (m *MessageBoxRepo) ListFailedBefore(ctx context.Context, t time.Time) ([]*model.MessageBoxModel, error) {
	messages := make([]*model.MessageBoxModel, 0)
	// 查询状态为发送中(Pending)且创建时间早于指定时间的消息
	if err := m.db.NewQuery(`
			SELECT
				id,
				status,
				message,
				source_request,
				source_type,
				destination_type,
				created_at,
				last_sent_at
			FROM message_box
			WHERE status = {:status}
			AND created_at < {:created_at}
		`).
		Bind(map[string]any{
			"status":     message_box_enum.Pending.Val(), // 发送中状态
			"created_at": t.Unix(),
		}).
		WithContext(ctx).
		All(&messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *MessageBoxRepo) UpdateByID(ctx context.Context, messageID int32, data map[string]any) error {
	// 如果是发送成功状态，更新最后发送时间
	_, err := m.db.Update("message_box", data, dbx.NewExp("id = {:id}", dbx.Params{"id": messageID})).
		WithContext(ctx).
		Execute()
	return err
}
