package service

import (
	"context"
	"time"

	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/domain"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/mappers"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/store"
	jwttoken "github.com/mizmorr/gw_currency/gw-currency-wallet/pkg/jwtToken"
	"github.com/pkg/errors"
)

func (ws *WalletService) RegisterUser(ctx context.Context, user *domain.RegisterRequest) error {
	storeUser, err := mappers.ToStoreUserFromRegister(user)
	if err != nil {
		return err
	}

	return ws.repo.CreateUser(ctx, storeUser)
}

func (ws *WalletService) LoginUser(ctx context.Context, user *domain.AuthorizationRequst) (*domain.TokenRepsonse, error) {
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
		Hash:      tokenHash,
		UserID:    userid,
		ExpiresAt: time.Now().Add(ws.optsJWT.RefreshExpiresTime),
		Revoked:   false,
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
	isEnough, err := ws.isBalanceEnough(ctx, userid, req.Currency, req.Amount)
	if err != nil {
		return nil, err
	}
	if !isEnough {
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

func (ws *WalletService) isBalanceEnough(ctx context.Context, userid int64, currencyCode string, amount float64) (bool, error) {
	getCurrencyReq := &store.CurrencyRequest{
		CurrencyCode: currencyCode,
		UserID:       userid,
	}
	currency, err := ws.repo.GetSpecificCurrency(ctx, getCurrencyReq)
	if err != nil {
		return false, err
	}
	if currency.Balance < amount {
		return false, nil
	}
	return true, nil
}

func (ws *WalletService) ExchangeRates(ctx context.Context) ([]*domain.RateResponse, error) {
	response, err := ws.exchanger.GetExchangeRates(ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (ws *WalletService) Exchange(ctx context.Context, userid int64, req *domain.ExchangeRequest) (*domain.ExchangeResponse, error) {
	isEnough, err := ws.isBalanceEnough(ctx, userid, req.BaseCurrency, req.Amount)
	if err != nil {
		return nil, err
	}
	if !isEnough {
		return nil, errors.New("Insufficient funds")
	}

	baseCurrencyRate, err := ws.exchanger.GetExchangeRate(ctx, req.BaseCurrency)
	if err != nil {
		return nil, err
	}

	targetCurrencyRate, err := ws.exchanger.GetExchangeRate(ctx, req.TargetCurrency)
	if err != nil {
		return nil, err
	}

	valueTpDeposit := ws.convertCurrency(req.Amount, baseCurrencyRate.Value, targetCurrencyRate.Value)

	err = ws.makeTransfer(ctx, req.BaseCurrency, req.TargetCurrency, req.Amount, valueTpDeposit, userid)
	if err != nil {
		return nil, err
	}
	newBalance, err := ws.repo.GetBalance(ctx, userid)
	if err != nil {
		return nil, err
	}

	return &domain.ExchangeResponse{
		Message:        "Exchange successful",
		ExchangeAmount: valueTpDeposit,
		NewBalance:     mappers.ToDomainBalance(newBalance),
	}, nil
}

func (ws *WalletService) convertCurrency(amount float64, rateFrom float64, rateTo float64) float64 {
	return amount * (rateTo / rateFrom)
}

func (ws *WalletService) makeTransfer(ctx context.Context,
	fromCurrencyCode, toCurrencyCode string, amountFrom, amountTo float64, userid int64,
) error {
	storeExchangeReq := &store.ExchangeBalance{
		UserID:       userid,
		FromCurrency: fromCurrencyCode,
		ToCurrency:   toCurrencyCode,
		ToAmount:     amountTo,
		FromAmount:   amountFrom,
	}

	return ws.repo.ExchangeCurrency(ctx, storeExchangeReq)
}

func (ws *WalletService) Refresh(ctx context.Context, req *domain.RefreshRequest) (*domain.TokenRepsonse, error) {
	err := jwttoken.Validate(req.TokenHash, []byte(ws.optsJWT.RefreshSecret))
	if err != nil {
		return nil, errors.Wrap(err, "The time of existence has expired. You need to log in again")
	}

	userID, err := jwttoken.GetUserID(req.TokenHash, []byte(ws.optsJWT.RefreshSecret))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user id")
	}

	tokenToCheck := &store.RefreshToken{
		Hash:   req.TokenHash,
		UserID: userID,
	}

	err = ws.repo.CheckRefreshToken(ctx, tokenToCheck)
	if err != nil {
		return nil, errors.Wrap(err, "invalid refresh token")
	}

	access, _, err := ws.generateTokens(userID)
	if err != nil {
		return nil, err
	}

	return &domain.TokenRepsonse{
		Access:  access,
		Refresh: req.TokenHash,
	}, nil
}
