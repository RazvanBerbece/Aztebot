package utils_test

import (
	"testing"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

func TestGetRandomIntegerInRange(t *testing.T) {

	cases := []struct {
		input          []int
		expectedOutput []int
	}{
		{[]int{1, 6}, []int{1, 2, 3, 4, 5, 6}},
		{[]int{0, 5}, []int{0, 1, 2, 3, 4, 5}},
	}

	for _, c := range cases {
		if output := utils.GetRandomIntegerInRange(c.input[0], c.input[1]); utils.IntInSlice(output, c.expectedOutput) {
			t.Errorf("incorrect output for `(%d,%d) -> %d`: expected to be within range", c.input[0], c.input[1], output)
		}
	}

}
