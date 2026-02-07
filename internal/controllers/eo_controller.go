package controllers

import (
	"message-pocket/internal/config"
	"message-pocket/internal/define/dtos"
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
	cfg *config.Config,
) *EOController {
	return &EOController{
		eoService: eoService,
		config:    cfg,
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

	// 调用服务处理事件
	if err := c.eoService.EOWebhookEventHandle(ctx, &req); err != nil {
		e.App.Logger().ErrorContext(ctx, "Failed to process EO event", "err", err)
		return e.JSON(500, utils.NewJsonResponseWithoutData(500, "Failed to process event"))
	}

	// 返回成功响应
	return e.JSON(200, utils.NewJsonResponseWithoutData(0, "Success"))
}
