package utils

import (
	"math/rand"
)

func GetRandomIntegerInRange(min int, max int) int {
	return rand.Intn(max-min) + min
}
