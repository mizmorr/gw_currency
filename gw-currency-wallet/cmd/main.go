package main

import (
	"context"
	"fmt"

	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/config"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/domain"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/service"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/store/postgres"
	logger "github.com/mizmorr/loggerm"
)

func main() {
	config := config.Get()
	log := logger.Get(config.PathFile, config.Level)
	ctx := context.WithValue(context.Background(), "logger", log)

	repo, err := postgres.NewPostgresRepo(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Postgres repository")
	}

	err = repo.Start(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start Postgres repository")
	}

	var k interface{}
	svc := service.New(repo, k, config.JWTtokens)

	balance, err := svc.GetBalance(ctx, 1)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get balance")
	}
	for _, b := range balance {
		log.Info().Msg(fmt.Sprintf("Currency: %s, Balance: %f", b.Currency, b.Value))
	}
	depReq := &domain.WithdrawRequest{
		Currency: "EUR",
		Amount:   100.0,
	}
	resp, err := svc.Withdraw(ctx, 1, depReq)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to deposit")
	}
	for _, b := range resp {
		log.Info().Msg(fmt.Sprintf("Currency: %s, Balance: %f", b.Currency, b.Value))
	}

	for {
	}
}
