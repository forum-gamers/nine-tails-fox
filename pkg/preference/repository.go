package preference

import (
	"context"
	"time"

	b "github.com/forum-gamers/nine-tails-fox/pkg/base"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewPreferenceRepo() PreferenceRepo {
	return &PreferenceRepoImpl{b.NewBaseRepo(b.GetCollection(b.Preference))}
}

func (r *PreferenceRepoImpl) Create(ctx context.Context, userId string) (UserPreference, error) {
	data := UserPreference{
		UserId:    userId,
		Tags:      []TagPreference{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if result, err := r.BaseRepo.Create(ctx, data); err != nil {
		return data, err
	} else {
		data.Id = result
	}
	return data, nil
}

func (r *PreferenceRepoImpl) FindByUserId(ctx context.Context, userId string) (data UserPreference, err error) {
	err = r.FindOneByQuery(ctx, bson.M{"userId": userId}, &data)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				return r.Create(ctx, userId)
			}
		}
	}
	return
}

func (r *PreferenceRepoImpl) UpdateTags(ctx context.Context, userId string, tags []TagPreference) error {
	_, err := r.UpdateOne(ctx, bson.M{"userId": userId}, bson.M{
		"$set": bson.M{
			"tags": tags,
		},
	})
	return err
}
