package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(c *gin.Context) {
	slog.Info("Request:", slog.String("Method", c.Request.Method), slog.String("Host", c.Request.Host))
	c.Next()
}
