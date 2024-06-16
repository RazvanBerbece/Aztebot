package utils_test

import (
	"testing"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

type TestStatsInput struct {
	NumMessagesSent      int
	NumSlashCommandsUsed int
	NumReactionsReceived int
	TsVc                 int
	TsMusic              int
}

func TestCalculateExperiencePointsFromStats(t *testing.T) {

	cases := []struct {
		input          TestStatsInput
		expectedOutput int
	}{
		{TestStatsInput{
			NumMessagesSent:      0,
			NumSlashCommandsUsed: 0,
			NumReactionsReceived: 0,
			TsVc:                 0,
			TsMusic:              0,
		}, 0},
		{TestStatsInput{
			NumMessagesSent:      1,
			NumSlashCommandsUsed: 0,
			NumReactionsReceived: 0,
			TsVc:                 0,
			TsMusic:              0,
		}, 0},
		{TestStatsInput{
			NumMessagesSent:      153,
			NumSlashCommandsUsed: 10,
			NumReactionsReceived: 7,
			TsVc:                 300,
			TsMusic:              0,
		}, 123},
	}

	for _, c := range cases {
		if output := utils.CalculateExperiencePointsFromStats(c.input.NumMessagesSent, c.input.NumSlashCommandsUsed, c.input.NumReactionsReceived, c.input.TsVc, c.input.TsMusic); output != c.expectedOutput {
			t.Errorf("incorrect output: expected `%d` but got `%d`", c.expectedOutput, output)
		}
	}

}
