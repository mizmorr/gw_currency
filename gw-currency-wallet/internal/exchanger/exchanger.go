package exchanger

import (
	"context"
	"time"

	pb "github.com/mizmorr/grpc_exchange/exchange"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/domain"
)

type RemoteExchanger interface {
	GetAllRates() (*pb.ExchangeRatesResponse, error)
	GetSpecificRate(code string) (*pb.ExchangeRateResponse, error)
}

type Cash interface {
	Set(ctx context.Context, key string, value float64, expiration time.Duration) error
	Get(ctx context.Context, key string) (float64, error)
}

type Exchanger struct {
	remote   RemoteExchanger
	cash     Cash
	cacheTTL time.Duration
	rates    []string
}

func New(remote RemoteExchanger, cash Cash, rates []string, ttl time.Duration) (*Exchanger, error) {
	return &Exchanger{
		remote:   remote,
		cash:     cash,
		cacheTTL: ttl,
		rates:    rates,
	}, nil
}

func (c *Exchanger) GetExchangeRate(ctx context.Context, currencyCode string) (*domain.RateResponse, error) {
	if rate, err := c.cash.Get(ctx, currencyCode); err == nil {
		return &domain.RateResponse{CurrencyCode: currencyCode, Value: rate}, nil
	}

	rate, err := c.remote.GetSpecificRate(currencyCode)
	if err != nil {
		return nil, err
	}
	_ = c.cash.Set(ctx, currencyCode, rate.Rate, c.cacheTTL)

	return &domain.RateResponse{
		CurrencyCode: rate.CurrencyCode,
		Value:        rate.Rate,
	}, nil
}

func (e *Exchanger) GetExchangeRates(ctx context.Context) ([]*domain.RateResponse, error) {
	var (
		notFound []string
		result   = make([]*domain.RateResponse, 0, len(e.rates))
	)
	notFound, result = e.scanCash(ctx)

	if len(notFound) == 0 {
		return result, nil
	}

	if len(notFound) < len(e.rates) {
		return e.mediumPath(ctx, notFound, result)
	}
	return e.slowPath(ctx, result)
}

func (e *Exchanger) scanCash(ctx context.Context) (notFound []string, result []*domain.RateResponse) {
	for _, currencyCode := range e.rates {
		rate, err := e.cash.Get(ctx, currencyCode)
		if err != nil {
			notFound = append(notFound, currencyCode)
		} else {
			result = append(result, &domain.RateResponse{
				CurrencyCode: currencyCode,
				Value:        rate,
			})
		}
	}
	return
}

func (e *Exchanger) mediumPath(ctx context.Context, codes []string, result []*domain.RateResponse) ([]*domain.RateResponse, error) {
	for _, currencyCode := range codes {
		rate, err := e.remote.GetSpecificRate(currencyCode)
		if err != nil {
			return nil, err
		}
		_ = e.cash.Set(ctx, currencyCode, rate.Rate, e.cacheTTL)

		result = append(result, &domain.RateResponse{
			CurrencyCode: rate.CurrencyCode,
			Value:        rate.Rate,
		})
	}
	return result, nil
}

func (e *Exchanger) slowPath(ctx context.Context, result []*domain.RateResponse) ([]*domain.RateResponse, error) {
	rates, err := e.remote.GetAllRates()
	if err != nil {
		return nil, err
	}

	for _, r := range rates.Rates {
		_ = e.cash.Set(ctx, r.CurrencyCode, r.Rate, e.cacheTTL)

		result = append(result, &domain.RateResponse{
			CurrencyCode: r.CurrencyCode,
			Value:        r.Rate,
		})
	}
	return result, nil
}
