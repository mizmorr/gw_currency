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

	repo.log.Info().Str("username", user.Username).Str("email", user.Email).Msg("Starting user creation")

	hashedPassword, err := hasher.MakeHash(user.Password)
	if err != nil {
		repo.log.Error().Err(err).Msg("Failed to hash password")
		return err
	}

	repo.log.Info().Msg(hashedPassword)

	err = repo.db.QueryRow(ctx, sqlCreateUser, user.Username, user.Email, hashedPassword).Scan(&userID)
	if err != nil {
		repo.log.Error().Err(err).Str("username", user.Username).Msg("Failed to insert user into database")
		return err
	}
	repo.log.Info().Int64("userID", userID).Msg("User created successfully")

	err = repo.db.QueryRow(ctx, sqlCreateWallet, userID).Scan(&walletID)
	if err != nil {
		repo.log.Error().Err(err).Int64("userID", userID).Msg("Failed to create wallet")
		return err
	}
	repo.log.Info().Int64("walletID", walletID).Msg("Wallet created successfully")

	row, err := repo.db.Exec(ctx, sqlSetBalances, walletID)
	if err != nil {
		repo.log.Error().Err(err).Int64("walletID", walletID).Msg("Failed to set initial wallet balances")
		return err
	}
	if row.RowsAffected() == 0 {
		repo.log.Warn().Int64("walletID", walletID).Msg("No balances were set, possible issue with wallet balance initialization")
		return errors.New("setting wallet balance failed")
	}

	repo.log.Info().Int64("userID", userID).Int64("walletID", walletID).Msg("User and wallet setup completed successfully")
	return nil
}

func (repo *PostgresRepo) Authentication(ctx context.Context, user *store.User) (int64, error) {
	var (
		sql      = "SELECT id,password from USERS WHERE username = $1"
		password string
		userID   int64
	)

	repo.log.Info().Str("username", user.Username).Msg("Starting authentication")

	err := repo.db.QueryRow(ctx, sql, user.Username).Scan(&userID, &password)
	if err != nil {
		repo.log.Warn().Str("username", user.Username).Err(err).Msg("Invalid credentials")
		return 0, errors.Wrap(err, "invalid credentials")
	}

	repo.log.Debug().Int64("userID", userID).Msg("User found in database")

	if !hasher.CheckPassword(user.Password, password) {
		repo.log.Warn().Str("username", user.Username).Msg("Incorrect password")
		return 0, errors.New("password is not correct")
	}

	repo.log.Info().Int64("userID", userID).Msg("Authentication successful")
	return userID, nil
}

func (repo *PostgresRepo) SetToken(ctx context.Context, refresh *store.RefreshToken) error {
	repo.log.Info().Int64("userID", refresh.UserID).Msg("Setting refresh token")
	err := repo.createRefreshToken(ctx, refresh.ExpiresAt, refresh.Hash, refresh.UserID)
	if err != nil {
		repo.log.Error().Err(err).Msg("Failed to set refresh token")
		return errors.Wrap(err, "invalid refresh token")
	}
	return nil
}

func (repo *PostgresRepo) createRefreshToken(ctx context.Context, expTime time.Time, refreshToken string, userid int64) error {
	repo.log.Debug().Int64("userID", userid).Msg("Creating refresh token")
	sql := `INSERT into refresh_tokens(user_id,token_hash,expires_at, revoked) VALUES
	($1,$2,$3,false)`
	_, err := repo.db.Exec(ctx, sql, userid, refreshToken, expTime)
	if err != nil {
		repo.log.Error().Err(err).Msg("Failed to create refresh token")
	}
	return err
}

func (repo *PostgresRepo) GetSpecificCurrency(ctx context.Context, req *store.CurrencyRequest) (*store.WalletCurrency, error) {
	repo.log.Info().Int64("userID", req.UserID).Str("currency", req.CurrencyCode).Msg("Fetching specific currency")
	var (
		sql = `SELECT
    	wb.id,
    	wb.currency,
    	wb.balance
		FROM
    	wallet_balances wb
		INNER JOIN
    	wallets w
		ON
    	wb.wallet_id = w.id
		WHERE
    	w.user_id = $1 AND
    	wb.currency = $2;`
		balance store.WalletCurrency
	)
	err := repo.db.QueryRow(ctx, sql, req.UserID, req.CurrencyCode).Scan(&balance.ID, &balance.Currency, &balance.Balance)
	if err == pgx.ErrNoRows {
		repo.log.Warn().Msg("Currency not found")
		return nil, errors.New("currency not found")
	} else if err != nil {
		repo.log.Error().Err(err).Msg("Failed to query currency")
		return nil, errors.Wrap(err, "failed to query currency")
	}
	repo.log.Debug().Interface("balance", balance).Msg("Currency retrieved successfully")
	return &balance, nil
}

func (repo *PostgresRepo) GetBalance(ctx context.Context, userid int64) ([]*store.WalletCurrency, error) {
	repo.log.Info().Int64("userID", userid).Msg("Fetching wallet balance")

	var (
		sql = `SELECT
    	wb.id,
    	wb.currency,
    	wb.balance
		FROM
    	wallet_balances wb
		INNER JOIN
    	wallets w
		ON
    	wb.wallet_id = w.id
		WHERE
    	w.user_id = $1;`
		balance []*store.WalletCurrency
	)
	rows, err := repo.db.Query(ctx, sql, userid)
	if err != nil {
		repo.log.Error().Err(err).Msg("Failed to query wallet balances")
		return nil, errors.Wrap(err, "failed to query wallet balances")
	}
	defer rows.Close()

	for rows.Next() {
		var currency store.WalletCurrency
		err = rows.Scan(&currency.ID, &currency.Currency, &currency.Balance)
		if err != nil {
			repo.log.Error().Err(err).Msg("Failed to scan row")
			return nil, errors.Wrap(err, "failed to scan row")
		}
		balance = append(balance, &currency)
	}
	repo.log.Debug().Int("count", len(balance)).Msg("Wallet balance fetched successfully")
	return balance, nil
}

func (repo *PostgresRepo) UpdateBalance(ctx context.Context, newBalance *store.UpdateBalance) error {
	operator, err := repo.getOperator(newBalance.Operation)
	repo.log.Info().Int64("userID", newBalance.UserID).Str("currency", newBalance.Currency).Str("operation", newBalance.Operation).Msg("Updating balance")

	if err != nil {

		repo.log.Error().Err(err).Msg("Invalid operation")
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
	repo.log.Info().Int64("userID", exchangeBody.UserID).Str("from", exchangeBody.FromCurrency).Str("to", exchangeBody.ToCurrency).Msg("Starting currency exchange")
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		repo.log.Error().Err(err).Msg("Failed to begin transaction")
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			repo.log.Error().Err(err).Msg("Rolling back transaction")
			_ = tx.Rollback(ctx)
		}
	}()
	err = repo.makeExchange(ctx, tx, exchangeBody)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		repo.log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.Wrap(err, "failed to commit transaction")
	}
	repo.log.Info().Msg("Currency exchange completed successfully")
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

func (repo *PostgresRepo) CheckRefreshToken(ctx context.Context, token *store.RefreshToken) error {
	repo.log.Debug().Int64("user_id", token.UserID).Str("token_hash", token.Hash).Msg("Checking refresh token")

	sql := `SELECT 1
FROM refresh_tokens
WHERE user_id = $1
  AND token_hash = $2
  AND expires_at > NOW()
LIMIT 1;`

	var exists int
	err := repo.db.QueryRow(ctx, sql, token.UserID, token.Hash).Scan(&exists)
	if err == pgx.ErrNoRows {
		repo.log.Warn().Int64("user_id", token.UserID).Msg("Refresh token not found or expired")
		return errors.New("token not found")
	} else if err != nil {
		repo.log.Error().Err(err).Int64("user_id", token.UserID).Msg("Failed to check refresh token")
		return errors.Wrap(err, "failed to check refresh token")
	}

	if exists == 0 {
		repo.log.Warn().Int64("user_id", token.UserID).Msg("Token has been revoked")
		return errors.New("token has been revoked")
	}

	repo.log.Info().Int64("user_id", token.UserID).Msg("Refresh token is valid")
	return nil
}

func (repo *PostgresRepo) handleRefreshTokenRepsonse(err error, exists int) error {
	if err == pgx.ErrNoRows {
		return errors.New("Token not found")
	}
	if err != nil {
		return errors.Wrap(err, "failed to check refresh token")
	}
	if exists == 0 {
		return errors.New("Token has been revoked")
	}
	return nil
}
