package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/config"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestBalanceAndWithdrawal(t *testing.T) {
	config := config.Get()
	address := net.JoinHostPort(config.HttpHost, config.HttpPort)
	serverURL := fmt.Sprintf("http://%s/api/v1", address)

	// 1. Регистрация пользователя
	user := domain.RegisterRequest{Username: "userforbalance", Email: "userforbalance@example.com", Password: "password"}
	userBody, _ := json.Marshal(user)
	resp, err := http.Post(serverURL+"/register", "application/json", bytes.NewBuffer(userBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 2. Логин
	loginReq := domain.AuthorizationRequest{Username: "userforbalance", Password: "password"}
	loginBody, _ := json.Marshal(loginReq)
	resp, err = http.Post(serverURL+"/login", "application/json", bytes.NewBuffer(loginBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var tokenResp domain.TokenResponse
	json.NewDecoder(resp.Body).Decode(&tokenResp)
	assert.NotEmpty(t, tokenResp.Access)

	// 3. Проверка баланса
	req, _ := http.NewRequest("GET", serverURL+"/wallet/balance", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.Access)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var balanceResp domain.BalanceResponse
	json.NewDecoder(resp.Body).Decode(&balanceResp)
	assert.Greater(t, balanceResp.Value, 0.0)

	// 4. Вывод средств
	withdrawalReq := domain.WithdrawRequest{Amount: 200.0, Currency: "USD"}
	withdrawalBody, _ := json.Marshal(withdrawalReq)
	req, _ = http.NewRequest("POST", serverURL+"/wallet/withdraw", bytes.NewBuffer(withdrawalBody))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Access)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
