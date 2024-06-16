package utils

func CalculateExperiencePointsFromStats(
	messagesSent int,
	slashCommandsUsed int,
	reactionsReceived int,
	tsVc int,
	tsMusic int,
	messageWeight float64,
	slashCommandWeight float64,
	reactionsReceivedWeight float64,
	timeSpentVCWeight float64,
	timeSpentMusicWeight float64,
) int {
	var totalExperience int = int(float64(messagesSent)*messageWeight +
		float64(slashCommandsUsed)*slashCommandWeight +
		float64(reactionsReceived)*reactionsReceivedWeight +
		float64(tsVc)*timeSpentVCWeight +
		float64(tsMusic)*timeSpentMusicWeight)

	return totalExperience
}
