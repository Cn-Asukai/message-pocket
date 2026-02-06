package controllers

import (
	"message-pocket/internal/config"
	"message-pocket/internal/dtos"
	"message-pocket/internal/services"
	"message-pocket/internal/utils"

	"github.com/pocketbase/pocketbase/core"
)

// EOController EO 控制器
type EOController struct {
	eoService     *services.EOService
	napcatService *services.NapCatService
	config        *config.Config
}

// NewEOController 创建 EO 控制器实例
func NewEOController(
	eoService *services.EOService,
	napcatService *services.NapCatService,
	cfg *config.Config,
) *EOController {
	return &EOController{
		eoService:     eoService,
		napcatService: napcatService,
		config:        cfg,
	}
}

// EOWebhookEvent 处理 EO Webhook 事件
func (c *EOController) EOWebhookEvent(e *core.RequestEvent) error {
	ctx := e.Request.Context()

	// 解析请求体
	var req dtos.EOEventRequest
	if err := e.BindBody(&req); err != nil {
		return err
	}
	e.App.Logger().InfoContext(ctx, "Received EO event", "request", req)

	groupID := config.GetConfig().NapCatConfig.GroupID

	message := c.eoService.TransformEventToMessage(&req)
	if err := c.napcatService.SendGroupMessage(ctx, groupID, message); err != nil {
		e.App.Logger().ErrorContext(ctx, "Failed to send group message", "err", err)
		return e.JSON(500, utils.NewJsonResponseWithoutData(1, "Failed to send group message"))
	}

	e.App.Logger().InfoContext(ctx, "Successfully sent notification for EO event", "event_type", req.EventType)

	// 返回成功响应
	return e.JSON(200, utils.NewJsonResponseWithoutData(0, "Success"))
}
