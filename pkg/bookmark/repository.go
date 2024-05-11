package bookmark

import (
	"context"

	h "github.com/forum-gamers/nine-tails-fox/helpers"
	b "github.com/forum-gamers/nine-tails-fox/pkg/base"
	"github.com/forum-gamers/nine-tails-fox/pkg/post"
	"github.com/forum-gamers/nine-tails-fox/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
)

func NewBookMarkRepo(q utils.QueryUtils) BookmarkRepo {
	return &BookmarkRepoImpl{b.NewBaseRepo(b.GetCollection(b.Bookmark)), q}
}

func (r *BookmarkRepoImpl) CreateOne(ctx context.Context, data *Bookmark) error {
	if result, err := r.Create(ctx, data); err != nil {
		return err
	} else {
		data.Id = result
	}
	return nil
}

func (r *BookmarkRepoImpl) FindOne(ctx context.Context, query any, result *Bookmark) error {
	return r.FindOneByQuery(ctx, query, result)
}

func (r *BookmarkRepoImpl) FindById(ctx context.Context, id primitive.ObjectID, result *Bookmark) error {
	return r.FindOneById(ctx, id, result)
}

func (r *BookmarkRepoImpl) DeleteOneById(ctx context.Context, id primitive.ObjectID) error {
	return r.BaseRepo.DeleteOneById(ctx, id)
}

func (r *BookmarkRepoImpl) FindByPostIdAndUserId(ctx context.Context, postId primitive.ObjectID, userId string, result *Bookmark) error {
	return r.FindOneByQuery(ctx, bson.M{"postId": postId, "userId": userId}, result)
}

func (r *BookmarkRepoImpl) FindMyBookmarks(ctx context.Context, postId primitive.ObjectID, userId string, query b.Pagination) (result []post.PostResponse, err error) {
	cursor, err := r.Aggregations(ctx, bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "userId", Value: userId}}}},
		bson.D{
			{Key: "$facet", Value: bson.D{
				{Key: "data", Value: bson.A{
					bson.D{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: 1}}}},
					r.NewSkip(int(query.Page-1) * int(query.Limit)),
					r.NewLimit(int(query.Limit)),
					r.NewLookup("post", "postId", "_id", "post"),
					r.NewRawUnwind("$post"),
					r.NewLookup("like", "post._id", "postId", "like"),
					r.NewLookup("share", "post._id", "postId", "share"),
					r.NewLookup("comment", "post._id", "postId", "comment"),
					bson.D{
						{Key: "$addFields", Value: bson.D{
							{Key: "countLike", Value: bson.D{{Key: "$size", Value: "$like"}}},
							{Key: "countShare", Value: bson.D{{Key: "$size", Value: "$share"}}},
							r.NewCountComment("countComment", "$comment"),
							r.IsDo("isLiked", "$like", userId),
							r.IsDo("isShared", "$share", userId),
						},
						},
					},
				},
				},
				{Key: "total", Value: bson.A{
					bson.D{{Key: "$count", Value: "total"}},
				},
				},
			},
			},
		},
		r.NewRawUnwind("$total"),
		r.NewRawUnwind("$data"),
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: "$data.post._id"},
				{Key: "userId", Value: "$data.post.userId"},
				{Key: "text", Value: "$data.post.text"},
				{Key: "media", Value: "$data.post.media"},
				{Key: "allowComment", Value: "$data.post.allowComment"},
				{Key: "isLiked", Value: "$data.isLiked"},
				{Key: "isShared", Value: "$data.isShared"},
				{Key: "countLike", Value: "$data.countLike"},
				{Key: "countShare", Value: "$data.countShare"},
				{Key: "countComment", Value: "$data.countComment"},
				{Key: "tags", Value: "$data.post.tags"},
				{Key: "privacy", Value: "$data.post.privacy"},
				{Key: "totalData", Value: "$total.total"},
			},
			},
		},
	},
	)

	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var data post.PostResponse
		if err = cursor.Decode(&data); err != nil {
			return
		}
		result = append(result, data)
	}

	if result == nil || len(result) < 1 {
		err = h.NewAppError(codes.NotFound, "data not found")
		return
	}
	return
}
