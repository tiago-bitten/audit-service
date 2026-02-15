package query

import (
	"context"
	"time"

	"github.com/tiago-bitten/audit-service/internal/project/domain/project"
)

type projectView struct {
	Name      string    `json:"name"`
	ProjectID string    `json:"project_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func toView(p *project.Project) projectView {
	return projectView{
		Name:      p.Name,
		ProjectID: p.ProjectID,
		Status:    string(p.Status),
		CreatedAt: p.CreatedAt,
	}
}

type FindAllProjectsHandler struct {
	projectRepo project.Repository
}

func NewFindAllProjectsHandler(projectRepo project.Repository) *FindAllProjectsHandler {
	return &FindAllProjectsHandler{
		projectRepo: projectRepo,
	}
}

func (h *FindAllProjectsHandler) Handle(ctx context.Context) ([]projectView, error) {
	projects, err := h.projectRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	views := make([]projectView, len(projects))
	for i, p := range projects {
		views[i] = toView(&p)
	}

	return views, nil
}
