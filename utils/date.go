package utils

import "time"

func ParseDate(dateString string) time.Time {
	t, _ := time.Parse("2006-01-02", dateString)
	return t
}

func FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}
