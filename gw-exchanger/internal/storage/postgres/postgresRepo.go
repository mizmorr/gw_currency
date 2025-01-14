package pg

import (
	"context"
	"time"

	logger "github.com/mizmorr/loggerm"
)

type PostgresRepo struct {
	db *pg
}

func NewPostgresRepo(ctx context.Context) (*PostgresRepo, error) {
	db, err := dial(ctx)
	if err != nil {
		return nil, err
	}

	repo := PostgresRepo{
		db: db,
	}

	go repo.keepAlive(ctx)

	return &repo, nil
}

const keepALiveTimeout = 5

func (repo *PostgresRepo) keepAlive(ctx context.Context) {
	log := logger.GetLoggerFromContext(ctx)

	log.Debug().Msg("Keeping database connection alive...")

	for {
		time.Sleep(time.Second * keepALiveTimeout)

		connectionLost := false

		conn, err := repo.db.Acquire(ctx)
		if err != nil {
			connectionLost = true
			log.Debug().Msg("[keepAlive] Lost connection, is trying to reconnect...")
		} else {
			conn.Release()
		}

		if connectionLost {
			repo.db, err = dial(ctx)
			if err != nil {
				log.Err(err).Msg("Failed to reconnect to PostgreSQL database")
			}
		}
	}
}
