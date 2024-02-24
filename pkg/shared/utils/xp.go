package utils

func CalculateExperiencePointsFromStats(
	messagesSent int,
	slashCommandsUsed int,
	reactionsReceived int,
	tsVc int,
	tsMusic int) int {

	var MessageWeight float64 = 0.5
	var SlashCommandWeight float64 = 0.45
	var ReactionsReceivedWeight float64 = 0.33
	var TimeSpentVCWeight float64 = 0.133
	var TimeSpentMusicWeight float64 = 0.1

	var totalExperience int = int(float64(messagesSent)*MessageWeight +
		float64(slashCommandsUsed)*SlashCommandWeight +
		float64(reactionsReceived)*ReactionsReceivedWeight +
		float64(tsVc)*TimeSpentVCWeight +
		float64(tsMusic)*TimeSpentMusicWeight)

	return totalExperience

}
