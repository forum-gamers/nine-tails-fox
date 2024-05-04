package reply

import "github.com/forum-gamers/nine-tails-fox/pkg/comment"

type ReplyService interface {
	CreatePayload(text, userId string) comment.ReplyComment
}

type ReplyServiceImpl struct{ Repo comment.CommentRepo }
