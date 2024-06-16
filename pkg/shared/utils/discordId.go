package utils

import "strings"

func GetDiscordIdFromMentionFormat(mention string) string {
	return strings.Trim(mention, "<@!>")
}
