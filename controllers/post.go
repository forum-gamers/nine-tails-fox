package controllers

import (
	"context"
	"sync"

	protobuf "github.com/forum-gamers/nine-tails-fox/generated/post"
	"github.com/forum-gamers/nine-tails-fox/pkg/comment"
	"github.com/forum-gamers/nine-tails-fox/pkg/like"
	"github.com/forum-gamers/nine-tails-fox/pkg/post"
	"github.com/forum-gamers/nine-tails-fox/pkg/share"
	"github.com/forum-gamers/nine-tails-fox/pkg/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostService struct {
	protobuf.UnimplementedPostServiceServer
	GetUser     func(ctx context.Context) user.User
	PostRepo    post.PostRepo
	PostService post.PostService
	LikeRepo    like.LikeRepo
	CommentRepo comment.CommentRepo
	ShareRepo   share.ShareRepo
}

func (s *PostService) CreatePost(ctx context.Context, req *protobuf.PostForm) (*protobuf.Post, error) {
	tags := []string{}
	if len(req.Text) > 0 {
		tags = s.PostService.GetPostTags(req.Text)
	}

	postMedias := make([]post.Media, 0)
	if len(req.Files) > 0 {
		for _, file := range req.Files {
			postMedias = append(postMedias, post.Media{
				Url:  file.Url,
				Type: file.ContentType,
				Id:   file.FileId,
			})
		}
	}

	userId := s.GetUser(ctx).Id
	post := s.PostService.CreatePostPayload(userId, req.Text, req.Privacy, req.AllowComment, postMedias, tags)

	s.PostRepo.Create(context.Background(), &post)
	resultMedia := make([]*protobuf.Media, 0)
	if len(post.Media) > 0 {
		for _, media := range post.Media {
			resultMedia = append(resultMedia, &protobuf.Media{
				Id:   media.Id,
				Url:  media.Url,
				Type: media.Type,
			})
		}
	}

	return &protobuf.Post{
		XId:          post.Id.Hex(),
		UserId:       post.UserId,
		Text:         post.Text,
		Media:        resultMedia,
		AllowComment: post.AllowComment,
		CreatedAt:    post.CreatedAt.String(),
		UpdatedAt:    post.UpdatedAt.String(),
		Tags:         post.Tags,
		Privacy:      post.Privacy,
	}, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *protobuf.PostIdPayload) (*protobuf.Messages, error) {
	if req.XId == "" {
		return nil, status.Error(codes.InvalidArgument, "_id is required")
	}

	postId, err := primitive.ObjectIDFromHex(req.XId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid ObjectId")
	}

	var data post.Post
	if err := s.PostRepo.FindById(ctx, postId, &data); err != nil {
		return nil, err
	}

	user := s.GetUser(ctx)
	if user.Id != data.UserId && user.AccountType != "Admin" {
		return nil, status.Error(codes.Unauthenticated, "Forbidden")
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
	handlers := []func(){
		func() {
			defer wg.Done()
			errCh <- s.LikeRepo.DeletePostLikes(dbCtx, data.Id)
		},
		func() {
			defer wg.Done()
			errCh <- s.PostRepo.DeleteOne(dbCtx, data.Id)
		},
		func() {
			defer wg.Done()
			errCh <- s.ShareRepo.DeleteMany(ctx, data.Id)
		},
		func() {
			defer wg.Done()
			errCh <- s.CommentRepo.DeleteMany(ctx, data.Id)
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

	return &protobuf.Messages{Message: "success"}, nil
}