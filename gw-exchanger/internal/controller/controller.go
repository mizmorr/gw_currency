package controller

import (
	"context"

	pb "github.com/mizmorr/grpc_exchange/exchange"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	GetAllRates(ctx context.Context) (*pb.ExchangeRatesResponse, error)
	GetRate(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error)
}

type ExchangeController struct {
	service Service
	pb.UnimplementedCurrencyExchangeServiceServer
}

func NewExchangeController(svc Service) *ExchangeController {
	return &ExchangeController{
		service: svc,
	}
}

func (c *ExchangeController) GetAllRates(ctx context.Context, req *pb.EmptyRequest) (*pb.ExchangeRatesResponse, error) {
	rateList, err := c.service.GetAllRates(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to get all rates")
	}
	return rateList, nil
}

func (c *ExchangeController) GetSpecificRate(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	rate, err := c.service.GetRate(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to get rate")
	}
	return rate, nil
}
