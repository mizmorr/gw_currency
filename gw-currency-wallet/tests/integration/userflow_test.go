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

func TestUserFlow(t *testing.T) {
	config := config.Get()

	address := net.JoinHostPort(config.HttpHost, config.HttpPort)

	serverURL := fmt.Sprintf("http://%s/api/v1", address)

	// 1. Регистрация пользователя
	user := domain.RegisterRequest{Username: "realnewuser", Email: "realnewuser@example.com", Password: "password"}
	userBody, _ := json.Marshal(user)
	resp, err := http.Post(serverURL+"/register", "application/json", bytes.NewBuffer(userBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 2. Логин
	loginReq := domain.AuthorizationRequest{Username: "realnewuser", Password: "password"}
	loginBody, _ := json.Marshal(loginReq)
	resp, err = http.Post(serverURL+"/login", "application/json", bytes.NewBuffer(loginBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var tokenResp domain.TokenResponse
	json.NewDecoder(resp.Body).Decode(&tokenResp)
	assert.NotEmpty(t, tokenResp.Access)

	// 3. Депозит средств
	depositReq := domain.DepositRequest{Amount: 1000.0, Currency: "USD"}
	depositBody, _ := json.Marshal(depositReq)
	req, _ := http.NewRequest("POST", serverURL+"/wallet/deposit", bytes.NewBuffer(depositBody))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Access)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 4. Обмен валюты
	exchangeReq := domain.ExchangeRequest{BaseCurrency: "USD", TargetCurrency: "EUR", Amount: 500.0}
	exchangeBody, _ := json.Marshal(exchangeReq)
	req, _ = http.NewRequest("POST", serverURL+"/exchange", bytes.NewBuffer(exchangeBody))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Access)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
