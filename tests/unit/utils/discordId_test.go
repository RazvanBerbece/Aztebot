package utils_test

import (
	"testing"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

func TestGetDiscordIdFromMentionFormat(t *testing.T) {

	cases := []struct {
		input          string
		expectedOutput string
	}{
		{"<@253893620837384192>", "253893620837384192"},
		{"<@2538936208373841920>", "2538936208373841920"},
		{"<@1234>", "1234"},
		{"1234", "1234"},
		{"2538936208373841920", "2538936208373841920"},
	}

	for _, c := range cases {
		if output := utils.GetDiscordIdFromMentionFormat(c.input); c.expectedOutput != output {
			t.Errorf("incorrect actual output %s: expected %s", output, c.expectedOutput)
		}
	}

}
