package post

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewPostService(repo PostRepo) PostService {
	return &PostServiceImpl{repo}
}

func (s *PostServiceImpl) InsertManyAndBindIds(ctx context.Context, datas []Post) error {
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

func (s *PostServiceImpl) GetPostTags(text string) []string {
	modified := text
	for _, p := range "!@#$%^&*)(_=+?.,;:'" {
		modified = strings.ReplaceAll(modified, string(p), " ")
	}
	return strings.Split(modified, " ")
}

func (s *PostServiceImpl) CreatePostPayload(userId, text, privacy string, allowComment bool, media []Media, tags []string) Post {
	return Post{
		UserId:       userId,
		Text:         text,
		Media:        media,
		AllowComment: allowComment,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Tags:         tags,
		Privacy:      privacy,
	}
}
