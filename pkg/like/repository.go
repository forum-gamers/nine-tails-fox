package like

import (
	"context"

	b "github.com/forum-gamers/nine-tails-fox/pkg/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewLikeRepo() LikeRepo {
	return &LikeRepoImpl{b.NewBaseRepo(b.GetCollection(b.Like))}
}

func (r *LikeRepoImpl) DeletePostLikes(ctx context.Context, postId primitive.ObjectID) error {
	return r.DeleteManyByQuery(ctx, bson.M{"postId": postId})
}

func (r *LikeRepoImpl) GetLikesByUserIdAndPostId(ctx context.Context, postId primitive.ObjectID, userId string, result *Like) error {
	return r.FindOneByQuery(ctx, bson.M{"userId": userId, "postId": postId}, &result)
}

func (r *LikeRepoImpl) AddLikes(ctx context.Context, like *Like) (primitive.ObjectID, error) {
	return r.Create(ctx, like)
}

func (r *LikeRepoImpl) DeleteLike(ctx context.Context, postId primitive.ObjectID, userId string) error {
	return r.DeleteOneByQuery(ctx, bson.M{"postId": postId, "userId": userId})
}

func (r *LikeRepoImpl) CreateMany(ctx context.Context, datas []any) (*mongo.InsertManyResult, error) {
	return r.InsertMany(ctx, datas)
}

func (r *LikeRepoImpl) GetSession() (mongo.Session, error) {
	return r.BaseRepo.GetSession()
}
