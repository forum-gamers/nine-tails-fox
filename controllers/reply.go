package controllers

import (
	"context"

	protobuf "github.com/forum-gamers/nine-tails-fox/generated/reply"
	"github.com/forum-gamers/nine-tails-fox/pkg/comment"
	"github.com/forum-gamers/nine-tails-fox/pkg/reply"
	"github.com/forum-gamers/nine-tails-fox/pkg/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReplyService struct {
	protobuf.UnimplementedReplyServiceServer
	GetUser        func(ctx context.Context) user.User
	CommentRepo    comment.CommentRepo
	CommentService comment.CommentService
	ReplyService   reply.ReplyService
}

func (s *ReplyService) CreateReply(ctx context.Context, req *protobuf.CommentForm) (*protobuf.Reply, error) {
	if req.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "text is required")
	}

	commentId, err := primitive.ObjectIDFromHex(req.CommentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid CommentId")
	}

	var commentData comment.Comment
	if err := s.CommentRepo.FindById(ctx, commentId, &commentData); err != nil {
		return nil, err
	}

	replyPayload := s.ReplyService.CreatePayload(req.Text, s.GetUser(ctx).Id)
	if err := s.CommentRepo.CreateReply(ctx, commentId, &replyPayload); err != nil {
		return nil, err
	}

	return &protobuf.Reply{
		XId:       replyPayload.Id.Hex(),
		Text:      replyPayload.Text,
		UserId:    replyPayload.UserId,
		CreatedAt: replyPayload.CreatedAt.Local().String(),
		UpdatedAt: replyPayload.UpdatedAt.Local().String(),
	}, nil
}

func (s *ReplyService) DeleteReply(ctx context.Context, req *protobuf.DeleteReplyPayload) (*protobuf.Messages, error) {
	if req.ReplyId == "" {
		return nil, status.Error(codes.InvalidArgument, "replyId is required")
	}

	if req.CommentId == "" {
		return nil, status.Error(codes.InvalidArgument, "commentId is required")
	}

	replyId, err := primitive.ObjectIDFromHex(req.ReplyId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid replyId")
	}

	commentId, err := primitive.ObjectIDFromHex(req.CommentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid commentId")
	}

	var data comment.ReplyComment
	if err := s.CommentRepo.FindReplyById(ctx, commentId, replyId, &data); err != nil {
		return nil, err
	}

	user := s.GetUser(ctx)
	if data.UserId != user.Id || user.AccountType != "Admin" {
		return nil, status.Error(codes.PermissionDenied, "Forbidden")
	}

	if err := s.CommentRepo.DeleteOneReply(ctx, commentId, replyId); err != nil {
		return nil, err
	}

	return &protobuf.Messages{Message: "success"}, nil
}
