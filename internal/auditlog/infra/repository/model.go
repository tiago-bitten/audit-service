package repository

import (
	"time"

	"github.com/tiago-bitten/audit-service/internal/auditlog/domain/auditlog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type auditLogModel struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty"`
	ItemID    string                 `bson:"item_id,omitempty"`
	UserID    string                 `bson:"user_id,omitempty"`
	Message   string                 `bson:"message"`
	GroupID   string                 `bson:"group_id,omitempty"`
	Object    map[string]interface{} `bson:"object,omitempty"`
	Date      time.Time              `bson:"date"`
	ProjectID string                 `bson:"project_id"`
}

func toModel(a *auditlog.AuditLog) *auditLogModel {
	return &auditLogModel{
		ItemID:    a.ItemID,
		UserID:    a.UserID,
		Message:   a.Message,
		GroupID:   a.GroupID,
		Object:    a.Object,
		Date:      a.Date,
		ProjectID: a.ProjectID,
	}
}

func toDomain(m *auditLogModel) *auditlog.AuditLog {
	return &auditlog.AuditLog{
		ID:        m.ID.Hex(),
		ItemID:    m.ItemID,
		UserID:    m.UserID,
		Message:   m.Message,
		GroupID:   m.GroupID,
		Object:    m.Object,
		Date:      m.Date,
		ProjectID: m.ProjectID,
	}
}
