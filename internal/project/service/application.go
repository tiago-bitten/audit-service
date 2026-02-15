package service

import (
	"github.com/tiago-bitten/audit-service/internal/project/app"
	"github.com/tiago-bitten/audit-service/internal/project/app/command"
	"github.com/tiago-bitten/audit-service/internal/project/app/query"
	"github.com/tiago-bitten/audit-service/internal/project/infra/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewApplication(db *mongo.Database) app.Application {
	projectRepo := repository.NewMongoProjectRepository(db)

	return app.Application{
		Commands: app.Commands{
			CreateProject: command.NewCreateProjectHandler(projectRepo),
		},
		Queries: app.Queries{
			FindAllProjects: query.NewFindAllProjectsHandler(projectRepo),
		},
	}
}
