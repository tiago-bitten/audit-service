package auditlog

import "context"

type Filter struct {
	ItemID  string
	UserID  string
	GroupID string
	Limit   int64
	Offset  int64
}

type Repository interface {
	Save(ctx context.Context, a *AuditLog) error
	FindByProject(ctx context.Context, projectID string, filter Filter) ([]AuditLog, error)
}
