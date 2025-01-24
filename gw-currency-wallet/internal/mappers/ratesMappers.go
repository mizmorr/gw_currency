package mappers

import (
	"errors"

	grpc_exchange "github.com/mizmorr/grpc_exchange/exchange"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/domain"
)

func ToDomainRatesResponse(resp *grpc_exchange.ExchangeRatesResponse) []*domain.RateResponse {
	rates := make([]*domain.RateResponse, 0, len(resp.Rates))

	for _, r := range resp.Rates {
		rate := exchangeRateToDomain(r)
		rates = append(rates, rate)
	}
	return rates
}

func ToDomainRateResponse(resp *grpc_exchange.ExchangeRateResponse) *domain.RateResponse {
	return &domain.RateResponse{
		CurrencyCode: resp.CurrencyCode,
		Value:        resp.Rate,
	}
}

func exchangeRateToDomain(rate *grpc_exchange.ExchangeRate) *domain.RateResponse {
	return &domain.RateResponse{
		CurrencyCode: rate.CurrencyCode,
		Value:        rate.Rate,
	}
}

func GetRateOfCurrency(resp *grpc_exchange.ExchangeRatesResponse, currencyCode string) (float64, error) {
	for _, rate := range resp.Rates {
		if rate.CurrencyCode == currencyCode {
			return rate.Rate, nil
		}
	}
	return 0, errors.New("No such currency found")
}
