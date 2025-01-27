package delivery

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/domain"
	"github.com/pkg/errors"
)

type WalletExchangeService interface {
	RegisterUser(ctx context.Context, user *domain.RegisterRequest) error
	LoginUser(ctx context.Context, user *domain.AuthorizationRequest) (*domain.TokenResponse, error)
	GetBalance(ctx context.Context, userid int64) ([]*domain.BalanceResponse, error)
	Deposit(ctx context.Context, userid int64, req *domain.DepositRequest) ([]*domain.BalanceResponse, error)
	Withdraw(ctx context.Context, userid int64, req *domain.WithdrawRequest) ([]*domain.BalanceResponse, error)
	ExchangeRates(ctx context.Context) ([]*domain.RateResponse, error)
	Exchange(ctx context.Context, userid int64, req *domain.ExchangeRequest) (*domain.ExchangeResponse, error)
	Refresh(ctx context.Context, req *domain.RefreshRequest) (*domain.TokenResponse, error)
}

type WalletController struct {
	service WalletExchangeService
}

func NewWalletController(service WalletExchangeService) *WalletController {
	return &WalletController{
		service: service,
	}
}

// @Summary Register a new user
// @Description Registers a new user with the provided credentials
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body domain.RegisterRequest true "User registration data"
// @Success 201 {object} map[string]string "user registered successfully"
// @Failure 400 {object} map[string]string "invalid request"
// @Router /register [post]
func (wc *WalletController) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := wc.service.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.Wrap(err, "failed to register user").Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// @Summary User login
// @Description Authenticates a user and returns JWT tokens
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body domain.AuthorizationRequest true "User credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "invalid credentials"
// @Router /login [post]
func (wc *WalletController) Login(c *gin.Context) {
	var req domain.AuthorizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	tokens, err := wc.service.LoginUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// @Summary Get user balance
// @Description Returns the balance of an authenticated user
// @Tags wallet
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "failed to get balance"
// @Router /balance [get]
func (wc *WalletController) GetBalance(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	balance, err := wc.service.GetBalance(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
		return
	}

	c.JSON(http.StatusOK, balance)
}

// @Summary Deposit funds
// @Description Deposits funds into the user's account
// @Tags wallet
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body domain.DepositRequest true "Deposit data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "failed to deposit"
// @Router /deposit [post]
func (wc *WalletController) Deposit(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	newBalance, err := wc.service.Deposit(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deposit"})
		return
	}

	c.JSON(http.StatusOK, newBalance)
}

// @Summary Withdraw funds
// @Description Withdraws funds from the user's account
// @Tags wallet
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body domain.WithdrawRequest true "Withdraw data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "failed to withdraw"
// @Router /withdraw [post]
func (wc *WalletController) Withdraw(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	newBalance, err := wc.service.Withdraw(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to withdraw"})
		return
	}

	c.JSON(http.StatusOK, newBalance)
}

// @Summary Refresh access token
// @Description Refreshes the user's access token using a refresh token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body domain.RefreshRequest true "Refresh token data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 401 {object} map[string]string "invalid refresh token"
// @Router /refresh [post]
func (wc *WalletController) Refresh(c *gin.Context) {
	var req domain.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	newTokens, err := wc.service.Refresh(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, newTokens)
}

// @Summary Get exchange rates
// @Description Fetches the latest exchange rates
// @Tags exchange
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string "failed to get exchange rates"
// @Router /exchange-rates [get]
func (wc *WalletController) ExchangeRatesHandler(c *gin.Context) {
	rates, err := wc.service.ExchangeRates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "failed to get exchange rates").Error()})
		return
	}

	c.JSON(http.StatusOK, rates)
}

// @Summary Exchange currency
// @Description Exchanges one currency for another
// @Tags exchange
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body domain.ExchangeRequest true "Exchange data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 400 {object} map[string]string "exchange failed"
// @Router /exchange [post]
func (wc *WalletController) ExchangeHandler(c *gin.Context) {
	var req domain.ExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	exchangeResponse, err := wc.service.Exchange(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.Wrap(err, "exchange failed").Error()})
		return
	}

	c.JSON(http.StatusOK, exchangeResponse)
}
