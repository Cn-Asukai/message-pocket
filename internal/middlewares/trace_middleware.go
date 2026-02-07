package middlewares

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

// TraceIDKey context中trace_id的键
const TraceIDKey = "trace_id"

// TraceMiddleware 生成trace_id并存入context的中间件
func TraceMiddleware() func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		// 生成trace_id
		traceID := generateTraceID()

		// 将trace_id存入context
		ctx := context.WithValue(e.Request.Context(), TraceIDKey, traceID)

		// 更新请求的context
		e.Request = e.Request.WithContext(ctx)

		// 记录日志（可选）
		e.App.Logger().DebugContext(ctx, "Trace middleware executed", "trace_id", traceID)

		return e.Next()
	}
}

// generateTraceID 生成唯一的trace_id
func generateTraceID() string {
	// 使用时间戳和随机数生成简单的trace_id
	timestamp := time.Now().UnixNano()
	random := rand.Intn(10000)
	return fmt.Sprintf("trace-%d-%04d", timestamp, random)
}

// GetTraceIDFromContext 从context中获取trace_id
func GetTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}
