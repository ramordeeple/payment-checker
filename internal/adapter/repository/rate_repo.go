package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"payment-checker/internal/domain"
	"time"
)

type RateRepo struct {
	db *sql.DB
}

func (r *RateRepo) GetRate(date time.Time, currency domain.CurrencyCode) (domain.Rate, error) {
	var rate domain.Rate

	query := `
        SELECT date, currency, nominal, value_scaled, cbr_id, num_code, name
        FROM rates
        WHERE date = $1 AND currency = $2
        LIMIT 1
    `

	row := r.db.QueryRow(query, date, string(currency))
	err := row.Scan(
		&rate.Date,
		&rate.Currency,
		&rate.Nominal,
		&rate.ValueScaled,
		&rate.CBRID,
		&rate.NumCode,
		&rate.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Rate{}, domain.ErrRateUnavailable
		}
		return domain.Rate{}, err
	}

	return rate, nil
}

func (r *RateRepo) GetRatesByDate(date time.Time) ([]domain.Rate, error) {
	rows, err := r.db.Query(`
		SELECT currency, nominal, value_scaled
		FROM rates
		WHERE date = $1`, date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []domain.Rate
	for rows.Next() {
		var currency string
		var nominal int32
		var valueScaled int64
		if err := rows.Scan(&currency, &nominal, &valueScaled); err != nil {
			return nil, err
		}

		rates = append(rates, domain.Rate{
			Date:        date,
			Currency:    domain.CurrencyCode(currency),
			Nominal:     nominal,
			ValueScaled: valueScaled,
		})
	}

	if len(rates) == 0 {
		return nil, domain.ErrRateNotFound
	}

	return rates, nil
}

func (r *RateRepo) GetCurrencyMeta(code domain.CurrencyCode) (domain.Currency, error) {
	var nameRU, numCode, cbrID string
	row := r.db.QueryRow(`
        SELECT name, num_code, cbr_id
        FROM currencies
        WHERE code = $1
    `, string(code))

	if err := row.Scan(&nameRU, &numCode, &cbrID); err != nil {
		if err == sql.ErrNoRows {
			return domain.Currency{}, fmt.Errorf("currency not found")
		}
		return domain.Currency{}, err
	}

	return domain.Currency{
		Code:    code,
		NameRU:  nameRU,
		NumCode: numCode,
		CBRID:   cbrID,
	}, nil
}

func NewRateRepo(database *sql.DB) *RateRepo {
	return &RateRepo{db: database}
}

func (r *RateRepo) HasCurrency(currency domain.CurrencyCode) bool {
	var exists bool
	row := r.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM rates WHERE currency = $1)`, currency)

	_ = row.Scan(&exists)

	return exists
}
