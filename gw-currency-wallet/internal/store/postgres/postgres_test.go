package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/store"

	"github.com/stretchr/testify/assert"
)

func getRepo(ctx context.Context) (*PostgresRepo, error) {
	repo, err := NewPostgresRepo(ctx)
	if err != nil {
		return nil, err
	}

	err = repo.Start(ctx)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func TestCreateUser(t *testing.T) {
	testStartTime := time.Now()

	ctx := context.Background()

	repo, err := getRepo(ctx)
	assert.NoError(t, err)

	defer repo.Stop(ctx)

	user := &store.User{Username: "testcreateuser", Email: "test@example.com", Password: "securepassword", CreatedAt: testStartTime}

	err = repo.CreateUser(ctx, user)

	assert.NoError(t, err)

	if err == nil {
		_, err = repo.db.Exec(ctx, "DELETE FROM users WHERE created_at > $1", testStartTime)
		assert.NoError(t, err)
	}
}

func TestAuthentication(t *testing.T) {
	ctx := context.Background()

	testStartTime := time.Now()

	repo, err := getRepo(ctx)
	assert.NoError(t, err)
	defer repo.Stop(ctx)

	// Тестовый пользователь
	user := &store.User{Username: "testuserauth", Password: "correctpassword", CreatedAt: testStartTime}

	// Создание тестового пользователя в базе (если нужно)
	err = repo.CreateUser(ctx, user)
	assert.NoError(t, err)

	// Тестирование успешной аутентификации
	userID, err := repo.Authentication(ctx, user)
	assert.NoError(t, err)
	assert.NotZero(t, userID)
	if err == nil {
		_, err = repo.db.Exec(ctx, "DELETE FROM users WHERE created_at > $1", testStartTime)
		assert.NoError(t, err)
	}

	// Тестирование аутентификации с неверным паролем
	invalidUser := &store.User{Username: "testuserauth", Password: "wrongpassword"}
	_, err = repo.Authentication(ctx, invalidUser)
	assert.Error(t, err)

	// Тестирование аутентификации с несуществующим пользователем
	nonExistentUser := &store.User{Username: "nonexistentuserauth", Password: "any"}
	_, err = repo.Authentication(ctx, nonExistentUser)
	assert.Error(t, err)
}

func TestSetToken(t *testing.T) {
	ctx := context.Background()
	repo, err := getRepo(ctx)
	testStartTime := time.Now()
	assert.NoError(t, err)
	defer repo.Stop(ctx)

	// Создание тестового refresh токена
	refreshToken := &store.RefreshToken{
		UserID:    1,
		Hash:      "samplehash",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: testStartTime,
	}

	// Тестирование успешного создания refresh токена
	err = repo.SetToken(ctx, refreshToken)
	assert.NoError(t, err)
	if err == nil {
		_, err = repo.db.Exec(ctx, "DELETE FROM refresh_tokens WHERE created_at > $1", testStartTime)
		assert.NoError(t, err)
	}
	// Тестирование ошибки при создании токена
	invalidToken := &store.RefreshToken{
		UserID:    0, // Неверный userID
		Hash:      "invalidhash",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Истекший токен
	}
	err = repo.SetToken(ctx, invalidToken)
	assert.Error(t, err)
}

func TestCheckRefreshToken(t *testing.T) {
	ctx := context.Background()
	repo, err := getRepo(ctx)
	testStartTime := time.Now()
	assert.NoError(t, err)
	defer repo.Stop(ctx)

	refreshToken := &store.RefreshToken{
		UserID:    2,
		Hash:      "validtokenhash",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: testStartTime,
	}

	// Тестирование успешного создания refresh токена
	err = repo.SetToken(ctx, refreshToken)
	assert.NoError(t, err)
	defer func() {
		if err == nil {
			_, err = repo.db.Exec(ctx, "DELETE FROM refresh_tokens WHERE created_at > $1", testStartTime)
			assert.NoError(t, err)
		}
	}()

	// Создание валидного refresh токена
	validToken := &store.RefreshToken{
		UserID: 2,
		Hash:   "validtokenhash",
	}

	// Проверка наличия токена в базе
	err = repo.CheckRefreshToken(ctx, validToken)
	assert.NoError(t, err)

	// Тестирование случая, когда токен не найден
	invalidToken := &store.RefreshToken{
		UserID: 2,
		Hash:   "invalidtokenhash",
	}
	err = repo.CheckRefreshToken(ctx, invalidToken)
	assert.Error(t, err)
}
