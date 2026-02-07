package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"message-pocket/internal/constants/message_box_enum"
	"message-pocket/internal/define/dtos"
	"message-pocket/internal/services/logic"
)

type EOService struct {
	messageBoxService *MessageBoxService
}

func NewEOService(
	messageBoxService *MessageBoxService,
) *EOService {
	return &EOService{
		messageBoxService: messageBoxService,
	}
}

func (s *EOService) EOWebhookEventHandle(ctx context.Context, event *dtos.EOEventRequest) error {
	// 获取消息类型标签
	messageTypeLabel := logic.GetMessageTypeLabel(event.EventType)

	// 构建详细消息
	message := fmt.Sprintf(`🚀 EdgeOne 部署事件
📋 事件类型: %s
📁 项目名称: %s
🌿 代码分支: %s
🆔 项目ID: %s
🆔 部署ID: %s
⏰ 时间: %s`,
		messageTypeLabel,
		event.ProjectName,
		event.RepoBranch,
		event.ProjectID,
		event.DeploymentID,
		event.Timestamp,
	)

	requestStr, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event to json: %w", err)
	}

	// 使用 MessageBoxService 保存并发送消息
	_, err = s.messageBoxService.SaveAndSendMessage(ctx, SaveMessageRequest{
		BizID:           event.DeploymentID,
		Message:         message,
		SourceRequest:   string(requestStr),
		SourceType:      message_box_enum.SourceTypeEO,
		DestinationType: message_box_enum.DestinationQQGroup,
	})
	if err != nil {
		return fmt.Errorf("failed to save and send message: %w", err)
	}

	slog.InfoContext(ctx, "Successfully sent notification for EO event", "event_type", event.EventType)
	return nil
}
