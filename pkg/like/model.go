package like

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Like struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId    string             `json:"userId" bson:"userId,omitempty"`
	PostId    primitive.ObjectID `json:"postId" bson:"postId,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type PostLikes struct {
	Id        primitive.ObjectID `json:"postId" bson:"_id"`
	TotalLike int                `json:"totalLike" bson:"totalLike"`
}
