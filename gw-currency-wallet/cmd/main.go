package main

import (
	"context"

	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/config"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/store"
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

	user := &store.User{
		Username: "bob",
		Email:    "example@example.com",
		Password: "password",
	}

	err = repo.CreateUser(ctx, user)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create user")
	}

	for {
	}
}
