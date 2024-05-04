package preference

import (
	"context"

	"github.com/forum-gamers/nine-tails-fox/pkg/base"
)

type PreferenceRepo interface {
	Create(ctx context.Context, userId string) (UserPreference, error)
	FindByUserId(ctx context.Context, userId string) (UserPreference, error)
	UpdateTags(ctx context.Context, userId string, tags []TagPreference) error
}

type PreferenceRepoImpl struct{ base.BaseRepo }

type PreferenceService interface {
	CreateUserNewTags(ctx context.Context, data UserPreference, newData []string) []TagPreference
}

type PreferenceServiceImpl struct{ Repo PreferenceRepo }
