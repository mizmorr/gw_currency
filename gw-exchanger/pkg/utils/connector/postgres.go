package connector

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	logger "github.com/mizmorr/loggerm"
	"github.com/pkg/errors"
)

func EstablishConnection(ctx context.Context, pool *pgxpool.Pool, attempts int, timeOut time.Duration) error {
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
		return errors.Wrap(err, "Cannot connect to PostgreSQL database")
	}
	return nil
}
