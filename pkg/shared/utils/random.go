package utils

import (
	"math/rand"
)

func GetRandomIntegerInRange(min int, max int, seed int64) int {
	r := rand.New(rand.NewSource(seed))
	return r.Intn(max-min) + min
}
