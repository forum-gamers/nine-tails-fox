package post

import (
	"context"
	"time"

	protobuf "github.com/forum-gamers/nine-tails-fox/generated/post"
	h "github.com/forum-gamers/nine-tails-fox/helpers"
	b "github.com/forum-gamers/nine-tails-fox/pkg/base"
	"github.com/forum-gamers/nine-tails-fox/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
)

func NewPostRepo(q utils.QueryUtils) PostRepo {
	return &PostRepoImpl{b.NewBaseRepo(b.GetCollection(b.Post)), q}
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

func (r *PostRepoImpl) GetPublicContent(ctx context.Context, userId string, query *protobuf.GetPostParams) ([]PostResponse, error) {
	now := time.Now().UTC()
	orQuery := bson.A{}

	if query.Tags != nil && len(query.Tags) > 0 {
		orQuery = append(orQuery, bson.D{
			{Key: "tags", Value: bson.D{
				{Key: "$in", Value: query.Tags},
			}},
		})
	}

	if query.UserIds != nil && len(query.UserIds) > 0 {
		orQuery = append(orQuery, bson.D{
			{Key: "userId", Value: bson.D{
				{Key: "$in", Value: query.UserIds},
			}},
		})
	}

	orQuery = append(orQuery, bson.D{})
	curr, err := r.BaseRepo.Aggregations(ctx, bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "createdAt",
				Value: bson.D{
					{Key: "$gte", Value: h.StartOfDay(now.AddDate(0, 0, -3))},
				},
			},
			{Key: "privacy", Value: "Public"},
			{Key: "$or", Value: orQuery},
		}}},
		bson.D{
			{Key: "$facet",
				Value: bson.D{
					{Key: "total",
						Value: bson.A{
							bson.D{{Key: "$count", Value: "total"}},
						},
					},
					{Key: "datas",
						Value: bson.A{
							r.NewSkip(int((query.Page - 1) * query.Limit)),
							r.NewLimit(int(query.Limit)),
							r.NewLookup("comment", "_id", "postId", "comment"),
							r.NewLookup("like", "_id", "postId", "like"),
							r.NewLookup("share", "_id", "postId", "share"),
							bson.D{
								{Key: "$addFields",
									Value: bson.D{
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
				},
			},
		},
		r.NewRawUnwind("$datas"),
		r.NewRawUnwind("$total"),
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: "$datas._id"},
					{Key: "userId", Value: "$datas.userId"},
					{Key: "text", Value: "$datas.text"},
					{Key: "media", Value: "$datas.media"},
					{Key: "allowComment", Value: "$datas.allowComment"},
					{Key: "createdAt", Value: "$datas.createdAt"},
					{Key: "updatedAt", Value: "$datas.updatedAt"},
					{Key: "countLike", Value: "$datas.countLike"},
					{Key: "countComment", Value: "$datas.countComment"},
					{Key: "countShare", Value: "$datas.countShare"},
					{Key: "isLiked", Value: "$datas.isLiked"},
					{Key: "isShared", Value: "$datas.isShared"},
					{Key: "tags", Value: "$datas.tags"},
					{Key: "privacy", Value: "$datas.privacy"},
					{Key: "totalData", Value: "$total.total"},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	defer curr.Close(ctx)

	var datas []PostResponse
	for curr.Next(ctx) {
		var data PostResponse
		if err := curr.Decode(&data); err != nil {
			return nil, err
		}
		datas = append(datas, data)
	}

	if len(datas) < 1 {
		return datas, h.NewAppError(codes.NotFound, "data not found")
	}

	return datas, nil
}

func (r *PostRepoImpl) GetUserPost(ctx context.Context, userId string, query *protobuf.Pagination) ([]PostResponse, error) {
	curr, err := r.Aggregations(ctx, bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "userId", Value: userId}}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: -1}}}},
		bson.D{
			{Key: "$facet",
				Value: bson.D{
					{Key: "total",
						Value: bson.A{
							bson.D{{Key: "$count", Value: "total"}},
						},
					},
					{Key: "datas",
						Value: bson.A{
							r.NewSkip(int((query.Page - 1) * query.Limit)),
							r.NewLimit(int(query.Limit)),
							r.NewLookup("comment", "_id", "postId", "comment"),
							r.NewLookup("like", "_id", "postId", "like"),
							r.NewLookup("share", "_id", "posId", "share"),
							bson.D{
								{Key: "$addFields",
									Value: bson.D{
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
				},
			},
		},
		bson.D{{Key: "$unwind", Value: "$datas"}},
		bson.D{{Key: "$unwind", Value: "$total"}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: "$datas._id"},
					{Key: "userId", Value: "$datas.userId"},
					{Key: "text", Value: "$datas.text"},
					{Key: "media", Value: "$datas.media"},
					{Key: "allowComment", Value: "$datas.allowComment"},
					{Key: "createdAt", Value: "$datas.createdAt"},
					{Key: "updatedAt", Value: "$datas.updatedAt"},
					{Key: "countLike", Value: "$datas.countLike"},
					{Key: "countComment", Value: "$datas.countComment"},
					{Key: "countShare", Value: "$datas.countShare"},
					{Key: "isLiked", Value: "$datas.isLiked"},
					{Key: "isShared", Value: "$datas.isShared"},
					{Key: "tags", Value: "$datas.tags"},
					{Key: "privacy", Value: "$datas.privacy"},
					{Key: "totalData", Value: "$total.total"},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	defer curr.Close(ctx)

	var datas []PostResponse
	for curr.Next(ctx) {
		var data PostResponse
		if err := curr.Decode(&data); err != nil {
			return nil, err
		}
		datas = append(datas, data)
	}

	if len(datas) < 1 {
		return datas, h.NewAppError(codes.NotFound, "data not found")
	}

	return datas, nil
}

func (r *PostRepoImpl) GetUserPostMedia(ctx context.Context, userId string, query *protobuf.Pagination) ([]PostResponse, error) {
	curr, err := r.Aggregations(ctx, bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "userId", Value: userId},
					{Key: "media", Value: bson.D{{Key: "$exists", Value: true}}},
				},
			},
		},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: -1}}}},
		bson.D{
			{Key: "$facet",
				Value: bson.D{
					{Key: "total",
						Value: bson.A{
							bson.D{{Key: "$count", Value: "total"}},
						},
					},
					{Key: "datas",
						Value: bson.A{
							r.NewSkip(int((query.Page - 1) * query.Limit)),
							r.NewLimit(int(query.Limit)),
							r.NewLookup("comment", "_id", "postId", "comment"),
							r.NewLookup("like", "_id", "postId", "like"),
							r.NewLookup("share", "_id", "posId", "share"),
							bson.D{
								{Key: "$addFields",
									Value: bson.D{
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
				},
			},
		},
		bson.D{{Key: "$unwind", Value: "$datas"}},
		bson.D{{Key: "$unwind", Value: "$total"}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: "$datas._id"},
					{Key: "userId", Value: "$datas.userId"},
					{Key: "text", Value: "$datas.text"},
					{Key: "media", Value: "$datas.media"},
					{Key: "allowComment", Value: "$datas.allowComment"},
					{Key: "createdAt", Value: "$datas.createdAt"},
					{Key: "updatedAt", Value: "$datas.updatedAt"},
					{Key: "countLike", Value: "$datas.countLike"},
					{Key: "countComment", Value: "$datas.countComment"},
					{Key: "countShare", Value: "$datas.countShare"},
					{Key: "isLiked", Value: "$datas.isLiked"},
					{Key: "isShared", Value: "$datas.isShared"},
					{Key: "tags", Value: "$datas.tags"},
					{Key: "privacy", Value: "$datas.privacy"},
					{Key: "totalData", Value: "$total.total"},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	defer curr.Close(ctx)

	var datas []PostResponse
	for curr.Next(ctx) {
		var data PostResponse
		if err := curr.Decode(&data); err != nil {
			return nil, err
		}

		datas = append(datas, data)
	}

	if len(datas) < 1 {
		return datas, h.NewAppError(codes.NotFound, "data not found")
	}

	return datas, nil
}

func (r *PostRepoImpl) GetTopTags(ctx context.Context, query *protobuf.Pagination) ([]TopTags, error) {
	cursor, err := r.Aggregations(ctx, bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "createdAt", Value: bson.D{{Key: "$gte", Value: h.StartOfDay(time.Now())}}}}}},
		r.NewRawUnwind("$tags"),
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$tags"},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "posts", Value: bson.D{{Key: "$addToSet", Value: "$_id"}}},
			},
			},
		},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
		bson.D{{Key: "$skip", Value: (query.Page - 1) * query.Limit}},
		bson.D{{Key: "$limit", Value: query.Limit}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var datas []TopTags
	for cursor.Next(ctx) {
		var data TopTags
		if err := cursor.Decode(&data); err != nil {
			return datas, err
		}
		datas = append(datas, data)
	}

	if len(datas) < 1 {
		return datas, h.NewAppError(codes.NotFound, "data not found")
	}

	return datas, nil
}
