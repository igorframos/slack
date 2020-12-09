package utils

import "time"

func NextBusinessDay(hour, min, sec int) time.Time {
	current := time.Now().AddDate(0, 0, 3)
	year, month, day := current.Date()

	today := time.Date(year, month, day, hour, min, sec, 0, current.Location())
	for today.Before(current) || !isBusinessDay(today) {
		today = today.AddDate(0, 0, 1)
	}

	return today
}

// TODO: Account for holidays.
func isBusinessDay(d time.Time) bool {
	return d.Weekday() != time.Sunday && d.Weekday() != time.Saturday
}
