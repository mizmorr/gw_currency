package service

import (
	"context"

	pb "github.com/mizmorr/grpc_exchange/exchange"
	"github.com/mizmorr/gw_currency/gw-exchanger/internal/storage"
	"github.com/mizmorr/gw_currency/gw-exchanger/internal/storage/model"
	"github.com/pkg/errors"
)

type ExchangerService struct {
	store storage.Repository
}

func NewExchangerService(store storage.Repository) *ExchangerService {
	return &ExchangerService{
		store: store,
	}
}

func (svc *ExchangerService) GetAllRates(ctx context.Context) (*pb.ExchangeRatesResponse, error) {
	rates, err := svc.store.GetAllRates(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all rates")
	}

	return model.RatesToResponse(
		rates,
	), nil
}

func (svc *ExchangerService) GetRate(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	rate, err := svc.store.GetRate(ctx, req.CurrencyCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get rate")
	}

	return rate.ToGRPC(), nil
}

func (svc *ExchangerService) Start(ctx context.Context) error {
	return svc.store.Start(ctx)
}

func (svc *ExchangerService) Stop(ctx context.Context) error {
	svc.store.Stop(ctx)
	return nil
}
