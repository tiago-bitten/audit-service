package command

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/tiago-bitten/audit-service/internal/auditlog/domain/auditlog"
)

var (
	ErrMessageRequired = errors.New("message is required")
	ErrDateRequired    = errors.New("date is required")
	ErrDateNotUTC      = errors.New("date must be in UTC (offset +00:00)")
)

type CreateAuditLogCommand struct {
	ItemID  string                 `json:"item_id"`
	UserID  string                 `json:"user_id"`
	Message string                 `json:"message"`
	GroupID string                 `json:"group_id"`
	Object  map[string]interface{} `json:"object"`
	Date    time.Time              `json:"date"`
}

func (c CreateAuditLogCommand) Validate() error {
	if strings.TrimSpace(c.Message) == "" {
		return ErrMessageRequired
	}

	if c.Date.IsZero() {
		return ErrDateRequired
	}

	if _, offset := c.Date.Zone(); offset != 0 {
		return ErrDateNotUTC
	}

	return nil
}

type CreateAuditLogHandler struct {
	auditLogRepo auditlog.Repository
}

func NewCreateAuditLogHandler(auditLogRepo auditlog.Repository) *CreateAuditLogHandler {
	return &CreateAuditLogHandler{
		auditLogRepo: auditLogRepo,
	}
}

func (h *CreateAuditLogHandler) Handle(ctx context.Context, cmd CreateAuditLogCommand, projectID string) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	a := auditlog.New(cmd.Message, projectID, cmd.ItemID, cmd.UserID, cmd.GroupID, cmd.Object, cmd.Date)
	return h.auditLogRepo.Save(ctx, a)
}
