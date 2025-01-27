package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mizmorr/gw_currency/gw-exchanger/internal/app"
)

func main() {
	if err := execute(); err != nil {
		panic(err)
	}
}

func execute() error {
	ctx := context.Background()
	app := app.New()

	if err := app.Start(ctx); err != nil {
		return err
	}

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-stopCh
	return app.Stop(ctx)
}
