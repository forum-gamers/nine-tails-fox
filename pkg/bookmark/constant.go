package bookmark

import (
	"context"

	"github.com/forum-gamers/nine-tails-fox/pkg/base"
	"github.com/forum-gamers/nine-tails-fox/pkg/post"
	"github.com/forum-gamers/nine-tails-fox/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookmarkRepo interface {
	CreateOne(ctx context.Context, data *Bookmark) error
	FindOne(ctx context.Context, query any, result *Bookmark) error
	FindById(ctx context.Context, id primitive.ObjectID, result *Bookmark) error
	DeleteOneById(ctx context.Context, id primitive.ObjectID) error
	FindByPostIdAndUserId(ctx context.Context, postId primitive.ObjectID, userId string, result *Bookmark) error
	FindMyBookmarks(ctx context.Context, postId primitive.ObjectID, userId string, query base.Pagination) (result []post.PostResponse, err error)
}

type BookmarkRepoImpl struct {
	base.BaseRepo
	utils.QueryUtils
}

type BookmarkService interface {
	CreatePayload(postId primitive.ObjectID, userId string) Bookmark
}

type BookmarkServiceImpl struct{ Repo BookmarkRepo }
