package utils

func CalculateExperiencePointsFromStats(
	messagesSent int,
	slashCommandsUsed int,
	reactionsReceived int,
	tsVc int,
	tsMusic int) int {

	var MessageWeight float64 = 1.0
	var SlashCommandWeight float64 = 0.65
	var ReactionsReceivedWeight float64 = 0.77
	var TimeSpentVCWeight float64 = 0.25
	var TimeSpentMusicWeight float64 = 0.125

	var totalExperience int = int(float64(messagesSent)*MessageWeight +
		float64(slashCommandsUsed)*SlashCommandWeight +
		float64(reactionsReceived)*ReactionsReceivedWeight +
		float64(tsVc)*TimeSpentVCWeight +
		float64(tsMusic)*TimeSpentMusicWeight)

	return totalExperience

}
