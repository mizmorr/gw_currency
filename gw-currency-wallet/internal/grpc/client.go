package grpc

import (
	"context"
	"time"

	pb "github.com/mizmorr/grpc_exchange/exchange"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/domain"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/mappers"
	"google.golang.org/grpc"
)

type ExchangerClient struct {
	client pb.CurrencyExchangeServiceClient
}

func NewExchangerClient(conn *grpc.ClientConn) *ExchangerClient {
	return &ExchangerClient{
		client: pb.NewCurrencyExchangeServiceClient(conn),
	}
}

func (c *ExchangerClient) GetAllRates() ([]*domain.RateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rate, err := c.client.GetAllRates(ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	return mappers.ToDomainRatesResponse(rate), nil
}

func (c *ExchangerClient) GetSpecificRate(code string) (*domain.RateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rate, err := c.client.GetSpecificRate(ctx, &pb.CurrencyRequest{CurrencyCode: code})
	if err != nil {
		return nil, err
	}

	return mappers.ToDomainRateResponse(rate), nil
}
