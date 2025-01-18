package postgres

import (
	"context"

	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/store"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/pkg/hasher"
	"github.com/pkg/errors"
)

func (repo *PostgresRepo) CreateUser(ctx context.Context, user *store.User) error {
	var (
		userID          int64
		walletID        int64
		sqlCreateUser   = `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) returning id`
		sqlCreateWallet = `INSERT INTO wallets (user_id) VALUES ($1) returning id`
		sqlSetBalances  = `INSERT INTO wallet_balances (wallet_id,currency,balance)
            VALUES
            ($1, 'RUB', 0),
            ($1, 'EUR', 0),
            ($1, 'USD', 0)`
	)

	hashedPassword, err := hasher.MakeHash(user.Password)
	if err != nil {
		return err
	}

	err = repo.db.QueryRow(ctx, sqlCreateUser, user.Username, user.Email, hashedPassword).Scan(&userID)
	if err != nil {
		return err
	}

	err = repo.db.QueryRow(ctx, sqlCreateWallet, userID).Scan(&walletID)
	if err != nil {
		return err
	}
	row, err := repo.db.Exec(ctx, sqlSetBalances, walletID)
	if err != nil {
		return err
	}
	if row.RowsAffected() == 0 {
		return errors.New("setting wallet balance faild")
	}

	return nil
}
