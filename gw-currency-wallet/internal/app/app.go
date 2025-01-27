package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/config"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/delivery"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/exchanger"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/grpc"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/middleware"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/service"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/store/postgres"
	httpserver "github.com/mizmorr/gw_currency/gw-currency-wallet/pkg/httpServer"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/pkg/redis"
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
	cmps   []component
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

	if err := a.setUp(ctx); err != nil {
		return err
	}

	okCh, errCh := make(chan interface{}), make(chan error)

	go func() {
		for _, comp := range a.cmps {
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
		for _, comp := range a.cmps {
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
		a.log.Info().Msg("Wallet service is stopped")
		return nil
	}
}

func (a *App) setUp(ctx context.Context) error {
	repo, err := postgres.NewPostgresRepo(ctx)
	if err != nil {
		return err
	}

	conn := grpc.NewConnection(a.config.GRPC.Host, a.config.GRPC.Port)

	remoteExchanger := grpc.NewExchangerClient(conn)

	cashExchanger := redis.NewRedisClient(ctx, a.config.Redis.Host, a.config.Redis.Port, a.config.Redis.Password)

	exchanger, err := exchanger.New(remoteExchanger, cashExchanger, a.config.CurrencyCodes, a.config.Redis.TTL)
	if err != nil {
		return err
	}

	service := service.New(repo, exchanger, a.config.JWTtokens)

	walletController := delivery.NewWalletController(service)

	handler := gin.New()

	authMiddleware := middleware.JWTAuthMiddleware(a.config.JWTtokens.AccessSecret)

	delivery.NewRouter(handler, authMiddleware, walletController)

	httpServer := httpserver.New(handler, a.config.HttpHost, a.config.HttpPort, a.config.ShutdownTimeout)

	a.cmps = append(a.cmps, component{Name: "server", Service: httpServer},
		component{Name: "repo", Service: repo},
		component{Name: "exchangerRemote", Service: remoteExchanger},
		component{Name: "exchangerCash", Service: cashExchanger},
	)

	return nil
}
