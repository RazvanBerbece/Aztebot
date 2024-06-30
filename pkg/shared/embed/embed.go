package embed

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Embed struct {
	*discordgo.MessageEmbed
}

// Constants for message embed character limits
const (
	EmbedLimitTitle       = 256
	EmbedLimitDescription = 2048
	EmbedLimitFieldValue  = 1024
	EmbedLimitFieldName   = 256
	EmbedLimitField       = 25
	EmbedLimitFooter      = 2048
	EmbedLimit            = 4000
)

// NewEmbed returns a new embed object
func NewEmbed() *Embed {
	return &Embed{&discordgo.MessageEmbed{}}
}

// SetTitle ...
func (e *Embed) SetTitle(name string) *Embed {
	e.Title = name
	return e
}

// SetDescription [desc]
func (e *Embed) SetDescription(description string) *Embed {
	if len(description) > 2048 {
		description = description[:2048]
	}
	e.Description = description
	return e
}

// AddField [name] [value]
func (e *Embed) AddField(name, value string, inline bool) *Embed {
	if len(value) > 1024 {
		value = value[:1024]
	}

	if len(name) > 1024 {
		name = name[:1024]
	}

	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})

	return e

}

// AddLineBreakField
func (e *Embed) AddLineBreakField() *Embed {
	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:   "\u200B",
		Value:  "",
		Inline: false,
	})

	return e
}

// AtTagEveryone
func (e *Embed) AtTagEveryone(hidden bool) *Embed {

	if hidden {
		e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
			Name:   "",
			Value:  "||@everyone||",
			Inline: false,
		})
	} else {
		e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
			Name:   "",
			Value:  "@everyone",
			Inline: false,
		})
	}

	return e
}

// SetFooter [Text] [iconURL]
func (e *Embed) SetFooter(args ...string) *Embed {
	iconURL := ""
	text := ""
	proxyURL := ""

	switch {
	case len(args) > 2:
		proxyURL = args[2]
		fallthrough
	case len(args) > 1:
		iconURL = args[1]
		fallthrough
	case len(args) > 0:
		text = args[0]
	case len(args) == 0:
		return e
	}

	e.Footer = &discordgo.MessageEmbedFooter{
		IconURL:      iconURL,
		Text:         text,
		ProxyIconURL: proxyURL,
	}

	return e
}

func (e *Embed) DecorateWithTimestampFooter(timeStringFormat string) *Embed {

	now := time.Now().Unix()

	var ts time.Time
	var timeString string

	ts = time.Unix(now, 0).UTC()
	timeString = ts.Format(timeStringFormat)

	e.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("timestamp: %s", timeString),
	}

	return e
}

// SetImage ...
func (e *Embed) SetImage(args ...string) *Embed {
	var URL string
	var proxyURL string

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		URL = args[0]
	}
	if len(args) > 1 {
		proxyURL = args[1]
	}
	e.Image = &discordgo.MessageEmbedImage{
		URL:      URL,
		ProxyURL: proxyURL,
	}
	return e
}

// SetThumbnail ...
func (e *Embed) SetThumbnail(args ...string) *Embed {
	var URL string
	var proxyURL string

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		URL = args[0]
	}
	if len(args) > 1 {
		proxyURL = args[1]
	}
	e.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:      URL,
		ProxyURL: proxyURL,
	}
	return e
}

// SetAuthor ...
func (e *Embed) SetAuthor(args ...string) *Embed {
	var (
		name     string
		iconURL  string
		URL      string
		proxyURL string
	)

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		name = args[0]
	}
	if len(args) > 1 {
		iconURL = args[1]
	}
	if len(args) > 2 {
		URL = args[2]
	}
	if len(args) > 3 {
		proxyURL = args[3]
	}

	e.Author = &discordgo.MessageEmbedAuthor{
		Name:         name,
		IconURL:      iconURL,
		URL:          URL,
		ProxyIconURL: proxyURL,
	}

	return e
}

// SetURL ...
func (e *Embed) SetURL(URL string) *Embed {
	e.URL = URL
	return e
}

// SetColor ...
func (e *Embed) SetColor(clr int) *Embed {
	e.Color = clr
	return e
}

// InlineAllFields sets all fields in the embed to be inline
func (e *Embed) InlineAllFields() *Embed {
	for _, v := range e.Fields {
		v.Inline = true
	}
	return e
}

// Truncate truncates any embed value over the character limit.
func (e *Embed) Truncate() *Embed {
	e.TruncateDescription()
	e.TruncateFields()
	e.TruncateFooter()
	e.TruncateTitle()
	return e
}

// TruncateFields truncates fields that are too long
func (e *Embed) TruncateFields() *Embed {
	if len(e.Fields) > EmbedLimitField {
		e.Fields = e.Fields[:EmbedLimitField]
	}

	for _, v := range e.Fields {

		if len(v.Name) > EmbedLimitFieldName {
			v.Name = v.Name[:EmbedLimitFieldName]
		}

		if len(v.Value) > EmbedLimitFieldValue {
			v.Value = v.Value[:EmbedLimitFieldValue]
		}

	}
	return e
}

// TruncateDescription ...
func (e *Embed) TruncateDescription() *Embed {
	if len(e.Description) > EmbedLimitDescription {
		e.Description = e.Description[:EmbedLimitDescription]
	}
	return e
}

// TruncateTitle ...
func (e *Embed) TruncateTitle() *Embed {
	if len(e.Title) > EmbedLimitTitle {
		e.Title = e.Title[:EmbedLimitTitle]
	}
	return e
}

// TruncateFooter ...
func (e *Embed) TruncateFooter() *Embed {
	if e.Footer != nil && len(e.Footer.Text) > EmbedLimitFooter {
		e.Footer.Text = e.Footer.Text[:EmbedLimitFooter]
	}
	return e
}

// GetApprovalActionRowForEmbed
func GetApprovalActionRowForEmbed(affirmativeCustomId string, negativeCustomId string) discordgo.ActionsRow {

	// Create accept and decline buttons
	acceptButton := discordgo.Button{
		Emoji: &discordgo.ComponentEmoji{
			Name: "👍🏽",
		},
		Label:    "Accept",
		Style:    discordgo.SuccessButton,
		CustomID: affirmativeCustomId,
		Disabled: false,
	}

	declineButton := discordgo.Button{
		Emoji: &discordgo.ComponentEmoji{
			Name: "👎🏽",
		},
		Label:    "Decline",
		Style:    discordgo.DangerButton,
		CustomID: negativeCustomId,
		Disabled: false,
	}

	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{&acceptButton, &declineButton},
	}

	return actionRow

}

// GetPaginationActionRowForEmbed
func GetPaginationActionRowForEmbed(previousPageCustomId string, nextPageCustomId string) discordgo.ActionsRow {

	previousPageButton := discordgo.Button{
		Emoji: &discordgo.ComponentEmoji{
			Name: "⬅️",
		},
		Label:    "Previous",
		Style:    discordgo.PrimaryButton,
		CustomID: previousPageCustomId,
		Disabled: false,
	}

	nextPageButton := discordgo.Button{
		Emoji: &discordgo.ComponentEmoji{
			Name: "➡️",
		},
		Label:    "Next",
		Style:    discordgo.PrimaryButton,
		CustomID: nextPageCustomId,
		Disabled: false,
	}

	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{&previousPageButton, &nextPageButton},
	}

	return actionRow

}
