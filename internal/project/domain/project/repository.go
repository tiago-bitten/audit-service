package project

import "context"

type Repository interface {
	Save(ctx context.Context, p *Project) error
	FindAll(ctx context.Context) ([]Project, error)
	FindByProjectID(ctx context.Context, projectID string) (*Project, error)
	FindByName(ctx context.Context, name string) (*Project, error)
}
