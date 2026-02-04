package main

import (
	"log"

	"message-pocket/internal/config"
	"message-pocket/internal/controllers"
	"message-pocket/internal/services"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	// 加载配置
	cfg := config.GetConfig()

	// 初始化服务
	napcatService := services.NewNapCatService(cfg)
	eoController := controllers.NewEOController(napcatService, cfg)

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
