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
