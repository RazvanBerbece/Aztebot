package globalsRepo

import "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"

// =============== RUNTIME VARIABLES (BOT APPLICATIONS) ===============

// Database tables repositories (1 connection per repository)
var RolesRepository = repositories.NewRolesRepository()
var UsersRepository = repositories.NewUsersRepository()
var UserStatsRepository = repositories.NewUsersStatsRepository()
var WarnsRepository = repositories.NewWarnsRepository()
