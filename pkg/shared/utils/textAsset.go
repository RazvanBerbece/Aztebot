package utils

import (
	"fmt"
	"os"
)

func GetTextFromFile(filepath string) string {
	b, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}
