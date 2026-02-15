package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tiago-bitten/audit-service/internal/project/app"
	"github.com/tiago-bitten/audit-service/internal/project/app/command"
	"github.com/tiago-bitten/audit-service/internal/project/domain/project"
)

type ProjectHandler struct {
	app app.Application
}

func NewProjectHandler(app app.Application) *ProjectHandler {
	return &ProjectHandler{app: app}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var cmd command.CreateProjectCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectID, err := h.app.Commands.CreateProject.Handle(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, project.ErrNameInUse) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, command.ErrNameRequired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		slog.Error("failed to create project", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "project_id": projectID})
}

func (h *ProjectHandler) FindAllProjects(c *gin.Context) {
	views, err := h.app.Queries.FindAllProjects.Handle(c.Request.Context())
	if err != nil {
		slog.Error("failed to find projects", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": views})
}
