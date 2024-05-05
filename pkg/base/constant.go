package base

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionName string

const (
	Post       CollectionName = "post"
	Like       CollectionName = "like"
	Comment    CollectionName = "comment"
	Reply      CollectionName = "replyComment"
	Share      CollectionName = "share"
	Log        CollectionName = "log"
	Bookmark   CollectionName = "bookmark"
	Preference CollectionName = "preference"
)

type BaseRepo interface {
	DeleteManyByQuery(ctx context.Context, filter any) error
	DeleteOneById(ctx context.Context, id primitive.ObjectID) error
	DeleteOneByQuery(ctx context.Context, query any) error
	FindOneById(ctx context.Context, id primitive.ObjectID, data any) error
	InsertMany(ctx context.Context, data []any) (*mongo.InsertManyResult, error)
	Create(ctx context.Context, data any) (primitive.ObjectID, error)
	FindOneByQuery(ctx context.Context, query any, result any) error
	UpdateOneByQuery(ctx context.Context, id primitive.ObjectID, query any) (*mongo.UpdateResult, error)
	UpdateOne(ctx context.Context, filter, update any) (*mongo.UpdateResult, error)
	FindByQuery(ctx context.Context, query any) (*mongo.Cursor, error)
	BulkUpdate(ctx context.Context, updateModel []mongo.WriteModel) (*mongo.BulkWriteResult, error)
	Aggregations(ctx context.Context, aggregation any) (*mongo.Cursor, error)
	GetSession() (mongo.Session, error)
	UpdateMany(ctx context.Context, filter any, update any) (*mongo.UpdateResult, error)
	DeleteMany(ctx context.Context, filter any) (*mongo.DeleteResult, error)
}

type BaseRepoImpl struct {
	DB *mongo.Collection
}
