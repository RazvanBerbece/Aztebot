package globalRepositories

import (
	aztebotRepositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
	aztemarketRepositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztemarket"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
)

// =============== RUNTIME VARIABLES (BOT APPLICATIONS) ===============

// Database tables repositories (1 connection per repository)
var RolesRepository = aztebotRepositories.NewRolesRepository()
var UsersRepository = aztebotRepositories.NewUsersRepository()
var UserStatsRepository = aztebotRepositories.NewUsersStatsRepository()
var WarnsRepository = aztebotRepositories.NewWarnsRepository()
var TimeoutsRepository = aztebotRepositories.NewTimeoutsRepository()
var MonthlyLeaderboardRepository = aztebotRepositories.NewMonthlyLeaderboardRepository()
var JailRepository = aztebotRepositories.NewJailRepository()
var ArcadeLadderRepository = aztebotRepositories.NewArcadeLadderRepository()
var UserRepRepository = aztebotRepositories.NewUserRepRepository()

var WalletsRepository = aztemarketRepositories.NewWalletsRepository(globalConfiguration.MySqlAztemarketRootConnectionString)
