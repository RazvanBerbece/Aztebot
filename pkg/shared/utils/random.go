package utils

import (
	"math/rand"
)

func GetRandomIntegerInRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func GetRandomFromArray(array []string) string {
	return array[rand.Intn(len(array))]
}

// len(weights) = len(array); each i_th element is associated
func GetRandomFromArrayWithWeights(array []string, weights []float64) string {
	return array[rand.Intn(len(array))]
}
