package utils_test

import (
	"testing"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

func TestRoleIsStaffRole(t *testing.T) {

	cases := []struct {
		input          int
		expectedOutput bool
	}{
		{1, false},
		{2, false},
		{3, true},
		{4, false},
		{5, true},
		{6, true},
		{7, true},
		{8, true},
		{19, true},
	}

	for _, c := range cases {
		if output := utils.RoleIsStaffRole(c.input); output != c.expectedOutput {
			t.Errorf("incorrect output for `%d`: expected `%t` but got `%t`", c.input, c.expectedOutput, output)
		}
	}

}

func TestGetCircleAndOrderFromRoleId(t *testing.T) {

	cases := []struct {
		input          int
		expectedOutput []int
	}{
		{1, []int{0, -1}},
		{2, []int{0, -1}},
		{3, []int{0, -1}},
		{4, []int{0, -1}},
		{5, []int{0, -1}},
		{6, []int{0, -1}},
		{7, []int{0, -1}},
		{8, []int{0, -1}},
		{9, []int{1, 1}},
		{10, []int{1, 1}},
		{11, []int{1, 1}},
		{12, []int{1, 1}},
		{13, []int{1, 2}},
		{14, []int{1, 2}},
		{15, []int{1, 2}},
		{16, []int{1, 3}},
		{17, []int{1, 3}},
		{18, []int{1, 3}},
		{19, []int{1, 3}},
	}

	for _, c := range cases {

		circle, order := utils.GetCircleAndOrderFromRoleId(c.input)

		if circle != c.expectedOutput[0] {
			t.Errorf("incorrect output for `%d`: expected `%d` but got `%d`", c.input, c.expectedOutput, circle)
		}

		if order != c.expectedOutput[1] {
			t.Errorf("incorrect output for `%d`: expected `%d` but got `%d`", c.input, c.expectedOutput, order)
		}

	}

}
