package auditlog

import "time"

type AuditLog struct {
	ID        string
	ItemID    string
	UserID    string
	Message   string
	GroupID   string
	Object    map[string]interface{}
	Date      time.Time
	ProjectID string
}

func New(message, projectID string, itemID, userID, groupID string, object map[string]interface{}, date time.Time) *AuditLog {
	return &AuditLog{
		ItemID:    itemID,
		UserID:    userID,
		Message:   message,
		GroupID:   groupID,
		Object:    object,
		Date:      date,
		ProjectID: projectID,
	}
}
