package repository

import (
	"context"

	"github.com/tiago-bitten/audit-service/internal/auditlog/domain/auditlog"
	sharedmongo "github.com/tiago-bitten/audit-service/internal/shared/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoAuditLogRepository struct {
	collection *mongo.Collection
}

func NewMongoAuditLogRepository(db *mongo.Database) *MongoAuditLogRepository {
	return &MongoAuditLogRepository{
		collection: db.Collection("auditlogs"),
	}
}

func (r *MongoAuditLogRepository) Save(ctx context.Context, a *auditlog.AuditLog) error {
	ctx, cancel := sharedmongo.WithTimeout(ctx)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, toModel(a))
	return err
}

func (r *MongoAuditLogRepository) FindByProject(ctx context.Context, projectID string, filter auditlog.Filter) ([]auditlog.AuditLog, error) {
	ctx, cancel := sharedmongo.WithTimeout(ctx)
	defer cancel()

	query := bson.M{"project_id": projectID}

	if filter.ItemID != "" {
		query["item_id"] = filter.ItemID
	}
	if filter.UserID != "" {
		query["user_id"] = filter.UserID
	}
	if filter.GroupID != "" {
		query["group_id"] = filter.GroupID
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "date", Value: -1}}).
		SetLimit(limit).
		SetSkip(filter.Offset)

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var models []auditLogModel
	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	result := make([]auditlog.AuditLog, len(models))
	for i, m := range models {
		result[i] = *toDomain(&m)
	}

	return result, nil
}
