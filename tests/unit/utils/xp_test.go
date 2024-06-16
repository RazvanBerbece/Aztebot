package utils_test

import (
	"math"
	"testing"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

type TestStatsInput struct {
	NumMessagesSent      int
	NumSlashCommandsUsed int
	NumReactionsReceived int
	TsVc                 int
	TsMusic              int
	MessageWeight        float64
	SlashCommandWeight   float64
	ReactionsWeight      float64
	TsVcWeight           float64
	TsMusicWeight        float64
}

func TestCalculateExperiencePointsFromStats(t *testing.T) {

	cases := []struct {
		input          TestStatsInput
		expectedOutput float64
	}{
		{TestStatsInput{
			NumMessagesSent:      0,
			NumSlashCommandsUsed: 0,
			NumReactionsReceived: 0,
			TsVc:                 0,
			TsMusic:              0,
			MessageWeight:        0.5,
			SlashCommandWeight:   0.45,
			ReactionsWeight:      0.33,
			TsVcWeight:           0.133,
			TsMusicWeight:        0.1,
		}, 0},
		{TestStatsInput{
			NumMessagesSent:      1,
			NumSlashCommandsUsed: 0,
			NumReactionsReceived: 0,
			TsVc:                 0,
			TsMusic:              0,
			MessageWeight:        0.5,
			SlashCommandWeight:   0.45,
			ReactionsWeight:      0.33,
			TsVcWeight:           0.133,
			TsMusicWeight:        0.1,
		}, 0.50},
		{TestStatsInput{
			NumMessagesSent:      153,
			NumSlashCommandsUsed: 10,
			NumReactionsReceived: 7,
			TsVc:                 300,
			TsMusic:              0,
			MessageWeight:        0.5,
			SlashCommandWeight:   0.45,
			ReactionsWeight:      0.33,
			TsVcWeight:           0.133,
			TsMusicWeight:        0.1,
		}, 123.21},
	}

	for _, c := range cases {
		output := utils.CalculateExperiencePointsFromStats(
			c.input.NumMessagesSent,
			c.input.NumSlashCommandsUsed,
			c.input.NumReactionsReceived,
			c.input.TsVc,
			c.input.TsMusic,
			c.input.MessageWeight,
			c.input.SlashCommandWeight,
			c.input.ReactionsWeight,
			c.input.TsVcWeight,
			c.input.TsMusicWeight,
		)

		const float64EqualityThreshold = 1e-9

		if math.Abs(c.expectedOutput-output) > float64EqualityThreshold {
			t.Errorf("difference between expected XP output and actual XP output is too big")
		}

	}

}
