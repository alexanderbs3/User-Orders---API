package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// Logger retorna um middleware Gin que loga todas as requisições.
// Usamos zap (biblioteca de log estruturado — JSON) em vez do log padrão do Go
// porque logs estruturados são indexáveis em ferramentas como Datadog/Loki/ELK.
func Logger(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path

        // c.Next() executa o handler principal + próximos middlewares
        c.Next()

        duration := time.Since(start)

        logger.Info("request",
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("duration", duration),
            zap.String("ip", c.ClientIP()),
        )
    }
}