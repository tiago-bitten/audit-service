package repository

import (
	"time"

	"github.com/tiago-bitten/audit-service/internal/project/domain/project"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type projectModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	ProjectID string             `bson:"project_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Status    string             `bson:"status"`
}

func toModel(p *project.Project) *projectModel {
	return &projectModel{
		Name:      p.Name,
		ProjectID: p.ProjectID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Status:    string(p.Status),
	}
}

func toDomain(m *projectModel) *project.Project {
	return &project.Project{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		ProjectID: m.ProjectID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Status:    project.Status(m.Status),
	}
}
