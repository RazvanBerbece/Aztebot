package utils

import "regexp"

func IsValidDiscordUserId(userId string) bool {

	pattern := regexp.MustCompile(`^\d{17,20}$`)

	// Check if the userID matches the regular expression pattern.
	return pattern.MatchString(userId)

}

func IsValidReasonMessage(msg string) bool {

	return len(msg) < 500

}
