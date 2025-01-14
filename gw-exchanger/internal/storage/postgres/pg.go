package pg

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mizmorr/gw_currency/gw-exchanger/internal/config"
	logger "github.com/mizmorr/loggerm"
)

type pg struct {
	*pgxpool.Pool
}

var (
	once       sync.Once
	pgInstance *pg
)

func dial(ctx context.Context) (*pg, error) {
	confg := config.Get()
	log := logger.GetLoggerFromContext(ctx)

	log.Debug().Msg("Database url checking...")
	if confg.PostgresURL == "" {
		return nil, errors.New("No database URL provided")
	}

	log.Debug().Msg("Database config parsing...")
	poolConfig, err := pgxpool.ParseConfig(confg.PostgresURL)
	if err != nil {
		return nil, errors.Wrap(err, "Parse config failed")
	}

	poolConfig.MaxConnIdleTime = confg.PostgresMaxIdleTime
	poolConfig.HealthCheckPeriod = confg.PostgresHealthCheckPeriod
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create pgx pool")
	}

	once.Do(func() {
		establishConnection(ctx, pool, confg.PostgresConnectAttempts, confg.PostgresTimeout)
	})

	return pgInstance, nil
}

func establishConnection(ctx context.Context, pool *pgxpool.Pool, attempts int, timeOut time.Duration) {
	log := logger.GetLoggerFromContext(ctx)
	var (
		conn *pgxpool.Conn
		err  error
	)

	for attempts > 0 {

		conn, err = pool.Acquire(ctx)
		if err == nil {
			log.Info().Msg("Connect to postgres is established")
			conn.Release()

			break
		}

		log.Error().Err(err).Msg("Failed to connect to pg, retrying...")

		time.Sleep(timeOut)

		attempts--
	}
	if err != nil {
		panic(errors.Wrap(err, "Cannot connect to PostgreSQL database"))
	}
	pgInstance = &pg{pool}
}
