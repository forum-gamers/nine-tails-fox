package share

import (
	"context"

	b "github.com/forum-gamers/nine-tails-fox/pkg/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewShareRepo() ShareRepo {
	return &ShareRepoImpl{b.NewBaseRepo(b.GetCollection(b.Share))}
}

func (r *ShareRepoImpl) DeleteMany(ctx context.Context, postId primitive.ObjectID) error {
	return r.DeleteManyByQuery(ctx, bson.M{"postId": postId})
}
