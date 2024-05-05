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

type PostResponse struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	UserId       string             `json:"userId" bson:"userId"`
	Text         string             `json:"text" bson:"text"`
	Media        []Media            `json:"media" bson:"media"`
	AllowComment bool               `json:"allowComment"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
	CountLike    int                `json:"countLike" bson:"countLike"`
	CountComment int                `json:"countComment" bson:"countComment"`
	CountShare   int                `json:"countShare" bson:"countShare"`
	IsLiked      bool               `json:"isLiked" bson:"isLiked"`
	IsShared     bool               `json:"isShared" bson:"isShared"`
	Tags         []string           `json:"tags" bson:"tags"`
	Privacy      string             `json:"privacy" bson:"privacy"`
	TotalData    int                `json:"totalData" bson:"totalData"`
}

type TopTags struct {
	Id    string               `json:"_id" bson:"_id"`
	Count int                  `json:"count" bson:"count"`
	Posts []primitive.ObjectID `json:"posts" bson:"posts"`
}
