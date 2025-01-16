package main

import (
	"context"
	"time"

	pg "github.com/mizmorr/gw_currency/gw-exchanger/internal/storage/postgres"
	logger "github.com/mizmorr/loggerm"
)

func main() {
	log := logger.Get("debug")
	ctx := context.WithValue(context.Background(), "logger", log)

	repo, err := pg.NewPostgresRepo(ctx)
	if err != nil {
		panic(err)
	}
	err = repo.Start(ctx)
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(time.Second * 15)
		repo.Stop(ctx)
	}()
	for {
	}
}
