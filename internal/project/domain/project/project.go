package project

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type Project struct {
	ID        string
	Name      string
	ProjectID string
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    Status
}

func New(name string) *Project {
	now := time.Now().UTC()
	return &Project{
		Name:      name,
		ProjectID: uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
		Status:    StatusActive,
	}
}
