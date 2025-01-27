package service

import (
	"context"

	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/config"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/domain"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/store"
)

type Repository interface {
	CreateUser(ctx context.Context, user *store.User) error
	Authentication(ctx context.Context, user *store.User) (int64, error)
	SetToken(ctx context.Context, refresh *store.RefreshToken) error
	GetBalance(ctx context.Context, userid int64) ([]*store.WalletCurrency, error)
	UpdateBalance(ctx context.Context, newBalance *store.UpdateBalance) error
	ExchangeCurrency(ctx context.Context, exchangeBody *store.ExchangeBalance) error
	GetSpecificCurrency(ctx context.Context, req *store.CurrencyRequest) (*store.WalletCurrency, error)
	CheckRefreshToken(ctx context.Context, token *store.RefreshToken) error
}

type RateExchanger interface {
	GetExchangeRate(ctx context.Context, currencyCode string) (*domain.RateResponse, error)
	GetExchangeRates(ctx context.Context) ([]*domain.RateResponse, error)
}

type WalletService struct {
	repo      Repository
	optsJWT   config.JWTtokens
	exchanger RateExchanger
}

func New(repo Repository, exch RateExchanger, tokensOpt config.JWTtokens) *WalletService {
	return &WalletService{
		repo:      repo,
		exchanger: exch,
		optsJWT:   tokensOpt,
	}
}
