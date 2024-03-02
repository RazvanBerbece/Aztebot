package testData

import (
	"fmt"
	"time"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/data/repositories"
	"github.com/brianvoe/gofakeit/v6"
)

func AddUser(r repositories.TimeoutsRepository, userId string) dataModels.User {

	stmt, _ := r.Conn.Db.Prepare(`
		INSERT INTO 
			Users(
				discordTag, 
				userId, 
				currentRoleIds, 
				currentCircle, 
				currentInnerOrder, 
				currentLevel, 
				currentExperience, 
				createdAt
			)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?);`)

	defer stmt.Close()

	fakeTag := gofakeit.Word()
	fakeRoleIds := "1,"
	now := time.Now().Unix()
	circle := "OUTER"
	_, _ = stmt.Exec(fakeTag, userId, fakeRoleIds, circle, nil, 0, 0, now)

	return dataModels.User{
		Id:                gofakeit.IntRange(100000, 200000),
		UserId:            userId,
		CurrentRoleIds:    fakeRoleIds,
		CurrentCircle:     circle,
		CurrentInnerOrder: nil,
		CurrentLevel:      0,
		CurrentExperience: 0,
		CreatedAt:         &now,
	}

}

func AddTimeoutForUser(r repositories.TimeoutsRepository, t *dataModels.Timeout) (*int64, error) {

	stmt, _ := r.Conn.Db.Prepare(`
		INSERT INTO 
			Timeouts(
				userId, 
				reason, 
				creationTimestamp,
				sTimeLength
			)
		VALUES(?, ?, ?, ?);`)

	defer stmt.Close()

	result, _ := stmt.Exec(t.UserId, t.Reason, t.CreationTimestamp, t.SDuration)

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	fmt.Println(id)

	return &id, nil

}
