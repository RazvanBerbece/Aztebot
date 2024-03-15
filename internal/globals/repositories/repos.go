package globalRepositories

import "github.com/RazvanBerbece/Aztebot/internal/data/repositories"

// =============== RUNTIME VARIABLES (BOT APPLICATIONS) ===============

// Database tables repositories (1 connection per repository)
var RolesRepository = repositories.NewRolesRepository()
var UsersRepository = repositories.NewUsersRepository()
var UserStatsRepository = repositories.NewUsersStatsRepository()
var WarnsRepository = repositories.NewWarnsRepository()
var TimeoutsRepository = repositories.NewTimeoutsRepository()
var MonthlyLeaderboardRepository = repositories.NewMonthlyLeaderboardRepository()
var JailRepository = repositories.NewJailRepository()
