package app

import (
	"github.com/tiago-bitten/audit-service/internal/project/app/command"
	"github.com/tiago-bitten/audit-service/internal/project/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateProject *command.CreateProjectHandler
}

type Queries struct {
	FindAllProjects *query.FindAllProjectsHandler
}
