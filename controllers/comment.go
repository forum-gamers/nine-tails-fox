package controllers

import (
	"context"

	protobuf "github.com/forum-gamers/nine-tails-fox/generated/comment"
	"github.com/forum-gamers/nine-tails-fox/pkg/comment"
	"github.com/forum-gamers/nine-tails-fox/pkg/post"
	"github.com/forum-gamers/nine-tails-fox/pkg/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommentService struct {
	protobuf.UnimplementedCommentServiceServer
	GetUser        func(ctx context.Context) user.User
	PostRepo       post.PostRepo
	CommentRepo    comment.CommentRepo
	CommentService comment.CommentService
}

func (s *CommentService) CreateComment(ctx context.Context, req *protobuf.CommentForm) (*protobuf.Comment, error) {
	switch true {
	case req.Text == "":
		return nil, status.Error(codes.InvalidArgument, "text is required")
	case req.PostId == "":
		return nil, status.Error(codes.InvalidArgument, "postId is required")
	default:
		break
	}

	postId, err := primitive.ObjectIDFromHex(req.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid PostId")
	}

	var postData post.Post
	if err := s.PostRepo.FindById(ctx, postId, &postData); err != nil {
		return nil, err
	}

	commentPayload := s.CommentService.CreatePayload(req.Text, postId, s.GetUser(ctx).Id)
	if err := s.CommentRepo.CreateComment(ctx, &commentPayload); err != nil {
		return nil, err
	}

	return &protobuf.Comment{
		XId:       commentPayload.Id.Hex(),
		Text:      commentPayload.Text,
		UserId:    commentPayload.UserId,
		PostId:    commentPayload.PostId.Hex(),
		CreatedAt: commentPayload.CreatedAt.Local().String(),
		UpdatedAt: commentPayload.UpdatedAt.Local().String(),
	}, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, req *protobuf.CommentIdPayload) (*protobuf.Messages, error) {
	if req.XId == "" {
		return nil, status.Error(codes.InvalidArgument, "_id is required")
	}

	commentId, err := primitive.ObjectIDFromHex(req.XId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid ObjectId")
	}

	var data comment.Comment
	if err := s.CommentRepo.FindById(ctx, commentId, &data); err != nil {
		return nil, err
	}

	if err := s.CommentRepo.DeleteOne(ctx, commentId); err != nil {
		return nil, err
	}

	return &protobuf.Messages{Message: "success"}, nil
}
