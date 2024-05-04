package comment

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewCommentService(repo CommentRepo) CommentService {
	return &CommentServiceImpl{repo}
}

func (s *CommentServiceImpl) CreatePayload(text string, postId primitive.ObjectID, userId string) Comment {
	return Comment{
		UserId:    userId,
		Text:      text,
		PostId:    postId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Reply:     []ReplyComment{},
	}
}

func (s *CommentServiceImpl) InsertManyAndBindIds(ctx context.Context, datas []Comment) error {
	var payload []any

	for _, data := range datas {
		payload = append(payload, data)
	}

	ids, err := s.Repo.CreateMany(ctx, payload)
	if err != nil {
		return err
	}

	for i := 0; i < len(ids.InsertedIDs); i++ {
		id := ids.InsertedIDs[i].(primitive.ObjectID)
		datas[i].Id = id
	}
	return nil
}
