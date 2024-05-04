package like

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewLikeService(repo LikeRepo) LikeService {
	return &LikeServiceImpl{repo}
}

func (s *LikeServiceImpl) InsertManyAndBindIds(ctx context.Context, likes []Like) error {
	var payload []any

	for _, data := range likes {
		data.Id = primitive.NilObjectID
		payload = append(payload, data)
	}

	ids, err := s.Repo.CreateMany(ctx, payload)
	if err != nil {
		return err
	}

	for i := 0; i < len(ids.InsertedIDs); i++ {
		id := ids.InsertedIDs[i].(primitive.ObjectID)
		likes[i].Id = id
	}
	return nil
}
