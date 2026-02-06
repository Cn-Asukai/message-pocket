package main

import (
	"context"
	"log"
	"log/slog"
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
	napcatService := services.NewNapCatService(cfg)
	eoService := services.NewEOService()
	eoController := controllers.NewEOController(eoService, napcatService, cfg)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// 保留原有的 hello 路由
		se.Router.GET("/hello/{name}", func(e *core.RequestEvent) error {
			name := e.Request.PathValue("name")
			return e.String(200, "Hello "+name)
		})

		// 添加 EO Webhook 路由
		se.Router.POST("/eo/webhook", eoController.EOWebhookEvent)

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
