package utils_test

import (
	"testing"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

func TestValidateDiscordUserId(t *testing.T) {

	cases := []struct {
		input          string
		expectedOutput bool
	}{
		{"", false},
		{"123", false},
		{"abc", false},
		{"-1", false},
		{"5736595333610209410", true},
		{"573659533361020941011", false},
		{"573659533361020941", true},
		{"57365953336102094", true},
		{"5736595333610209", false},
	}

	for _, c := range cases {
		if output := utils.IsValidDiscordUserId(c.input); output != c.expectedOutput {
			t.Errorf("incorrect output for `%s`: expected `%t` but got `%t`", c.input, c.expectedOutput, output)
		}
	}

}

func TestIsValidReasonMessage(t *testing.T) {

	cases := []struct {
		input          string
		expectedOutput bool
	}{
		{"", true},
		{"123", true},
		{"abc", true},
		{"-1", true},
		{"This is a real, valid reason", true},
		{"573659533361020941057365953336102094105736595333610209410573659533361020941057365953336102094105736595333610209410573659533361020941057365953336102094105736595333610209410573659533361020941057365953336102094105736595333610209410573659533361020941057365953336102094105736595333610209410573659533361020941057365953336102094105736595333610209410573659533361020941057365953336102094105736595333610209410573659533361020941057365953336102094105736595333610209410573659533361020941057365953336102094105736595333610209410573659533361020941057365953336102094105736595333610209410573659533361020941057365953336102094105736595333610209410573659533361020941057365953336102094105736595333610209410573659533361020941057365953336102094105736595333610209410", false},
	}

	for _, c := range cases {
		if output := utils.IsValidReasonMessage(c.input); output != c.expectedOutput {
			t.Errorf("incorrect output for `%s`: expected `%t` but got `%t`", c.input, c.expectedOutput, output)
		}
	}

}
