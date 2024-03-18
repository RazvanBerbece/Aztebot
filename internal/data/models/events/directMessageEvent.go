package events

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
)

type DirectMessageEvent struct {
	UserId string
	Title  *string
	Text   *string
	Embed  *embed.Embed
}
