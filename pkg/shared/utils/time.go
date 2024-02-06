package utils

import (
	"math"
	"strconv"
)

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

func StringToInt64(timeString string) (*int64, error) {
	i, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func StringToInt(timeString string) (*int, error) {
	i, err := strconv.Atoi(timeString)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func StringToFloat64(timeString string) (*float64, error) {
	i, err := strconv.ParseFloat(timeString, 64)
	if err != nil {
		return nil, err
	}
	return &i, nil
}
