package preference

import "strings"

func (t *UserPreference) IsContainsTag(tag string) bool {
	for _, data := range t.Tags {
		if strings.EqualFold(strings.ToLower(data.Value), strings.ToLower(tag)) {
			return true
		}
	}
	return false
}
