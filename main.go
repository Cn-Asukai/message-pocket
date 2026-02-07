package main

import (
	"context"
	"log"
	"log/slog"
	"message-pocket/internal/middlewares"
	"message-pocket/internal/repo"
	"os"
	"strings"

	"message-pocket/internal/config"
	"message-pocket/internal/controllers"
	"message-pocket/internal/services"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "message-pocket/migrations"
)

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	// 从 context 中提取 traceID（假设之前已存入）
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		r.AddAttrs(slog.String("trace_id", traceID))
	}
	return h.Handler.Handle(ctx, r)
}

func main() {
	slog.SetDefault(slog.New(&ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})}))

	app := pocketbase.New()

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	// 加载配置
	cfg := config.GetConfig()

	// 初始化服务
	messageBoxRepo := repo.NewMessageBoxRepo(app.DB())
	napcatService := services.NewNapCatService(cfg)
	messageBoxService := services.NewMessageBoxService(napcatService, messageBoxRepo)
	eoService := services.NewEOService(messageBoxService)
	eoController := controllers.NewEOController(eoService, cfg)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		apiGroup := se.Router.Group("api")
		{
			// 添加 Trace 中间件（最先执行）
			apiGroup.BindFunc(middlewares.TraceMiddleware())
			// 添加 Token 验证中间件
			apiGroup.BindFunc(middlewares.TokenAuthMiddleware())
			// 添加 EO Webhook 路由
			apiGroup.POST("/eo/webhook", eoController.EOWebhookEvent)
		}

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
