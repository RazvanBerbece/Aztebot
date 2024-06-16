package utils

import repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"

func CleanupRepositories(rolesRepository *repositories.RolesRepository, usersRepository *repositories.UsersRepository, userStatsRepository *repositories.UsersStatsRepository, warnsRepository *repositories.WarnsRepository, timeoutsRepository *repositories.TimeoutsRepository) {

	if rolesRepository != nil {
		rolesRepository.Conn.SqlDb.Close()
	}

	if usersRepository != nil {
		usersRepository.Conn.SqlDb.Close()
	}

	if userStatsRepository != nil {
		userStatsRepository.Conn.SqlDb.Close()
	}

	if warnsRepository != nil {
		warnsRepository.Conn.SqlDb.Close()
	}

	if timeoutsRepository != nil {
		timeoutsRepository.Conn.SqlDb.Close()
	}

}
