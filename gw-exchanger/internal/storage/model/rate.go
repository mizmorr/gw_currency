package model

import pb "github.com/mizmorr/grpc_exchange/exchange"

type Rate struct {
	CurrencyCode string
	Value        float64
}

func (r Rate) ToGRPC() *pb.ExchangeRateResponse {
	return &pb.ExchangeRateResponse{
		CurrencyCode: r.CurrencyCode,
		Rate:         r.Value,
	}
}

func RatesToResponse(rates []*Rate) *pb.ExchangeRatesResponse {
	var exchangeRates []*pb.ExchangeRate

	for _, r := range rates {
		exchangeRates = append(exchangeRates, &pb.ExchangeRate{
			CurrencyCode: r.CurrencyCode,
			Rate:         r.Value,
		})
	}

	return &pb.ExchangeRatesResponse{
		Rates: exchangeRates,
	}
}
