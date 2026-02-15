package http

import (
	"github.com/gin-gonic/gin"
	"github.com/tiago-bitten/audit-service/internal/auditlog/app"
)

func RegisterAuditLogRoutes(g *gin.RouterGroup, app app.Application) {
	handler := NewAuditLogHandler(app)

	g.POST("/auditlogs", handler.CreateAuditLog)
	g.GET("/auditlogs", handler.FindAuditLogs)
}
