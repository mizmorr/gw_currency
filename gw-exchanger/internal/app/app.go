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

type component struct {
	Name    string
	Service lifecycle.Lifecycle
}

type App struct {
	log    *logger.Logger
	config *config.Config
	comps  []component
}

func New() *App {
	var (
		config = config.Get()
		log    = logger.Get(config.PathFile, config.Level)
	)
	return &App{
		log:    log,
		config: config,
	}
}

func (a *App) Start(ctx context.Context) error {
	if _, ok := ctx.Value("logger").(*logger.Logger); !ok {
		ctx = context.WithValue(ctx, "logger", a.log)
	}

	repo, err := pg.NewPostgresRepo(ctx)
	if err != nil {
		return err
	}

	svc := service.NewExchangerService(repo)

	control := controller.NewExchangeController(svc)

	server, err := server.New(ctx, control, a.config.Host, string(a.config.Port))
	if err != nil {
		return err
	}
	okCh, errCh := make(chan interface{}), make(chan error)

	a.comps = []component{
		{Name: "server", Service: server},
		{Name: "service", Service: svc},
	}

	go func() {
		for _, comp := range a.comps {
			err := comp.Service.Start(ctx)
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

func (a *App) Stop(ctx context.Context) error {
	if _, ok := ctx.Value("logger").(*logger.Logger); !ok {
		ctx = context.WithValue(ctx, "logger", a.log)
	}

	a.log.Info().Msg("Graceful shutdown is running..")

	okCh, errCh := make(chan interface{}), make(chan error)

	go func() {
		for _, comp := range a.comps {
			err := comp.Service.Stop(ctx)
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
		a.log.Info().Msg("Exchange service is stopped")
		return nil
	}
}
