package grpc

import (
	"context"
	"time"

	pb "github.com/mizmorr/grpc_exchange/exchange"
	logger "github.com/mizmorr/loggerm"
	"google.golang.org/grpc"
)

type ExchangerClient struct {
	client pb.CurrencyExchangeServiceClient
	conn   *grpc.ClientConn
}

func NewExchangerClient(conn *grpc.ClientConn) *ExchangerClient {
	return &ExchangerClient{
		conn: conn,
	}
}

func (c *ExchangerClient) Start(ctx context.Context) error {
	log := logger.GetLoggerFromContext(ctx)

	c.client = pb.NewCurrencyExchangeServiceClient(c.conn)

	log.Info().Msg("GRPC client is started")

	return nil
}

func (c *ExchangerClient) Stop(ctx context.Context) error {
	log := logger.GetLoggerFromContext(ctx)

	log.Info().Msg("GRPC client is stopping..")

	return c.conn.Close()
}

func (c *ExchangerClient) GetAllRates(cont context.Context) (*pb.ExchangeRatesResponse, error) {
	ctx, cancel := context.WithTimeout(cont, 5*time.Second)
	defer cancel()

	rates, err := c.client.GetAllRates(ctx, &pb.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	return rates, nil
}

func (c *ExchangerClient) GetSpecificRate(cont context.Context, code string) (*pb.ExchangeRateResponse, error) {
	ctx, cancel := context.WithTimeout(cont, 5*time.Second)
	defer cancel()

	rate, err := c.client.GetSpecificRate(ctx, &pb.CurrencyRequest{CurrencyCode: code})
	if err != nil {
		return nil, err
	}

	return rate, nil
}
