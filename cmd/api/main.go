package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	auditlogHttp "github.com/tiago-bitten/audit-service/internal/auditlog/api/http"
	auditlogService "github.com/tiago-bitten/audit-service/internal/auditlog/service"
	projectHttp "github.com/tiago-bitten/audit-service/internal/project/api/http"
	projectRepository "github.com/tiago-bitten/audit-service/internal/project/infra/repository"
	projectService "github.com/tiago-bitten/audit-service/internal/project/service"
	"github.com/tiago-bitten/audit-service/internal/shared/config"
	"github.com/tiago-bitten/audit-service/internal/shared/http/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	appConfig := config.LoadConfig()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(appConfig.MongoURL))
	if err != nil {
		slog.Error("failed to connect to mongodb", "err", err)
		os.Exit(1)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		slog.Error("failed to ping mongodb", "err", err)
		os.Exit(1)
	}

	slog.Info("connected to mongodb")

	db := client.Database(appConfig.Database)

	auditlogApp := auditlogService.NewApplication(db)
	projectApp := projectService.NewApplication(db)
	projectRepo := projectRepository.NewMongoProjectRepository(db)

	projectAuthMiddleware := middleware.NewProjectAuthMiddleware(projectRepo)
	adminAuthMiddleware := middleware.NewAdminAuthMiddleware(appConfig.AdminAPIKey)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	protectedGroup := r.Group("/v1")
	protectedGroup.Use(projectAuthMiddleware.Handle())
	auditlogHttp.RegisterAuditLogRoutes(protectedGroup, auditlogApp)

	adminGroup := r.Group("/v1")
	adminGroup.Use(adminAuthMiddleware.Handle())
	projectHttp.RegisterProjectRoutes(adminGroup, projectApp)

	srv := &http.Server{
		Addr:    ":" + appConfig.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	slog.Info("server started", "port", appConfig.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
	}

	if err := client.Disconnect(ctx); err != nil {
		slog.Error("mongodb disconnect error", "err", err)
	}

	slog.Info("server stopped")
}
