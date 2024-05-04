package post

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	Url  string `json:"url" bson:"url,omitempty"`
	Type string `json:"type" bson:"type,omitempty"`
	Id   string `json:"id" bson:"id,omitempty"`
}

type Post struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId       string             `json:"userId" bson:"userId,omitempty"`
	Text         string             `json:"text" bson:"text"`
	Media        []Media            `json:"media" bson:"media"`
	AllowComment bool               `json:"allowComment" bson:"allowComment" default:"true"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
	Tags         []string           `json:"tags" bson:"tags,omitempty"`
	Privacy      string             `json:"privacy" bson:"privacy" default:"Public"`
}
