package like

import (
	"context"

	"github.com/forum-gamers/nine-tails-fox/pkg/base"
	"github.com/forum-gamers/nine-tails-fox/pkg/post"
	"github.com/forum-gamers/nine-tails-fox/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LikeRepo interface {
	DeletePostLikes(ctx context.Context, postId primitive.ObjectID) error
	GetLikesByUserIdAndPostId(ctx context.Context, postId primitive.ObjectID, userId string, result *Like) error
	AddLikes(ctx context.Context, like *Like) (primitive.ObjectID, error)
	DeleteLike(ctx context.Context, postId primitive.ObjectID, userId string) error
	CreateMany(ctx context.Context, datas []any) (*mongo.InsertManyResult, error)
	GetSession() (mongo.Session, error)
	FindUserLikedPost(ctx context.Context, userId string, in base.Pagination) ([]post.PostResponse, error)
	CountPostLikes(ctx context.Context, ids []primitive.ObjectID) ([]PostLikes, error)
}

type LikeRepoImpl struct {
	base.BaseRepo
	utils.QueryUtils
}

type LikeService interface {
	InsertManyAndBindIds(ctx context.Context, likes []Like) error
}

type LikeServiceImpl struct{ Repo LikeRepo }
