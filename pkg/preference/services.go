package preference

import (
	"context"
	"time"
)

func NewPreferenceService(r PreferenceRepo) PreferenceService {
	return &PreferenceServiceImpl{r}
}

func (s *PreferenceServiceImpl) CreateUserNewTags(ctx context.Context, data UserPreference, newData []string) []TagPreference {
	newTags := data.Tags
	for _, newData := range newData {
		if !data.IsContainsTag(newData) {
			newTags = append(newTags, TagPreference{
				Value:     newData,
				CreatedAt: time.Now(),
			})
		}
	}
	return newTags
}
