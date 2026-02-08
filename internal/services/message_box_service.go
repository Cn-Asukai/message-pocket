package services

import (
	"context"
	"fmt"
	"log/slog"
	"message-pocket/internal/config"
	"message-pocket/internal/constants/message_box_enum"
	"message-pocket/internal/define/model"
	"message-pocket/internal/repo"
	"time"

	"github.com/samber/do/v2"
)

type MessageBoxService struct {
	napcatService  *NapCatService
	messageBoxRepo repo.IMessageBoxRepo
}

// SaveMessageRequest 保存消息的请求参数
type SaveMessageRequest struct {
	BizID           string
	Message         string
	SourceRequest   string
	SourceType      message_box_enum.SourceType
	DestinationType message_box_enum.DestinationType
}

func NewMessageBoxService(
	napcatService *NapCatService,
	messageBoxRepo repo.IMessageBoxRepo,
) *MessageBoxService {
	return &MessageBoxService{
		napcatService:  napcatService,
		messageBoxRepo: messageBoxRepo,
	}
}

func ProvideMessageBoxService(i do.Injector) (*MessageBoxService, error) {
	napCatService := do.MustInvoke[*NapCatService](i)
	messageBoxRepo := do.MustInvoke[repo.IMessageBoxRepo](i)
	return NewMessageBoxService(napCatService, messageBoxRepo), nil
}

// SaveAndSendMessage 保存并发送消息
func (s *MessageBoxService) SaveAndSendMessage(
	ctx context.Context,
	req SaveMessageRequest,
) (*model.MessageBoxModel, error) {
	// 保存消息到数据库
	createMessageIn := repo.CreateMessageIn{
		BizID:           req.BizID,
		Message:         req.Message,
		SourceRequest:   req.SourceRequest,
		SourceType:      req.SourceType,
		DestinationType: req.DestinationType,
	}
	messageBox, err := s.messageBoxRepo.Create(ctx, createMessageIn)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to save message",
			"err", err,
			"created_message_in", createMessageIn)
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	slog.InfoContext(ctx, "Successfully saved message",
		"message_id", messageBox.ID,
		"biz_id", req.BizID)

	// 发送消息
	if err := s.SendMessage(ctx, messageBox); err != nil {
		if err := s.messageSentFailureProcess(ctx, messageBox.ID, err); err != nil {
			slog.ErrorContext(ctx, "messageSentFailureProcess finished with error", err)
		}
		// 注意：这里返回了messageBox，即使发送失败，消息也已经保存
		return messageBox, fmt.Errorf("message saved but failed to send: %w", err)
	}

	err = s.messageSentSuccessProcess(ctx, messageBox.ID)
	if err != nil {
		return nil, fmt.Errorf("message sent success but change message status failed: %w", err)
	}

	return messageBox, nil
}

// SendMessage 发送消息，根据destination_type决定发送方式
func (s *MessageBoxService) SendMessage(ctx context.Context, messageBox *model.MessageBoxModel) error {
	// 根据目的地类型选择发送方式
	switch messageBox.DestinationType {
	case message_box_enum.DestinationQQGroup:
		return s.sendToQQGroup(ctx, messageBox)
	default:
		return fmt.Errorf("unsupported destination type: %v", messageBox.DestinationType)
	}
}

// sendToQQGroup 发送消息到QQ群
func (s *MessageBoxService) sendToQQGroup(ctx context.Context, messageBox *model.MessageBoxModel) error {
	groupID := config.GetConfig().NapCatConfig.GroupID

	if err := s.napcatService.SendGroupMessage(ctx, groupID, messageBox.Message); err != nil {
		slog.ErrorContext(ctx, "Failed to send message to QQ group",
			"err", err,
			"message_id", messageBox.ID,
			"group_id", groupID)
		return fmt.Errorf("failed to send message to QQ group: %w", err)
	}

	slog.InfoContext(ctx, "Successfully sent message to QQ group",
		"message_id", messageBox.ID,
		"group_id", groupID)
	return nil
}

// 状态为发送中，且创建时间超过一分钟的即为发送失败的消息
func (s *MessageBoxService) findFailedMessages(ctx context.Context) ([]*model.MessageBoxModel, error) {
	// 查找创建时间超过1分钟的发送中消息
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)
	return s.messageBoxRepo.ListFailedBefore(ctx, oneMinuteAgo)
}

// 修改消息状态为发送成功
func (s *MessageBoxService) messageSentSuccessProcess(ctx context.Context, messageID int32) error {
	return s.messageBoxRepo.UpdateByID(ctx, messageID, map[string]any{
		"status":       message_box_enum.Sent,
		"last_sent_at": time.Now().Unix(),
	})
}

func (s *MessageBoxService) messageSentFailureProcess(ctx context.Context, messageID int32, err error) error {
	return s.messageBoxRepo.UpdateByID(ctx, messageID, map[string]any{
		"last_sent_at": time.Now().Unix(),
		"last_error":   err.Error(),
	})
}

func (s *MessageBoxService) MessageRetry(ctx context.Context) error {
	sentFailedMessages, err := s.findFailedMessages(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to find failed messages", "err", err)
		return fmt.Errorf("failed to find failed messages: %w", err)
	}

	if len(sentFailedMessages) == 0 {
		slog.InfoContext(ctx, "No failed messages to retry")
		return nil
	}

	slog.InfoContext(ctx, "Found failed messages to retry", "count", len(sentFailedMessages))

	for _, sentFailedMessage := range sentFailedMessages {
		if err = s.SendMessage(ctx, sentFailedMessage); err != nil {
			slog.ErrorContext(ctx, "Failed to resend message",
				"err", err,
				"message_id", sentFailedMessage.ID,
				"biz_id", sentFailedMessage.BizID)
			// 继续尝试其他消息，不立即返回错误
			continue
		}

		err = s.messageSentSuccessProcess(ctx, sentFailedMessage.ID)
		if err != nil {
			slog.ErrorContext(ctx, "Message resent success but change message status failed",
				"err", err,
				"message_id", sentFailedMessage.ID)
			// 继续处理其他消息
		}
		slog.InfoContext(ctx, "Successfully resent message",
			"message_id", sentFailedMessage.ID,
			"biz_id", sentFailedMessage.BizID)
	}

	return nil
}
