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
	LoginUser(ctx context.Context, user *domain.AuthorizationRequst) (*domain.TokenRepsonse, error)
	GetBalance(ctx context.Context, userid int64) ([]*domain.BalanceResponse, error)
	Deposit(ctx context.Context, userid int64, req *domain.DepositRequest) ([]*domain.BalanceResponse, error)
	Withdraw(ctx context.Context, userid int64, req *domain.WithdrawRequest) ([]*domain.BalanceResponse, error)
	ExchangeRates(ctx context.Context) ([]*domain.RateResponse, error)
	Exchange(ctx context.Context, userid int64, req *domain.ExchangeRequest) (*domain.ExchangeResponse, error)
	Refresh(ctx context.Context, req *domain.RefreshRequest) (*domain.TokenRepsonse, error)

	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type WalletController struct {
	service WalletExchangeService
}

func NewWalletController(service WalletExchangeService) *WalletController {
	return &WalletController{
		service: service,
	}
}

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

func (wc *WalletController) Login(c *gin.Context) {
	var req domain.AuthorizationRequst
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

func (wc *WalletController) ExchangeRatesHandler(c *gin.Context) {
	rates, err := wc.service.ExchangeRates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "failed to get exchange rates").Error()})
		return
	}

	c.JSON(http.StatusOK, rates)
}

func (wc *WalletController) ExchangeHandler(c *gin.Context) {
	var req domain.ExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID, exists := c.Get("userID")
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
