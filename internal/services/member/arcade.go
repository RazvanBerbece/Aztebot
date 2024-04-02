package member

import (
	"fmt"
	"time"

	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	globalMessaging "github.com/RazvanBerbece/Aztebot/internal/globals/messaging"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

func GiveArcadeWin(userId string, arcadeName string, notificationChannelId string) error {

	entryExists := globalRepositories.ArcadeLadderRepository.EntryExists(userId)
	if entryExists < 1 {
		if entryExists == -1 {
			return fmt.Errorf("an error ocurred while checking if arcade ladder entry exists")
		}
		err := globalRepositories.ArcadeLadderRepository.AddNewLadderEntry(userId)
		if err != nil {
			return err
		}
	}

	timestamp := time.Now().Unix()

	err := globalRepositories.ArcadeLadderRepository.AddWin(userId)
	if err != nil {
		return err
	}

	// Also send notification to announce winner in the designated channel
	user, err := globalRepositories.UsersRepository.GetUser(userId)
	if err != nil {
		return err
	}

	notificationEmbed := embed.NewEmbed().
		SetTitle("👾🎮    Arcade Winner Announcement").
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		AddField("", fmt.Sprintf("`%s` has won the `%s` event !", user.DiscordTag, arcadeName), false).
		AddField("", fmt.Sprintf("on `%s`", utils.FormatUnixAsString(timestamp, "Mon, 02 Jan 2006")), false).
		AtTagEveryone()

	globalMessaging.NotificationsChannel <- events.NotificationEvent{
		TargetChannelId: notificationChannelId,
		Type:            "EMBED_PASSTHROUGH",
		Embed:           notificationEmbed,
	}

	return nil

}
