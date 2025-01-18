package app

import (
	"context"

	"github.com/mizmorr/gw_currency/gw-exchanger/internal/config"
	"github.com/mizmorr/gw_currency/gw-exchanger/internal/controller"
	"github.com/mizmorr/gw_currency/gw-exchanger/internal/server"
	"github.com/mizmorr/gw_currency/gw-exchanger/internal/service"
	pg "github.com/mizmorr/gw_currency/gw-exchanger/internal/storage/postgres"
	"github.com/mizmorr/gw_currency/gw-exchanger/pkg/utils/lifecycle"
	logger "github.com/mizmorr/loggerm"
)

type companent struct {
	Name    string
	Service lifecycle.Lifecycle
}

type App struct {
	log    *logger.Logger
	config *config.Config
	comps  []companent
	ctx    context.Context
}

func New(c context.Context) *App {
	var (
		config = config.Get()
		log    = logger.Get(config.PathFile, config.Level)
		ctx    = context.WithValue(c, "logger", log)
	)
	return &App{
		log:    log,
		config: config,
		ctx:    ctx,
	}
}

func (a *App) Start() error {
	repo, err := pg.NewPostgresRepo(a.ctx)
	if err != nil {
		return err
	}

	svc := service.NewExchangerService(repo)

	control := controller.NewExchangeController(svc)

	server, err := server.New(a.ctx, control, a.config.Host, string(a.config.Port))
	if err != nil {
		return err
	}
	okCh, errCh := make(chan interface{}), make(chan error)

	a.comps = []companent{
		{Name: "server", Service: server},
		{Name: "service", Service: svc},
	}

	go func() {
		for _, comp := range a.comps {
			err := comp.Service.Start(a.ctx)
			if err != nil {
				errCh <- err
				return
			}
		}
		okCh <- struct{}{}
	}()
	select {
	case err := <-errCh:
		return err
	case <-okCh:
		a.log.Info().Msg("Service started")
		return nil
	}
}

func (a *App) Stop() error {
	a.log.Info().Msg("Graceful shoutdown is running..")
	okCh, errCh := make(chan interface{}), make(chan error)
	go func() {
		for _, comp := range a.comps {
			err := comp.Service.Stop(a.ctx)
			if err != nil {
				errCh <- err
			}
		}
		okCh <- struct{}{}
	}()
	select {
	case err := <-errCh:
		return err
	case <-okCh:
		a.log.Info().Msg("Service stopped")
		return nil
	}
}
