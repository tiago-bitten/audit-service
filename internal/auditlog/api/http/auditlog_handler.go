package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tiago-bitten/audit-service/internal/auditlog/app"
	"github.com/tiago-bitten/audit-service/internal/auditlog/app/command"
	"github.com/tiago-bitten/audit-service/internal/auditlog/domain/auditlog"
)

type AuditLogHandler struct {
	app app.Application
}

func NewAuditLogHandler(app app.Application) *AuditLogHandler {
	return &AuditLogHandler{app: app}
}

func (h *AuditLogHandler) CreateAuditLog(c *gin.Context) {
	var cmd command.CreateAuditLogCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectID := c.GetString("projectID")
	err := h.app.Commands.CreateAuditLog.Handle(c.Request.Context(), cmd, projectID)
	if err != nil {
		if errors.Is(err, command.ErrMessageRequired) ||
			errors.Is(err, command.ErrDateRequired) ||
			errors.Is(err, command.ErrDateNotUTC) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		slog.Error("failed to create audit log", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true})
}

func (h *AuditLogHandler) FindAuditLogs(c *gin.Context) {
	projectID := c.GetString("projectID")

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	filter := auditlog.Filter{
		ItemID:  c.Query("item_id"),
		UserID:  c.Query("user_id"),
		GroupID: c.Query("group_id"),
		Limit:   limit,
		Offset:  offset,
	}

	views, err := h.app.Queries.FindAuditLogs.Handle(c.Request.Context(), projectID, filter)
	if err != nil {
		slog.Error("failed to find audit logs", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": views})
}
