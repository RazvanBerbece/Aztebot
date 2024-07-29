package repositories

import (
	"fmt"

	databaseconn "github.com/RazvanBerbece/Aztebot/internal/data/connection"
	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
)

type GlobalGainRatesRepository struct {
	Conn databaseconn.AztebotDbContext
}

func NewGlobalGainRatesRepository() *GlobalGainRatesRepository {
	repo := new(GlobalGainRatesRepository)
	repo.Conn.Connect()
	return repo
}

func (r GlobalGainRatesRepository) UpdateXpGlobalGainRate(activityId string, multiplier float64) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE GlobalGainRates SET 
			multiplierXp = ?
		WHERE activityId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(multiplier, activityId)
	if err != nil {
		return err
	}

	return nil
}

func (r GlobalGainRatesRepository) UpdateCoinsGlobalGainRate(activityId string, multiplier float64) error {

	stmt, err := r.Conn.SqlDb.Prepare(`
		UPDATE GlobalGainRates SET 
			multiplierCoins = ?
		WHERE activityId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(multiplier, activityId)
	if err != nil {
		return err
	}

	return nil
}

func (r GlobalGainRatesRepository) GetGlobalGainRates() ([]dax.GlobalGainRate, error) {

	var rates []dax.GlobalGainRate

	rows, err := r.Conn.SqlDb.Query("SELECT * FROM GlobalGainRates")
	if err != nil {
		return nil, fmt.Errorf("GetGlobalGainRates: %v", err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var rate dax.GlobalGainRate
		if err := rows.Scan(&rate.ActivityId, &rate.MultiplierXp, &rate.MultiplierCoins); err != nil {
			return nil, fmt.Errorf("GetGlobalGainRates: %v", err)
		}
		rates = append(rates, rate)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetGlobalGainRates: %v", err)
	}

	return rates, nil

}
