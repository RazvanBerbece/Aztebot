package member

import (
	"fmt"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	globalState "github.com/RazvanBerbece/Aztebot/internal/globals/state"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/dm"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func SendDirectMessageToMember(s *discordgo.Session, userId string, msg string) error {
	errDm := dm.DmUser(s, userId, msg)
	if errDm != nil {
		fmt.Printf("Error sending DM to member with UID %s: %v\n", userId, errDm)
		return errDm
	}
	return nil
}

func SendDirectSimpleEmbedToMember(s *discordgo.Session, userId string, title string, text string) error {

	simpleEmbed := utils.SimpleEmbed(title, text)

	errDm := dm.DmEmbedUser(s, userId, *simpleEmbed[0])
	if errDm != nil {
		fmt.Printf("Error sending embed DM to member with UID %s: %v\n", userId, errDm)
		return errDm
	}
	return nil
}

func SendDirectEmbedToMember(s *discordgo.Session, userId string, embed embed.Embed) error {
	errDm := dm.DmEmbedUser(s, userId, *embed.MessageEmbed)
	if errDm != nil {
		fmt.Printf("Error sending embed DM to member with UID %s: %v\n", userId, errDm)
		return errDm
	}
	return nil
}

func SendDirectComplexEmbedToMember(s *discordgo.Session, userId string, embed embed.Embed, actionsRow discordgo.ActionsRow, pageSize int) error {

	originalAllFields := make([]*discordgo.MessageEmbedField, len(embed.Fields))
	copy(originalAllFields, embed.Fields)

	// Only show fields from page 1 in the beginning
	embed.Fields = embed.Fields[0:pageSize]
	msg, err := dm.DmEmbedComplexUser(s, userId, *embed.MessageEmbed, actionsRow)
	if err != nil {
		fmt.Printf("Error sending embed DM to member with UID %s: %v\n", userId, err)
		return err
	}

	// Keep paginated embeds in-memory to enable handling on button presses
	// TODO: Clean the map once in a while !!
	globalState.EmbedsToPaginate[msg.ID] = dataModels.EmbedData{
		FieldData:   &originalAllFields, // all fields
		CurrentPage: 1,                  // same for all complex paginated embeds
	}

	return nil
}
