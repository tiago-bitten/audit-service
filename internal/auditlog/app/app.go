package app

import (
	"github.com/tiago-bitten/audit-service/internal/auditlog/app/command"
	"github.com/tiago-bitten/audit-service/internal/auditlog/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateAuditLog *command.CreateAuditLogHandler
}

type Queries struct {
	FindAuditLogs *query.FindAuditLogsHandler
}
