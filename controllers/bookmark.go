package controllers

import (
	"context"

	"github.com/forum-gamers/nine-tails-fox/generated"
	protobuf "github.com/forum-gamers/nine-tails-fox/generated/bookmark"
	"github.com/forum-gamers/nine-tails-fox/pkg/base"
	"github.com/forum-gamers/nine-tails-fox/pkg/bookmark"
	"github.com/forum-gamers/nine-tails-fox/pkg/post"
	"github.com/forum-gamers/nine-tails-fox/pkg/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookmarkService struct {
	protobuf.UnimplementedBookmarkServiceServer
	GetUser         func(ctx context.Context) user.User
	PostRepo        post.PostRepo
	BookmarkRepo    bookmark.BookmarkRepo
	BookmarkService bookmark.BookmarkService
}

func (s *BookmarkService) CreateBookmark(ctx context.Context, req *protobuf.PostIdPayload) (*protobuf.Bookmark, error) {
	if req.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "postId is required")
	}

	postId, err := primitive.ObjectIDFromHex(req.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid PostId")
	}

	var postData post.Post
	if err := s.PostRepo.FindById(ctx, postId, &postData); err != nil {
		return nil, err
	}

	userId := s.GetUser(ctx).Id
	var bookmarkData bookmark.Bookmark
	if err := s.BookmarkRepo.FindByPostIdAndUserId(ctx, postId, userId, &bookmarkData); err != nil {
		if e, ok := status.FromError(err); ok && e.Code() != codes.NotFound {
			return nil, err
		}
	}

	if bookmarkData.Id != primitive.NilObjectID {
		return nil, status.Error(codes.AlreadyExists, "Conflict")
	}

	data := s.BookmarkService.CreatePayload(postId, userId)
	if err := s.BookmarkRepo.CreateOne(ctx, &data); err != nil {
		return nil, err
	}

	return &protobuf.Bookmark{
		XId:       data.Id.Hex(),
		PostId:    data.PostId.Hex(),
		UserId:    data.UserId,
		CreatedAt: data.CreatedAt.Local().String(),
		UpdatedAt: data.UpdatedAt.Local().String(),
	}, nil
}

func (s *BookmarkService) DeleteBookmark(ctx context.Context, req *protobuf.IdPayload) (*protobuf.Messages, error) {
	if req.XId == "" {
		return nil, status.Error(codes.InvalidArgument, "_id is required")
	}

	bookmarkId, err := primitive.ObjectIDFromHex(req.XId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid ObjectID")
	}

	var data bookmark.Bookmark
	if err := s.BookmarkRepo.FindById(ctx, bookmarkId, &data); err != nil {
		return nil, err
	}

	if err := s.BookmarkRepo.DeleteOneById(ctx, bookmarkId); err != nil {
		return nil, err
	}
	return &protobuf.Messages{Message: "success"}, nil
}

func (s *BookmarkService) GetMyBookmarks(ctx context.Context, req *protobuf.PaginationWithPostId) (*protobuf.RespWithMetadata, error) {
	if req.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "postId is required")
	}

	postId, err := primitive.ObjectIDFromHex(req.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid objectId")
	}

	id := s.GetUser(ctx).Id
	data, err := s.BookmarkRepo.FindMyBookmarks(ctx, postId, id, base.Pagination{
		Page:  uint32(req.Page),
		Limit: uint32(req.Limit),
	})
	if err != nil {
		return nil, err
	}

	return &protobuf.RespWithMetadata{
		TotalData: int64(data[0].TotalData),
		Page:      req.Page,
		Limit:     req.Limit,
		Data:      generated.ParseBookmarkPostRespToProto(data),
	}, nil
}
