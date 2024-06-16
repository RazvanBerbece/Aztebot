package events

import (
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
)

type DirectMessageEvent struct {
	UserId string
	Embed  *embed.Embed
	Title  *string
	Text   *string
}
