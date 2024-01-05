package utils

import "math"

func HumanReadableTimeLength(totalSeconds float64) (int, int, int, int) {

	const day = 86400
	const hour = 3600
	const minute = 60

	days := math.Floor(totalSeconds / day)
	hours := math.Floor((totalSeconds - days*day) / hour)
	minutes := math.Floor((totalSeconds - days*day - hours*hour) / minute)
	seconds := totalSeconds - days*day - hours*hour - minutes*minute

	return int(days), int(hours), int(minutes), int(seconds)

}
