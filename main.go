package main

import (
	"log"
	"net"
	"os"

	cc "github.com/forum-gamers/nine-tails-fox/controllers"
	"github.com/forum-gamers/nine-tails-fox/database"
	bookmarkProto "github.com/forum-gamers/nine-tails-fox/generated/bookmark"
	commentProto "github.com/forum-gamers/nine-tails-fox/generated/comment"
	likeProto "github.com/forum-gamers/nine-tails-fox/generated/like"
	postProto "github.com/forum-gamers/nine-tails-fox/generated/post"
	replyProto "github.com/forum-gamers/nine-tails-fox/generated/reply"
	h "github.com/forum-gamers/nine-tails-fox/helpers"
	"github.com/forum-gamers/nine-tails-fox/interceptors"
	"github.com/forum-gamers/nine-tails-fox/pkg/bookmark"
	"github.com/forum-gamers/nine-tails-fox/pkg/comment"
	"github.com/forum-gamers/nine-tails-fox/pkg/like"
	"github.com/forum-gamers/nine-tails-fox/pkg/post"
	"github.com/forum-gamers/nine-tails-fox/pkg/preference"
	"github.com/forum-gamers/nine-tails-fox/pkg/reply"
	"github.com/forum-gamers/nine-tails-fox/pkg/share"
	"github.com/forum-gamers/nine-tails-fox/utils"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	h.PanicIfError(godotenv.Load())
	database.Connection()

	address := os.Getenv("PORT")
	if address == "" {
		address = "50052"
	}

	lis, err := net.Listen("tcp", ":"+address)
	if err != nil {
		log.Fatalf("Failed to listen : %s", err.Error())
	}

	query := utils.NewQueryUtils()

	//repository
	postRepo := post.NewPostRepo(query)
	likeRepo := like.NewLikeRepo(query)
	commentRepo := comment.NewCommentRepo(query)
	shareRepo := share.NewShareRepo()
	userPreferenceRepo := preference.NewPreferenceRepo()
	bookmarkRepo := bookmark.NewBookMarkRepo(query)

	//services
	postService := post.NewPostService(postRepo)
	userPreferenceService := preference.NewPreferenceService(userPreferenceRepo)
	commentService := comment.NewCommentService(commentRepo)
	bookmarkService := bookmark.NewBookMarkService(bookmarkRepo)
	replyService := reply.NewReplyService(commentRepo)

	interceptor := interceptors.NewInterCeptor()
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor.Logging, interceptor.UnaryAuthentication),
	)

	postProto.RegisterPostServiceServer(grpcServer, &cc.PostService{
		GetUser:     interceptor.GetUserFromCtx,
		PostRepo:    postRepo,
		PostService: postService,
		LikeRepo:    likeRepo,
		CommentRepo: commentRepo,
		ShareRepo:   shareRepo,
	})
	likeProto.RegisterLikeServiceServer(grpcServer, &cc.LikeService{
		GetUser:               interceptor.GetUserFromCtx,
		LikeRepo:              likeRepo,
		PostRepo:              postRepo,
		UserPreferenceRepo:    userPreferenceRepo,
		UserPreferenceService: userPreferenceService,
	})
	commentProto.RegisterCommentServiceServer(grpcServer, &cc.CommentService{
		GetUser:        interceptor.GetUserFromCtx,
		PostRepo:       postRepo,
		CommentRepo:    commentRepo,
		CommentService: commentService,
	})
	bookmarkProto.RegisterBookmarkServiceServer(grpcServer, &cc.BookmarkService{
		GetUser:         interceptor.GetUserFromCtx,
		PostRepo:        postRepo,
		BookmarkRepo:    bookmarkRepo,
		BookmarkService: bookmarkService,
	})
	replyProto.RegisterReplyServiceServer(grpcServer, &cc.ReplyService{
		GetUser:        interceptor.GetUserFromCtx,
		CommentRepo:    commentRepo,
		CommentService: commentService,
		ReplyService:   replyService,
	})

	log.Printf("Starting to serve in port : %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %s", err.Error())
	}
}
