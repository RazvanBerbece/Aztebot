package member

import (
	"fmt"

	globalsRepo "github.com/RazvanBerbece/Aztebot/internal/globals/repo"
	"github.com/bwmarrin/discordgo"
)

func KickMember(s *discordgo.Session, guildId string, userId string) error {
	// Delete member from server
	err := s.GuildMemberDelete(guildId, userId)
	if err != nil {
		fmt.Println("Error kicking member from guild:", err)
		return err
	}
	// Delete member-related entries from the databases
	err = DeleteAllMemberData(userId)
	if err != nil {
		fmt.Printf("Error deleting member %s data from DB tables: %v", userId, err)
		return err
	}
	return nil
}

func DeleteAllMemberData(userId string) error {
	err := globalsRepo.UserStatsRepository.DeleteUserStats(userId)
	if err != nil {
		fmt.Printf("Error deleting member %s stats from DB: %v", userId, err)
		return err
	}
	err = globalsRepo.UsersRepository.DeleteUser(userId)
	if err != nil {
		fmt.Printf("Error deleting user %s from DB: %v", userId, err)
		return err
	}
	err = globalsRepo.WarnsRepository.DeleteAllWarningsForUser(userId)
	if err != nil {
		fmt.Printf("Error deleting user %s warnings from DB: %v", userId, err)
		return err
	}
	err = globalsRepo.TimeoutsRepository.ClearTimeoutForUser(userId)
	if err != nil {
		fmt.Printf("Error deleting user %s active timeouts from DB: %v", userId, err)
		return err
	}
	err = globalsRepo.TimeoutsRepository.ClearArchivedTimeoutsForUser(userId)
	if err != nil {
		fmt.Printf("Error deleting user %s archived timeouts from DB: %v", userId, err)
		return err
	}
	err = globalsRepo.MonthlyLeaderboardRepository.DeleteEntry(userId)
	if err != nil {
		fmt.Printf("Error deleting user %s monthly leaderboard entry from DB: %v", userId, err)
		return err
	}
	err = globalsRepo.JailRepository.RemoveUserFromJail(userId)
	if err != nil {
		fmt.Printf("Error deleting user %s jail entry from DB: %v", userId, err)
		return err
	}

	return nil
}
