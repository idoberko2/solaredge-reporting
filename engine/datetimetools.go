package engine

import "time"

func ComputeStartNextMonth(t time.Time) time.Time {
	t = t.AddDate(0, 1, 0)
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func ComputeEndOfMonth(t time.Time) time.Time {
	return ComputeStartNextMonth(t).Add(-24 * time.Hour)
}

func Min(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t2
	}

	return t1
}
