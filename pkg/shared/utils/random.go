package utils

import (
	"math/rand"
	"time"
)

func GetRandomIntegerInRange(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}

func GetRandomFromArray(array []string) string {
	return array[rand.Intn(len(array))]
}
