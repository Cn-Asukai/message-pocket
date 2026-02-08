package main

import (
	"context"
	"log"
	"log/slog"
	"message-pocket/internal/cron"
	"message-pocket/internal/middlewares"
	"message-pocket/internal/repo"
	"os"
	"strings"

	"message-pocket/internal/config"
	"message-pocket/internal/controllers"
	"message-pocket/internal/services"

	_ "message-pocket/migrations"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/samber/do/v2"
)

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	// 从 context 中提取 traceID
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

	injector := Inject(app)

	// 定时任务初始化
	cron.Init(app, injector)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		apiGroup := se.Router.Group("api")
		{
			eoController := do.MustInvoke[*controllers.EOController](injector)
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

func Inject(app core.App) do.Injector {

	injector := do.New()

	// load config
	cfg := config.GetConfig()

	// controller
	do.Provide(injector, controllers.ProvideEOController)

	// service
	do.Provide(injector, services.ProvideEOService)
	do.Provide(injector, services.ProvideMessageBoxService)
	do.Provide(injector, services.ProvideNapCatService)

	// repo
	do.Provide(injector, repo.ProvideMessageBoxRepo)

	// other
	do.ProvideValue(injector, app.DB())
	do.ProvideValue(injector, cfg)

	return injector
}
