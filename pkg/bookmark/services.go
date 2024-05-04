package bookmark

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewBookMarkService(r BookmarkRepo) BookmarkService {
	return &BookmarkServiceImpl{r}
}

func (s *BookmarkServiceImpl) CreatePayload(postId primitive.ObjectID, userId string) Bookmark {
	return Bookmark{
		PostId:    postId,
		UserId:    userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
