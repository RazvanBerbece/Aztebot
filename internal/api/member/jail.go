package member

import (
	"fmt"
	"time"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/data/models/events"
	"github.com/RazvanBerbece/Aztebot/internal/globals"
	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/globals/repo"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
	"github.com/bwmarrin/discordgo"
)

func JailMember(s *discordgo.Session, guildId string, userId string, reason string, jailRoleName string, notificationChannelId string) (*dataModels.JailedUser, *dataModels.User, error) {

	var err error

	// Ensure that a user won't be jailed twice
	isJailedResult := globalsRepo.JailRepository.UserIsJailed(userId)
	if isJailedResult <= 0 {
		if isJailedResult == -1 {
			return nil, nil, fmt.Errorf("could not verify whether user `%s` is already jailed", userId)
		}
	} else {
		return nil, nil, fmt.Errorf("a user cannot be jailed twice")
	}

	user, err := globalsRepo.UsersRepository.GetUser(userId)
	if err != nil {
		fmt.Printf("Failed to JailMember %s: %v", userId, err)
		return nil, nil, err
	}

	currentTimestamp := time.Now()

	// Pick a random task to assign to the jailed user
	taskToFree := utils.GetRandomFromArray(globals.JailTasks)

	// Build a record of the jailed user for the command feedback
	var jailedRecord *dataModels.JailedUser = &dataModels.JailedUser{
		UserId:         userId,
		JailedAt:       currentTimestamp.Unix(),
		TaskToComplete: taskToFree,
		Reason:         reason,
	}

	// Add User to Jail in the DB
	err = globalsRepo.JailRepository.AddUserToJail(userId, reason, taskToFree, currentTimestamp.Unix(), user.CurrentRoleIds)
	if err != nil {
		fmt.Printf("Failed to JailMember %s: %v", userId, err)
		return nil, nil, err
	}

	// Remove all roles from Discord user to restrict access
	err = RemoveAllDiscordUserRoles(s, guildId, userId)
	if err != nil {
		fmt.Printf("Failed to JailMember %s: %v", userId, err)
		return nil, nil, err
	}

	// Remove all roles from OTA member in the database
	err = RemoveAllMemberRoles(userId)
	if err != nil {
		fmt.Printf("Failed to JailMember %s: %v", userId, err)
		return nil, nil, err
	}

	// Give designated Jailed Discord role to member
	err = GiveDiscordRoleToMember(s, guildId, userId, jailRoleName)
	if err != nil {
		fmt.Printf("Failed to JailMember %s: %v", userId, err)
		return nil, nil, err
	}

	// Send notification about jailing on designated channel
	notificationEmbed := embed.NewEmbed().
		SetTitle("👮🏽‍♀️⛓️    A New Prisoner Has Arrived").
		AddField("Known As", user.DiscordTag, false).
		AddField("Convicted Because", reason, false).
		AddField("Tasked With", taskToFree, false).
		AddField("Convincted At", currentTimestamp.String(), false)

	globals.NotificationsChannel <- events.NotificationEvent{
		TargetChannelId: notificationChannelId,
		Type:            "EMBED_PASSTHROUGH",
		Embed:           notificationEmbed,
	}

	// Send Jail DM to jailed user
	dmEmbed := embed.NewEmbed().
		SetTitle("👮🏽‍♀️⛓️    You have been jailed.").
		AddField("", fmt.Sprintf("You have been jailed on: `%s`, for the following reason: `%s`.\n\nYour rights have been stripped but you can still communicate via the designated Jail channel. In order to be released from Jail, you'll need to complete the task you have been randomly assgined when you were jailed.\n\nYour assigned task is: `%s`.\n\nThe staff supervisors will guide you through the process and the implications.", currentTimestamp.String(), reason, taskToFree), false)

	go SendDirectEmbedToMember(s, userId, *dmEmbed)

	return jailedRecord, user, nil

}

func UnjailMember(s *discordgo.Session, guildId string, userId string, jailRoleName string, notificationChannelId string) (*dataModels.JailedUser, *dataModels.User, error) {

	var err error

	// Make sure that a user can't be unjailed if not in jail at a certain point in time
	isJailedResult := globalsRepo.JailRepository.UserIsJailed(userId)
	if isJailedResult <= 0 {
		return nil, nil, fmt.Errorf("cannot unjail a user who is not in jail. user `%s` not found in jail", userId)
	}

	user, err := globalsRepo.UsersRepository.GetUser(userId)
	if err != nil {
		fmt.Printf("Failed to UnjailMember (Retrieve OTA Member) %s: %v", userId, err)
		return nil, nil, err
	}

	jailedUser, err := globalsRepo.JailRepository.GetJailedUser(userId)
	if err != nil {
		fmt.Printf("Failed to UnjailMember (Retrieve Jailed Member Entry) %s: %v", userId, err)
		return nil, nil, err
	}

	// Remove User from Jail in the DB
	err = globalsRepo.JailRepository.RemoveUserFromJail(userId)
	if err != nil {
		fmt.Printf("Failed to UnjailMember (Remove user from OTA Jail) %s: %v", userId, err)
		return nil, nil, err
	}

	// Give roles back to member to return permsisions
	err = AddRolesToDiscordUser(s, guildId, userId, utils.GetRoleIdsFromRoleString(jailedUser.RoleIdsBeforeJail))
	if err != nil {
		fmt.Printf("Failed to UnjailMember (Return roles before jail) %s: %v", userId, err)
		return nil, nil, err
	}

	// Remove designated Jailed Discord role from member
	err = RemoveDiscordRoleFromMember(s, guildId, userId, jailRoleName)
	if err != nil {
		fmt.Printf("Failed to UnjailMember (Remove Jailed Role From Discord) %s: %v", userId, err)
		return nil, nil, err
	}

	// Send notification about unjailing on designated channel
	notificationEmbed := embed.NewEmbed().
		SetTitle("👮🏽‍♀️⛓️    A Prisoner Has Been Released !").
		AddField("Known As", user.DiscordTag, false).
		AddField("Convicted Because", jailedUser.Reason, false).
		AddField("Completed Release Task", jailedUser.TaskToComplete, false).
		AddField("Convincted At", utils.FormatUnixAsString(jailedUser.JailedAt, "Mon, 02 Jan 2006 15:04:05 MST"), false)

	globals.NotificationsChannel <- events.NotificationEvent{
		TargetChannelId: notificationChannelId,
		Type:            "EMBED_PASSTHROUGH",
		Embed:           notificationEmbed,
	}

	// Send Unjail DM to jailed user
	dmEmbed := embed.NewEmbed().
		SetTitle("👮🏽‍♀️⛓️    You have been unjailed.").
		AddField("", "You have been unjailed for completing your release task! Well done.", false)

	go SendDirectEmbedToMember(s, userId, *dmEmbed)

	return jailedUser, user, nil

}
