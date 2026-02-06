package middlewares

import (
	"message-pocket/internal/config"
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

func TokenAuthMiddleware() func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		authHeader := e.Request.Header.Get("Authorization")

		split := strings.Split(authHeader, " ")
		if len(split) != 2 {
			return e.UnauthorizedError("unauthorized", nil)
		}

		token := split[1]

		c := config.GetConfig()
		if c.ServerConfig.OpenToken != token {
			return e.UnauthorizedError("unauthorized", nil)
		}

		return e.Next()
	}
}
