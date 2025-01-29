package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/store"
	logger "github.com/mizmorr/loggerm"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func setupMockRepo(t *testing.T) (*PostgresRepo, pgxmock.PgxPoolIface) {
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()
	logger := logger.Get("test.log", "debug")
	repo := &PostgresRepo{db: mockDB, log: logger}
	defer mockDB.Close()
	return repo, mockDB
}

func TestCreateUser(t *testing.T) {
	testRepo, mockDB := setupMockRepo(t)

	ctx := context.Background()
	user := &store.User{Username: "testcreateuser", Email: "test@example.com", Password: "securepassword"}

	mockDB.ExpectQuery(`INSERT INTO users`).
		WithArgs(user.Username, user.Email, pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(1)))

	mockDB.ExpectQuery(`INSERT INTO wallets`).
		WithArgs(int64(1)).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(1001)))

	mockDB.ExpectExec(`INSERT INTO wallet_balances`).
		WithArgs(int64(1001)).
		WillReturnResult(pgxmock.NewResult("INSERT", 3))

	err := testRepo.CreateUser(ctx, user)

	assert.NoError(t, err)

	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestSetToken(t *testing.T) {
	repo, mockDB := setupMockRepo(t)
	ctx := context.Background()
	refresh := &store.RefreshToken{UserID: 1, Hash: "somehash", ExpiresAt: time.Now().Add(time.Hour)}

	mockDB.ExpectExec("INSERT into refresh_tokens").
		WithArgs(refresh.UserID, refresh.Hash, refresh.ExpiresAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.SetToken(ctx, refresh)
	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestGetBalance(t *testing.T) {
	repo, mockDB := setupMockRepo(t)
	ctx := context.Background()

	mockDB.ExpectQuery("SELECT wb.id, wb.currency, wb.balance").
		WithArgs(int64(1)).
		WillReturnRows(pgxmock.NewRows([]string{"id", "currency", "balance"}).
			AddRow(int64(100), "USD", 500.0).
			AddRow(int64(101), "EUR", 300.0))

	balance, err := repo.GetBalance(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, balance, 2)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestCheckRefreshToken(t *testing.T) {
	repo, mockDB := setupMockRepo(t)
	ctx := context.Background()
	token := &store.RefreshToken{UserID: 1, Hash: "somehash"}

	mockDB.ExpectQuery("SELECT 1 FROM refresh_tokens").
		WithArgs(token.UserID, token.Hash).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(1))

	err := repo.CheckRefreshToken(ctx, token)
	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}
