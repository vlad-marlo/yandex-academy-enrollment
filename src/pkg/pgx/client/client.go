package client

import (
	"context"
	"fmt"
	pgxzap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/retryer"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

type (
	Config interface {
		URI() string
	}
	Client struct {
		pool *pgxpool.Pool
		log  *zap.Logger
	}
)

const (
	RetryAttempts = 4
	RetryDelay    = 2 * time.Second
)

// New opens new postgres connection, configures it and return prepared client.
func New(lc fx.Lifecycle, cfg Config, log *zap.Logger) (*Client, error) {
	var pool *pgxpool.Pool
	log.Info("initializing postgres client with config", zap.Any("cfg", cfg))

	c, err := pgxpool.ParseConfig(
		cfg.URI(),
	)
	if err != nil {
		return nil, fmt.Errorf("error while parsing db uri: %w", err)
	}

	var lvl = tracelog.LogLevelError
	c.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxzap.NewLogger(log),
		LogLevel: lvl,
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), c)
	if err != nil {
		return nil, fmt.Errorf("postgres: init pgxpool: %w", err)
	}

	cli := &Client{
		pool: pool,
		log:  log,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return retryer.TryWithAttemptsCtx(ctx, pool.Ping, RetryAttempts, RetryDelay)
		},
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})
	log.Info("created postgres client")
	return cli, nil
}

// L return global client logger.
//
// If client is nil object then global logger will be returned.
func (cli *Client) L() *zap.Logger {
	if cli == nil {
		zap.L().Error("unexpectedly got nil client dereference")
		return zap.L()
	}
	return cli.log
}

// P returns client's configured logger.
//
// If client is nil object then will be returned nil pool.
func (cli *Client) P() *pgxpool.Pool {
	if cli == nil {
		zap.L().Error("unexpectedly got nil client dereference")
		return nil
	}
	return cli.pool
}
