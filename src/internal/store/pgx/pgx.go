package pgx

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/service/production"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/pgx"
	"go.uber.org/zap"
)

var _ production.Store = (*Store)(nil)

type Store struct {
	log  *zap.Logger
	pool *pgxpool.Pool
}

func New(cli pgx.Client) *Store {
	return &Store{
		log:  cli.L(),
		pool: cli.P(),
	}
}
