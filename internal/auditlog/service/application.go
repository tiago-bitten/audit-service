package service

import (
	"github.com/tiago-bitten/audit-service/internal/auditlog/app"
	"github.com/tiago-bitten/audit-service/internal/auditlog/app/command"
	"github.com/tiago-bitten/audit-service/internal/auditlog/app/query"
	"github.com/tiago-bitten/audit-service/internal/auditlog/infra/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewApplication(db *mongo.Database) app.Application {
	auditLogRepo := repository.NewMongoAuditLogRepository(db)

	return app.Application{
		Commands: app.Commands{
			CreateAuditLog: command.NewCreateAuditLogHandler(auditLogRepo),
		},
		Queries: app.Queries{
			FindAuditLogs: query.NewFindAuditLogsHandler(auditLogRepo),
		},
	}
}
