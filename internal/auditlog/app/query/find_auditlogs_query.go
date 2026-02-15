package query

import (
	"context"
	"time"

	"github.com/tiago-bitten/audit-service/internal/auditlog/domain/auditlog"
)

type auditLogView struct {
	ID        string                 `json:"id"`
	ItemID    string                 `json:"item_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Message   string                 `json:"message"`
	GroupID   string                 `json:"group_id,omitempty"`
	Object    map[string]interface{} `json:"object,omitempty"`
	Date      time.Time              `json:"date"`
	ProjectID string                 `json:"project_id"`
}

func toView(a *auditlog.AuditLog) auditLogView {
	return auditLogView{
		ID:        a.ID,
		ItemID:    a.ItemID,
		UserID:    a.UserID,
		Message:   a.Message,
		GroupID:   a.GroupID,
		Object:    a.Object,
		Date:      a.Date,
		ProjectID: a.ProjectID,
	}
}

type FindAuditLogsHandler struct {
	auditLogRepo auditlog.Repository
}

func NewFindAuditLogsHandler(auditLogRepo auditlog.Repository) *FindAuditLogsHandler {
	return &FindAuditLogsHandler{
		auditLogRepo: auditLogRepo,
	}
}

func (h *FindAuditLogsHandler) Handle(ctx context.Context, projectID string, filter auditlog.Filter) ([]auditLogView, error) {
	logs, err := h.auditLogRepo.FindByProject(ctx, projectID, filter)
	if err != nil {
		return nil, err
	}

	views := make([]auditLogView, len(logs))
	for i, l := range logs {
		views[i] = toView(&l)
	}

	return views, nil
}
