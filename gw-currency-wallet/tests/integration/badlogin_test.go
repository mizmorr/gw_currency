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

func TestInvalidLogin(t *testing.T) {
	config := config.Get()
	address := net.JoinHostPort(config.HttpHost, config.HttpPort)
	serverURL := fmt.Sprintf("http://%s/api/v1", address)

	// 1. Регистрация пользователя
	user := domain.RegisterRequest{Username: "userforlogin", Email: "userforlogin@example.com", Password: "password"}
	userBody, _ := json.Marshal(user)
	resp, err := http.Post(serverURL+"/register", "application/json", bytes.NewBuffer(userBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 2. Неверный логин
	loginReq := domain.AuthorizationRequest{Username: "userforlogin", Password: "wrongpassword"}
	loginBody, _ := json.Marshal(loginReq)
	resp, err = http.Post(serverURL+"/login", "application/json", bytes.NewBuffer(loginBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
