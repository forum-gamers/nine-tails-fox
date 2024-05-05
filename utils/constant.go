package utils

import "go.mongodb.org/mongo-driver/bson"

type QueryUtils interface {
	NewLookup(from, localField, foreignField, as string) bson.D
	NewSkip(val int) bson.D
	NewRawUnwind(val string) bson.D
	NewLimit(val int) bson.D
	IsDo(key, input, userId string) bson.E
	NewCountComment(key, field string) bson.E
}

type QueryUtilsImpl struct{}
