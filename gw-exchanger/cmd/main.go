package main

import (
	"context"

	pg "github.com/mizmorr/gw_currency/gw-exchanger/internal/storage/postgres"
	logger "github.com/mizmorr/loggerm"
)

func main() {
	log := logger.Get("debug")
	ctx := context.WithValue(context.Background(), "logger", log)

	_, err := pg.NewPostgresRepo(ctx)
	if err != nil {
		panic(err)
	}
	for {
	}
}
