package preference

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPreference struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	UserId    string             `json:"userId" bson:"userId"`
	Tags      []TagPreference    `json:"tags" bson:"tags"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type TagPreference struct {
	Value     string    `json:"value" bson:"value"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}
