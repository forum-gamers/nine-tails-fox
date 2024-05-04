package post

import (
	"context"

	"github.com/forum-gamers/nine-tails-fox/pkg/base"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepo interface {
	Create(ctx context.Context, data *Post) error
	FindById(ctx context.Context, id primitive.ObjectID, data *Post) error
	GetSession() (mongo.Session, error)
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	CreateMany(ctx context.Context, datas []any) (*mongo.InsertManyResult, error)
}

type PostRepoImpl struct{ base.BaseRepo }

type PostService interface {
	InsertManyAndBindIds(ctx context.Context, datas []Post) error
	GetPostTags(text string) []string
	CreatePostPayload(userId, text, privacy string, allowComment bool, media []Media, tags []string) Post
}

type PostServiceImpl struct{ Repo PostRepo }
