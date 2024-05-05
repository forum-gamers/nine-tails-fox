package helpers

import "time"

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func StartOfYear(d time.Time) time.Time {
	return time.Date(d.Year(), 1, 1, 0, 0, 0, 0, d.Location())
}

func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year()+1, 1, 1, 0, 0, 0, 0, t.Location()).Add(-time.Second)
}

func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func JakartaDate(t time.Time) time.Time {
	return t.UTC().Truncate(24 * time.Hour).Add(24 * time.Hour).In(time.FixedZone("Asia/Jakarta", 7*60*60))
}

func CalcDifference(now, target time.Time) int {
	return int(
		time.Date(
			now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
			Sub(
				time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, target.Location())).
			Hours() / 24,
	)
}

func EndOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location()).Add(-time.Second)
}

func EndOfDay(t time.Time) time.Time {
	t = t.Add(24 * time.Hour)

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Add(-time.Second)
}
