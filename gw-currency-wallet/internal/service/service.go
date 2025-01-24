package service

import (
	"context"

	grpc_exchange "github.com/mizmorr/grpc_exchange/exchange"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/config"
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

	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Cash interface{}

type ClientGRPC interface {
	GetAllRates() (*grpc_exchange.ExchangeRatesResponse, error)
	GetSpecificRate(code string) (*grpc_exchange.ExchangeRateResponse, error)
}

type WalletService struct {
	repo    Repository
	cash    Cash
	optsJWT config.JWTtokens
	grpc    ClientGRPC
}

func New(repo Repository, cash Cash, tokensOpt config.JWTtokens, client ClientGRPC) *WalletService {
	return &WalletService{
		repo:    repo,
		cash:    cash,
		optsJWT: tokensOpt,
		grpc:    client,
	}
}
