package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
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

func (repo *PostgresRepo) CheckPassword(ctx context.Context, user *store.User) error {
	var (
		sql      = "SELECT password from USERS WHERE username = $1"
		password string
	)

	err := repo.db.QueryRow(ctx, sql, user.Username).Scan(&password)
	if err != nil {
		return errors.Wrap(err, "invalid credentials")
	}

	if !hasher.CheckPassword(user.Password, password) {
		return errors.New("password is not correct")
	}

	return nil
}

func (repo *PostgresRepo) AuthenticateUser(ctx context.Context, refresh *store.RefreshToken) error {
	err := repo.createRefreshToken(ctx, refresh.ExpiresAt, refresh.RefreshHash, refresh.UserID)
	if err != nil {
		return errors.Wrap(err, "invalid refresh token")
	}

	return nil
}

func (repo *PostgresRepo) createRefreshToken(ctx context.Context, expTime time.Time, refreshToken string, userid int64) error {
	sql := `INSERT into refresh_tokens(user_id,refresh_hash,expires_at, revoked) VALUES
	($1,$2,$3,false)`
	_, err := repo.db.Exec(ctx, sql, userid, refreshToken, expTime)
	return err
}

func (repo *PostgresRepo) GetBalance(ctx context.Context, userid int64) ([]*store.WalletBalance, error) {
	var (
		sql     = `select wallet_balances.id,wallet_balances.currency,wallet_balances.balance from wallet_balances join wallets on wallet_balances.wallet_id=wallets.id where wallets.id=$1`
		wallets []*store.WalletBalance
		wallet  store.WalletBalance
	)
	rows, err := repo.db.Query(ctx, sql, userid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query wallet balances")
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&wallet.ID, &wallet.Currency, &wallet.Balance)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		wallets = append(wallets, &wallet)
	}

	return wallets, nil
}

func (repo *PostgresRepo) UpdateBalance(ctx context.Context, newBalance *store.UpdateBalance) error {
	operator, err := repo.getOperator(newBalance.Operation)
	if err != nil {
		return err
	}

	return repo.updateBalance(
		ctx, newBalance.Amount,
		newBalance.UserID, newBalance.Currency,
		operator)
}

func (repo *PostgresRepo) getOperator(operation string) (string, error) {
	var operator string
	switch operation {
	case "deposit":
		operator = "+"
	case "withdraw":
		operator = "-"
	default:
		return "", errors.New("invalid operation")
	}
	return operator, nil
}

func (repo *PostgresRepo) updateBalance(ctx context.Context, amount float64, userid int64, currency, operator string) error {
	sql := repo.getSqlForChangeBalance(operator)

	row, err := repo.db.Exec(ctx, sql, amount, currency, userid)
	if err != nil {
		return errors.Wrap(err, "failed to update balance")
	}
	if row.RowsAffected() == 0 {
		return errors.New("user or currency not found")
	}
	return nil
}

func (repo *PostgresRepo) getSqlForChangeBalance(operator string) string {
	return `WITH wallet_ids AS (
    SELECT id
    FROM wallets
    WHERE user_id = $3
)
UPDATE wallet_balances
SET balance = balance ` + operator + ` $1
WHERE currency = $2
AND wallet_id IN (SELECT id FROM wallet_ids);`
}

func (repo *PostgresRepo) ExchangeCurrency(ctx context.Context, exchangeBody *store.ExchangeBalance) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transactionw")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	err = repo.makeExchange(ctx, tx, exchangeBody)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (repo *PostgresRepo) makeExchange(ctx context.Context, tx pgx.Tx, exchangeBody *store.ExchangeBalance) error {
	_, err := tx.Exec(ctx, repo.getSqlForChangeBalance("-"), exchangeBody.FromAmount, exchangeBody.FromCurrency, exchangeBody.UserID)
	if err != nil {
		return errors.Wrapf(err, "failed to update %s balance", exchangeBody.FromCurrency)
	}
	_, err = tx.Exec(ctx, repo.getSqlForChangeBalance("+"), exchangeBody.ToAmount, exchangeBody.ToCurrency, exchangeBody.UserID)
	if err != nil {
		return errors.Wrap(err, "failed to update target currency balance")
	}
	return nil
}
