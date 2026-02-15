package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tiago-bitten/audit-service/internal/project/domain/project"
)

type ProjectAuthMiddleware struct {
	projectRepo project.Repository
}

func NewProjectAuthMiddleware(projectRepo project.Repository) *ProjectAuthMiddleware {
	return &ProjectAuthMiddleware{projectRepo: projectRepo}
}

func (m *ProjectAuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.GetHeader("X-Project-Id")
		if projectID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
			return
		}

		p, err := m.projectRepo.FindByProjectID(c.Request.Context(), projectID)
		if err != nil || p == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
			return
		}

		if p.Status != project.StatusActive {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "project inactive"})
			return
		}

		c.Set("projectID", p.ProjectID)
		c.Next()
	}
}
