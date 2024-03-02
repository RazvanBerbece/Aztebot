package utils

import "github.com/RazvanBerbece/Aztebot/internal/data/repositories"

func CleanupRepositories(rolesRepository *repositories.RolesRepository, usersRepository *repositories.UsersRepository, userStatsRepository *repositories.UsersStatsRepository, warnsRepository *repositories.WarnsRepository, timeoutsRepository *repositories.TimeoutsRepository) {

	if rolesRepository != nil {
		rolesRepository.Conn.Db.Close()
	}

	if usersRepository != nil {
		usersRepository.Conn.Db.Close()
	}

	if userStatsRepository != nil {
		userStatsRepository.Conn.Db.Close()
	}

	if warnsRepository != nil {
		warnsRepository.Conn.Db.Close()
	}

	if timeoutsRepository != nil {
		timeoutsRepository.Conn.Db.Close()
	}

}
