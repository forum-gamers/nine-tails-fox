package post

import (
	"context"

	protobuf "github.com/forum-gamers/nine-tails-fox/generated/post"
	"github.com/forum-gamers/nine-tails-fox/pkg/base"
	"github.com/forum-gamers/nine-tails-fox/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepo interface {
	Create(ctx context.Context, data *Post) error
	FindById(ctx context.Context, id primitive.ObjectID, data *Post) error
	GetSession() (mongo.Session, error)
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	CreateMany(ctx context.Context, datas []any) (*mongo.InsertManyResult, error)
	GetPublicContent(ctx context.Context, userId string, query *protobuf.GetPostParams) ([]PostResponse, error)
	GetUserPost(ctx context.Context, userId string, query *protobuf.Pagination) ([]PostResponse, error)
	GetUserPostMedia(ctx context.Context, userId string, query *protobuf.Pagination) ([]PostResponse, error)
	GetTopTags(ctx context.Context, query *protobuf.Pagination) ([]TopTags, error)
}

type PostRepoImpl struct {
	base.BaseRepo
	utils.QueryUtils
}

type PostService interface {
	InsertManyAndBindIds(ctx context.Context, datas []Post) error
	GetPostTags(text string) []string
	CreatePostPayload(userId, text, privacy string, allowComment bool, media []Media, tags []string) Post
}

type PostServiceImpl struct{ Repo PostRepo }
