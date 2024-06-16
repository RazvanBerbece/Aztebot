package utils

import (
	"math"
	"time"
)

func HumanReadableDuration(totalSeconds float64) (int, int, int, int) {

	const day = 86400
	const hour = 3600
	const minute = 60

	days := math.Floor(totalSeconds / day)
	hours := math.Floor((totalSeconds - days*day) / hour)
	minutes := math.Floor((totalSeconds - days*day - hours*hour) / minute)
	seconds := totalSeconds - days*day - hours*hour - minutes*minute

	return int(days), int(hours), int(minutes), int(seconds)

}

func FormatUnixAsString(timestamp int64, format string) string {

	var ts time.Time
	var timeString string

	ts = time.Unix(timestamp, 0).UTC()
	timeString = ts.Format(format) // e.g -> "Mon, 02 Jan 2006 15:04:05 MST"

	return timeString

}
