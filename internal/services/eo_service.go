package services

import (
	"fmt"
	"message-pocket/internal/dtos"
)

type EOService struct {
}

func NewEOService() *EOService {
	return &EOService{}
}

func (s *EOService) TransformEventToMessage(event *dtos.EOEventRequest) string {
	// 获取消息类型标签
	messageTypeLabel := s.getMessageTypeLabel(event.EventType)

	// 发送群消息
	message := fmt.Sprintf("🚀 EO 有新事件: %s", messageTypeLabel)

	return message
}

// getMessageTypeLabel 根据事件类型获取中文标签
func (s *EOService) getMessageTypeLabel(eventType string) string {
	switch eventType {
	case "deployment.created":
		return "开始部署"
	default:
		return eventType
	}
}
