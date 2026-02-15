package command

import (
	"context"
	"errors"
	"strings"

	"github.com/tiago-bitten/audit-service/internal/project/domain/project"
)

var ErrNameRequired = errors.New("project name is required")

type CreateProjectCommand struct {
	Name string `json:"name"`
}

func (c CreateProjectCommand) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return ErrNameRequired
	}
	return nil
}

type CreateProjectHandler struct {
	projectRepo project.Repository
}

func NewCreateProjectHandler(projectRepo project.Repository) *CreateProjectHandler {
	return &CreateProjectHandler{
		projectRepo: projectRepo,
	}
}

func (h *CreateProjectHandler) Handle(ctx context.Context, cmd CreateProjectCommand) (string, error) {
	if err := cmd.Validate(); err != nil {
		return "", err
	}

	existing, err := h.projectRepo.FindByName(ctx, cmd.Name)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return "", project.ErrNameInUse
	}

	p := project.New(cmd.Name)
	if err := h.projectRepo.Save(ctx, p); err != nil {
		return "", err
	}

	return p.ProjectID, nil
}
