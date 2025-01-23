package service

import (
	"context"
	"errors"
	"time"

	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/config"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/domain"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/mappers"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/store"
	jwttoken "github.com/mizmorr/gw_currency/gw-currency-wallet/pkg/jwtToken"
)

type Repository interface {
	CreateUser(ctx context.Context, user *store.User) error
	Authentication(ctx context.Context, user *store.User) (int64, error)
	SetToken(ctx context.Context, refresh *store.RefreshToken) error
	GetBalance(ctx context.Context, userid int64) ([]*store.WalletCurrency, error)
	UpdateBalance(ctx context.Context, newBalance *store.UpdateBalance) error
	ExchangeCurrency(ctx context.Context, exchangeBody *store.ExchangeBalance) error

	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Cash interface{}

type WalletService struct {
	repo    Repository
	cash    Cash
	optsJWT config.JWTtokens
}

func New(repo Repository, cash Cash, tokensOpt config.JWTtokens) *WalletService {
	return &WalletService{
		repo:    repo,
		cash:    cash,
		optsJWT: tokensOpt,
	}
}

func (ws *WalletService) CreateUser(ctx context.Context, user *domain.RegisterRequest) error {
	storeUser, err := mappers.ToStoreUserFromRegister(user)
	if err != nil {
		return err
	}

	return ws.repo.CreateUser(ctx, storeUser)
}

func (ws *WalletService) Authorization(ctx context.Context, user *domain.AuthorizationRequst) (*domain.TokenRepsonse, error) {
	storeUser, err := mappers.ToStoreUserFromAuthorize(user)
	if err != nil {
		return nil, err
	}
	var id int64

	id, err = ws.repo.Authentication(ctx, storeUser)
	if err != nil {
		return nil, err
	}
	access, refresh, err := ws.generateTokens(id)
	if err != nil {
		return nil, err
	}
	err = ws.storeRefreshToken(ctx, id, refresh)
	if err != nil {
		return nil, err
	}
	return &domain.TokenRepsonse{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func (ws *WalletService) generateTokens(userid int64) (string, string, error) {
	tokenOpts := &jwttoken.TokensOption{
		UserID:        userid,
		AccessExp:     ws.optsJWT.AccessExpiresTime,
		RefreshExp:    ws.optsJWT.RefreshExpiresTime,
		SecretAccess:  ws.optsJWT.AccessSecret,
		SecretRefresh: ws.optsJWT.RefreshSecret,
	}

	access, refresh, err := jwttoken.GenerateTokens(tokenOpts)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (ws *WalletService) storeRefreshToken(ctx context.Context, userid int64, tokenHash string) error {
	refreshToken := &store.RefreshToken{
		RefreshHash: tokenHash,
		UserID:      userid,
		ExpiresAt:   time.Now().Add(ws.optsJWT.RefreshExpiresTime),
		Revoked:     false,
	}

	return ws.repo.SetToken(ctx, refreshToken)
}

func (ws *WalletService) GetBalance(ctx context.Context, userid int64) ([]*domain.BalanceResponse, error) {
	balance, err := ws.repo.GetBalance(ctx, userid)
	if err != nil {
		return nil, err
	}
	return mappers.ToDomainBalance(balance), nil
}

func (ws *WalletService) Deposit(ctx context.Context, userid int64, req *domain.DepositRequest) ([]*domain.BalanceResponse, error) {
	depositInStore := mappers.ToStoreDepositBalance(userid, req)

	err := ws.repo.UpdateBalance(ctx, depositInStore)
	if err != nil {
		return nil, err
	}

	newBalance, err := ws.repo.GetBalance(ctx, userid)
	if err != nil {
		return nil, err
	}
	return mappers.ToDomainBalance(newBalance), nil
}

func (ws *WalletService) Withdraw(ctx context.Context, userid int64, req *domain.WithdrawRequest) ([]*domain.BalanceResponse, error) {
	balance, err := ws.repo.GetBalance(ctx, userid)
	if err != nil {
		return nil, err
	}

	isBalanceCorrect := ws.checkCurrency(balance, req.Currency, req.Amount)
	if !isBalanceCorrect {
		return nil, errors.New("Insufficient funds")
	}

	withdrawInStore := mappers.ToStoreWithdrawBalance(userid, req)

	err = ws.repo.UpdateBalance(ctx, withdrawInStore)
	if err != nil {
		return nil, err
	}

	newBalance, err := ws.repo.GetBalance(ctx, userid)
	if err != nil {
		return nil, err
	}
	return mappers.ToDomainBalance(newBalance), nil
}

func (ws *WalletService) checkCurrency(balance []*store.WalletCurrency, currency string, amount float64) bool {
	for _, b := range balance {
		if b.Currency == currency {
			if b.Balance < amount {
				return false
			}
		}
	}
	return true
}

func (ws *WalletService) ExchangeRates(ctx context.Context)
