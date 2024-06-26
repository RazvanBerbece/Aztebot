package testData

import repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"

func RemoveUser(r repositories.TimeoutsRepository, userId string) {

	query := "DELETE FROM Users WHERE userId = ?"

	_, _ = r.Conn.SqlDb.Exec(query, userId)

}

func RemoveUserStats(r repositories.TimeoutsRepository, userId string) {

	query := "DELETE FROM UserStats WHERE userId = ?"

	_, _ = r.Conn.SqlDb.Exec(query, userId)

}

func RemoveUserWarns(r repositories.TimeoutsRepository, userId string) {

	query := "DELETE FROM Warns WHERE userId = ?"

	_, _ = r.Conn.SqlDb.Exec(query, userId)

}

func RemoveUserArchivedTimeouts(r repositories.TimeoutsRepository, userId string) {

	query := "DELETE FROM TimeoutsArchive WHERE userId = ?"

	_, _ = r.Conn.SqlDb.Exec(query, userId)

}

func RemoveUserTimeout(r repositories.TimeoutsRepository, userId string) {

	query := "DELETE FROM Timeouts WHERE userId = ?"

	_, _ = r.Conn.SqlDb.Exec(query, userId)

}
