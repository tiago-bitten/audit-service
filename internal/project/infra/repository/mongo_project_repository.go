package repository

import (
	"context"

	"github.com/tiago-bitten/audit-service/internal/project/domain/project"
	sharedmongo "github.com/tiago-bitten/audit-service/internal/shared/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoProjectRepository struct {
	collection *mongo.Collection
}

func NewMongoProjectRepository(db *mongo.Database) *MongoProjectRepository {
	return &MongoProjectRepository{
		collection: db.Collection("projects"),
	}
}

func (r *MongoProjectRepository) Save(ctx context.Context, p *project.Project) error {
	ctx, cancel := sharedmongo.WithTimeout(ctx)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, toModel(p))
	return err
}

func (r *MongoProjectRepository) FindAll(ctx context.Context) ([]project.Project, error) {
	ctx, cancel := sharedmongo.WithTimeout(ctx)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var models []projectModel
	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	result := make([]project.Project, len(models))
	for i, m := range models {
		result[i] = *toDomain(&m)
	}

	return result, nil
}

func (r *MongoProjectRepository) FindByProjectID(ctx context.Context, projectID string) (*project.Project, error) {
	ctx, cancel := sharedmongo.WithTimeout(ctx)
	defer cancel()

	var m projectModel
	found, err := sharedmongo.FindOneDocument(ctx, r.collection, bson.M{"project_id": projectID}, &m)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return toDomain(&m), nil
}

func (r *MongoProjectRepository) FindByName(ctx context.Context, name string) (*project.Project, error) {
	ctx, cancel := sharedmongo.WithTimeout(ctx)
	defer cancel()

	var m projectModel
	found, err := sharedmongo.FindOneDocument(ctx, r.collection, bson.M{"name": name}, &m)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return toDomain(&m), nil
}
