package member

import (
	"database/sql"
	"fmt"
	"time"

	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
	globalRepositories "github.com/RazvanBerbece/Aztebot/internal/globals/repositories"
	"github.com/bwmarrin/discordgo"
)

func GetMemberTimeouts(userId string) (*dax.Timeout, []dax.ArchivedTimeout, error) {

	// Result variables
	var activeTimeoutResult *dax.Timeout = nil
	var archivedTimeoutResults []dax.ArchivedTimeout = []dax.ArchivedTimeout{}

	// Active timeout
	activeTimeout, err := globalRepositories.TimeoutsRepository.GetUserTimeout(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			activeTimeoutResult = nil
		} else {
			return nil, nil, err
		}
	}
	activeTimeoutResult = activeTimeout

	// Archived timeouts
	archivedTimeoutResults, err = globalRepositories.TimeoutsRepository.GetAllArchivedTimeoutsForUser(userId)
	if err != nil {
		return nil, nil, err
	}

	return activeTimeoutResult, archivedTimeoutResults, nil

}

func GiveTimeoutToMemberWithId(s *discordgo.Session, guildId string, userId string, reason string, creationTimestamp int64, sTimeoutLength float64) error {

	result := globalRepositories.TimeoutsRepository.GetTimeoutsCountForUser(userId)
	if result > 0 {
		return fmt.Errorf("a user cannot be given more than 1 timeout at a time")
	}

	// If the user is on their 10th timeout
	numArchivedTimeouts := globalRepositories.TimeoutsRepository.GetArchivedTimeoutsCountForUser(userId)
	if numArchivedTimeouts == 9 {
		// ban them instead
		err := s.GuildBanCreateWithReason(guildId, userId, "Received 10th and final timeout", 1)
		if err != nil {
			fmt.Println("Error banning user on 10th timeout: ", err)
			return err
		}
		// and clean DB related entries
		err = DeleteAllMemberData(userId)
		if err != nil {
			fmt.Println("Error deleting user data on 10th timeout: ", err)
			return err
		}
	}

	err := globalRepositories.TimeoutsRepository.SaveTimeout(userId, reason, creationTimestamp, int(sTimeoutLength))
	if err != nil {
		fmt.Printf("Error ocurred while storing timeout for user: %s\n", err)
		return fmt.Errorf(err.Error())
	}

	// Give actual Discord timeout to member
	timeoutExpiryTimestamp := time.Now().Add(time.Second * time.Duration(sTimeoutLength))
	err = s.GuildMemberTimeout(guildId, userId, &timeoutExpiryTimestamp)
	if err != nil {
		fmt.Println("Error timing out user: ", err)
		return fmt.Errorf("%v", err)
	}

	return nil

}

func ClearMemberActiveTimeout(s *discordgo.Session, guildId string, userId string) error {

	err := globalRepositories.TimeoutsRepository.ClearTimeoutForUser(userId)
	if err != nil {
		return err
	}

	err = s.GuildMemberTimeout(guildId, userId, nil)
	if err != nil {
		fmt.Println("Error timing out user: ", err)
		return fmt.Errorf("%v", err)
	}

	return nil

}

// TODO: Implement timeout appeals in private DMs between bot application and guild member.
func AppealTimeout(guildId string, userId string) error {

	activeTimeout, _, err := GetMemberTimeouts(userId)
	if err != nil {
		timeoutError := fmt.Errorf("an error ocurred while retrieving timeout data for user with ID %s: %v", userId, err)
		return timeoutError
	}

	if activeTimeout == nil {
		return fmt.Errorf("no active timeout was found for user with ID `%s`", userId)
	}

	// TODO: Etc etc.

	return nil

}
