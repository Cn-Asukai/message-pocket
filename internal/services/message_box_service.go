package services

import (
	"context"
	"fmt"
	"log/slog"
	"message-pocket/internal/config"
	"message-pocket/internal/constants/message_box_enum"
	"message-pocket/internal/define/model"
	"message-pocket/internal/repo"
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

// SaveAndSendMessage 保存并发送消息
func (s *MessageBoxService) SaveAndSendMessage(
	ctx context.Context,
	req SaveMessageRequest,
) (*model.MessageBoxModel, error) {
	// 保存消息到数据库
	messageBox, err := s.messageBoxRepo.CreateMessage(repo.CreateMessageIn{
		BizID:           req.BizID,
		Message:         req.Message,
		SourceRequest:   req.SourceRequest,
		SourceType:      req.SourceType,
		DestinationType: req.DestinationType,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to save message",
			"err", err,
			"biz_id", req.BizID)
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	slog.InfoContext(ctx, "Successfully saved message",
		"message_id", messageBox.ID,
		"biz_id", req.BizID)

	// 发送消息
	if err := s.SendMessage(ctx, messageBox); err != nil {
		// 注意：这里返回了messageBox，即使发送失败，消息也已经保存
		// 调用者可以根据需要处理发送失败的情况
		return messageBox, fmt.Errorf("message saved but failed to send: %w", err)
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
