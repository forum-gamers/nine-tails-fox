package post

import (
	"context"

	b "github.com/forum-gamers/nine-tails-fox/pkg/base"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewPostRepo() PostRepo {
	return &PostRepoImpl{b.NewBaseRepo(b.GetCollection(b.Post))}
}

func (r *PostRepoImpl) Create(ctx context.Context, data *Post) error {
	result, err := r.BaseRepo.Create(ctx, data)
	if err != nil {
		return err
	}
	data.Id = result
	return nil
}

func (r *PostRepoImpl) FindById(ctx context.Context, id primitive.ObjectID, data *Post) error {
	return r.FindOneById(ctx, id, data)
}

func (r *PostRepoImpl) GetSession() (mongo.Session, error) {
	return r.BaseRepo.GetSession()
}

func (r *PostRepoImpl) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	return r.DeleteOneById(ctx, id)
}

func (r *PostRepoImpl) CreateMany(ctx context.Context, datas []any) (*mongo.InsertManyResult, error) {
	return r.InsertMany(ctx, datas)
}
