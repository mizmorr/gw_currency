package grpc

import (
	"context"
	"time"

	pb "github.com/mizmorr/grpc_exchange/exchange"
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

func (c *ExchangerClient) GetAllRates() (*pb.ExchangeRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rates, err := c.client.GetAllRates(ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	return rates, nil
}

func (c *ExchangerClient) GetSpecificRate(code string) (*pb.ExchangeRateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rate, err := c.client.GetSpecificRate(ctx, &pb.CurrencyRequest{CurrencyCode: code})
	if err != nil {
		return nil, err
	}

	return rate, nil
}
