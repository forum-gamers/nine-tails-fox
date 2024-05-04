package controllers

import (
	"context"
	"sync"
	"time"

	protobuf "github.com/forum-gamers/nine-tails-fox/generated/like"
	"github.com/forum-gamers/nine-tails-fox/pkg/like"
	"github.com/forum-gamers/nine-tails-fox/pkg/post"
	"github.com/forum-gamers/nine-tails-fox/pkg/preference"
	"github.com/forum-gamers/nine-tails-fox/pkg/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LikeService struct {
	protobuf.UnimplementedLikeServiceServer
	GetUser               func(ctx context.Context) user.User
	LikeRepo              like.LikeRepo
	PostRepo              post.PostRepo
	UserPreferenceRepo    preference.PreferenceRepo
	UserPreferenceService preference.PreferenceService
}

func (s *LikeService) CreateLike(ctx context.Context, in *protobuf.LikeIdPayload) (*protobuf.Like, error) {
	if in.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "postId is required")
	}

	postId, err := primitive.ObjectIDFromHex(in.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid PostId")
	}

	var post post.Post
	if err := s.PostRepo.FindById(ctx, postId, &post); err != nil {
		return nil, err
	}

	var data like.Like
	userId := s.GetUser(ctx).Id
	if err := s.LikeRepo.GetLikesByUserIdAndPostId(ctx, postId, userId, &data); err != nil {
		if e, ok := status.FromError(err); ok && e.Code() != codes.NotFound {
			return nil, err
		}
	}

	if data.Id != primitive.NilObjectID {
		return nil, status.Error(codes.AlreadyExists, "Conflict")
	}

	userPreference, err := s.UserPreferenceRepo.FindByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	session, err := s.PostRepo.GetSession()
	if err != nil {
		return nil, status.Error(codes.Unavailable, "Failed get session")
	}
	defer session.EndSession(ctx)

	dbCtx := mongo.NewSessionContext(ctx, session)
	if err := session.StartTransaction(); err != nil {
		return nil, status.Error(codes.Unavailable, "Failed start DB Operations")
	}

	var wg sync.WaitGroup
	errCh := make(chan error)
	result := like.Like{
		UserId:    userId,
		PostId:    postId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	handlers := []func(){
		func() {
			defer wg.Done()
			id, err := s.LikeRepo.AddLikes(dbCtx, &result)
			if err != nil {
				errCh <- err
				return
			}
			result.Id = id
			errCh <- nil
		},
		func() {
			defer wg.Done()
			errCh <- s.UserPreferenceRepo.UpdateTags(dbCtx, userId, s.UserPreferenceService.CreateUserNewTags(dbCtx, userPreference, post.Tags))
		},
	}

	for _, handler := range handlers {
		wg.Add(1)
		go handler()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			session.AbortTransaction(dbCtx)
			return nil, err
		}
	}

	if err := session.CommitTransaction(dbCtx); err != nil {
		session.AbortTransaction(dbCtx)
		return nil, err
	}

	return &protobuf.Like{
		XId:       result.Id.Hex(),
		UserId:    result.UserId,
		PostId:    result.PostId.Hex(),
		CreatedAt: result.CreatedAt.Local().String(),
		UpdatedAt: result.UpdatedAt.Local().String(),
	}, nil
}

func (s *LikeService) DeleteLike(ctx context.Context, in *protobuf.LikeIdPayload) (*protobuf.Messages, error) {
	if in.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "postId is required")
	}

	postId, err := primitive.ObjectIDFromHex(in.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid ObjectId")
	}

	userId := s.GetUser(ctx).Id
	var data like.Like
	if err := s.LikeRepo.GetLikesByUserIdAndPostId(ctx, postId, userId, &data); err != nil {
		return nil, err
	}

	if err := s.LikeRepo.DeleteLike(ctx, postId, userId); err != nil {
		return nil, err
	}

	return &protobuf.Messages{Message: "success"}, nil
}
