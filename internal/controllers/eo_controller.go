package controllers

import (
	"fmt"
	"log"

	"message-pocket/internal/config"
	"message-pocket/internal/dtos"
	"message-pocket/internal/services"

	"github.com/pocketbase/pocketbase/core"
)

// EOController EO 控制器
type EOController struct {
	napcatService *services.NapCatService
	config        *config.Config
}

// NewEOController 创建 EO 控制器实例
func NewEOController(napcatService *services.NapCatService, cfg *config.Config) *EOController {
	return &EOController{
		napcatService: napcatService,
		config:        cfg,
	}
}

// getMessageTypeLabel 根据事件类型获取中文标签
func (c *EOController) getMessageTypeLabel(eventType string) string {
	switch eventType {
	case "deployment.created":
		return "开始部署"
	default:
		return eventType
	}
}

// EOWebhookEvent 处理 EO Webhook 事件
func (c *EOController) EOWebhookEvent(e *core.RequestEvent) error {
	// 解析请求体
	var req dtos.EOEventRequest
	if err := e.BindBody(&req); err != nil {
		return err
	}
	log.Printf("Received EO event: %+v", req)

	// 获取消息类型标签
	messageTypeLabel := c.getMessageTypeLabel(req.EventType)

	// 发送群消息
	message := fmt.Sprintf("🚀 EO 有新事件: %s", messageTypeLabel)
	groupID := ""

	if err := c.napcatService.SendGroupMessage(groupID, message); err != nil {
		log.Printf("Failed to send group message: %v", err)
		return e.String(500, fmt.Sprintf(`{"status": 1, "message": "Failed to send notification: %v"}`, err))
	}

	log.Printf("Successfully sent notification for EO event: %s", req.EventType)

	// 返回成功响应
	return e.String(200, `{"status": 0, "message": "ok"}`)
}
