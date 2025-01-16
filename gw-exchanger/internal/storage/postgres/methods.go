package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mizmorr/gw_currency/gw-exchanger/internal/storage/model"
	"github.com/pkg/errors"
)

func (repo *PostgresRepo) GetRate(ctx context.Context, code string) (*model.Rate, error) {
	sql := `SELECT currency_code, rate FROM exchange_rate WHERE currency_code = $1`
	row := repo.db.QueryRow(ctx, sql, code)

	var rate model.Rate
	err := row.Scan(&rate.CurrencyCode, &rate.Value)
	if err == pgx.ErrNoRows {
		return nil, errors.Wrapf(err, "rate not found for code %s", code)
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to scan row")
	}

	return &rate, nil
}

func (repo *PostgresRepo) GetAllRates(ctx context.Context) ([]*model.Rate, error) {
	sql := `SELECT currency_code, rate FROM exchange_rate`
	rows, err := repo.db.Query(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query rates")
	}
	defer rows.Close()

	var rates []*model.Rate
	for rows.Next() {
		var rate model.Rate
		err := rows.Scan(&rate.CurrencyCode, &rate.Value)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		rates = append(rates, &rate)
	}

	return rates, nil
}

func (repo *PostgresRepo) setRate(ctx context.Context, rate *model.Rate) error {
	sql := "update exchange_rate set rate=$1 where currency_code = $2"
	row, err := repo.db.Exec(ctx, sql, rate.Value, rate.CurrencyCode)
	if err != nil {
		return errors.Wrap(err, "failed to set rate")
	}
	if row.RowsAffected() == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
