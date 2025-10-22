package repository

import (
	"database/sql"
	"payment-checker/internal/domain"
	"time"
)

type RateRepo struct {
	db *sql.DB
}

func NewRateRepo(database *sql.DB) *RateRepo {
	return &RateRepo{db: database}
}

func (r *RateRepo) GetRate(date time.Time, currency domain.CurrencyCode) (domain.Rate, error) {
	var nominal int32
	var valueScaled int64

	row := r.db.QueryRow(`
    SELECT nominal, value_scaled
    FROM rates
    WHERE date = $1 AND currency = $2`,
		date.Format("2006-01-02"), string(currency))

	if err := row.Scan(&nominal, &valueScaled); err != nil {
		if err == sql.ErrNoRows {
			return domain.Rate{}, domain.ErrRateNotFound
		}
		return domain.Rate{}, err
	}

	return domain.Rate{
		Date:        date,
		Currency:    currency,
		Nominal:     nominal,
		ValueScaled: valueScaled,
	}, nil
}

func (r *RateRepo) HasCurrency(currency domain.CurrencyCode) bool {
	var exists bool
	row := r.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM rates WHERE currency = $1)`, currency)

	_ = row.Scan(&exists)

	return exists
}
