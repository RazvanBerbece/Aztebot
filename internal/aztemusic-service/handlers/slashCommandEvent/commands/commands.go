package commands

import (
	"fmt"
	"os"

	slashHandlers "github.com/RazvanBerbece/Aztebot/internal/aztemusic-service/handlers/slashCommandEvent/commands/slashHandlers"
	"github.com/bwmarrin/discordgo"
)

var appName = os.Getenv("APP_NAME")

var AztemusicSlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "music_ping",
		Description: fmt.Sprintf("Basic ping slash interaction for %s", appName),
	},
	{
		Name:        "music_join",
		Description: fmt.Sprintf("Command which instructs %s to join the designated Radio channel.", appName),
	},
	{
		Name:        "music_disconnect",
		Description: fmt.Sprintf("Command which instructs %s to disconnect from the designated Radio channel.", appName),
	},
	{
		Name:        "music_play",
		Description: fmt.Sprintf("Command which instructs %s to start playing music to the designated Radio channel.", appName),
	},
}

var AztemusicSlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"music_ping":       slashHandlers.HandleSlashMusicPing,
	"music_join":       slashHandlers.HandleSlashMusicJoin,
	"music_disconnect": slashHandlers.HandleSlashMusicDisconnect,
	"music_play":       slashHandlers.HandleSlashMusicPlay,
}
