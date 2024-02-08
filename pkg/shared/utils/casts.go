package utils

import "strconv"

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
