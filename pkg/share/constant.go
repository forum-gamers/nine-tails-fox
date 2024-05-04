package share

import (
	"context"

	"github.com/forum-gamers/nine-tails-fox/pkg/base"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShareRepo interface {
	DeleteMany(ctx context.Context, postId primitive.ObjectID) error
}

type ShareRepoImpl struct{ base.BaseRepo }
